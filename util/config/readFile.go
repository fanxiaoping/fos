package config

import (
	"errors"
	"github.com/fanxiaoping/fos/util/validator"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

)

var (
	confMap *config
)

//	config	读取配置文件
type config struct {
	Conf map[string]string
	lock *sync.RWMutex
	file string
}

func createConf() *config {
	r := new(config)
	r.Conf = make(map[string]string)
	return r
}

func (c *config) getResource(dir string) error {
	cont, err := ioutil.ReadFile(dir)

	if err != nil {
		return err
	}
	conf := strings.Split(string(cont), "\n")
	c.lock.RLock()
	defer c.lock.RUnlock()

	for _, val := range conf {
		val = strings.TrimSpace(val)
		val = strings.TrimRight(val, "\r")
		a := c.trimSpace(val)
		if len(a) > 0 && a[0] == '#' {
			continue
		}
		key, keyVal, err := c.splitEqual(a)
		if err == nil {
			c.Conf[key] = keyVal
		}
	}
	return nil
}

func (c *config) trimSpace(str string) string {
	var rst []byte

	by := []byte(str)
	bgn := 0
	for _, val := range by {
		if val == '"' && bgn == 0 {
			bgn = 1
		} else if val == '"' && bgn == 1 {
			bgn = 0
		}

		if bgn == 0 && val == ' ' {

		} else {
			rst = append(rst, val)
		}
	}

	return string(rst)
}

func (c *config) Set(key, value string) error {
	c.lock.RLock()
	defer c.lock.RUnlock()
	if _, ok := c.Conf[key]; ok {

		// modify configuration variable values.
		fd, err := os.OpenFile(c.file, os.O_RDWR, 666)
		if err != nil {
			return err
		}
		newValue := key + "=" + value

		rd := make([]byte, 1024)
		type rdstruct struct {
			line   string
			index  int
			offset int
		}

		var lines []rdstruct
		var tmp rdstruct
		var line []byte
		var index = -1
		var offset = 0
		var lstIndex = 0
		for {
			size, err := fd.Read(rd)
			if err != nil && err != io.EOF {
				return err
			}

			if size == 0 {
				if len(line) > 0 {
					tmp.index = lstIndex
					tmp.line = string(line)
					tmp.offset = offset + 1
					lines = append(lines, tmp)
				}
				break
			}

			for _, val := range rd[0:size] {
				if val == '\x0a' {
					index++
					offset++
					tmp.index = lstIndex
					lstIndex = index
					tmp.offset = offset
					if len(line) == 0 {
						continue
					}
					if line[len(line)-1] == '\x0d' {
						tmp.line = string(line[:len(line)-1])
						tmp.offset -= 1
					} else {
						tmp.line = string(line)
					}
					lines = append(lines, tmp)
					line = make([]byte, 0)
					offset = 0
				} else {
					index++
					offset++
					line = append(line, val)
				}
			}
		}

		for _, val := range lines {
			t1 := c.trimSpace(val.line)
			mkey, _, _ := c.splitEqual(t1)

			if mkey == key {
				//te := make([]byte, val.offset)
				var aprst []byte
				fd.Seek(int64(val.index+val.offset), 0)
				for {
					tmp := make([]byte, 1024)
					n, err := fd.Read(tmp)
					if err != io.EOF && err != nil {
						return err
					}
					if n == 0 {
						break
					}
					aprst = append(aprst, tmp[0:n]...)
				}
				if val.index == 0 {
					fd.Seek(int64(val.index), 0)
				} else {
					fd.Seek(int64(val.index+1), 0)
				}

				if len(newValue) < val.offset {
					if val.offset-len(newValue) > 1 {
						var tb []byte = make([]byte, val.offset-len(newValue)-1)
						for i := 0; i < val.offset-len(newValue)-1; i++ {
							tb[i] = '\x20'
						}
						newValue = newValue + string(tb) + "\n"
					} else {
						newValue = newValue + "\n"
					}
				} else {
					newValue = newValue + "\n"
				}

				if len(aprst) > 0 && aprst[0] == '\x0a' {
					fd.WriteString(newValue + string(aprst[1:]))
				} else {
					fd.WriteString(newValue + string(aprst))
				}
				fd.Close()
			} else {
				//fmt.Println(mkey, mkeyVal, err)
			}
		}
		return nil
	} else {
		// add configuration variable values.
		op := key + "=" + value + "\n"
		fd, err := os.OpenFile(c.file, os.O_APPEND, 666)
		_, err = fd.WriteString(op)
		fd.Close()
		return err
	}
	c.getResource(c.file)
	return nil
}

func (c *config) Get(key string) (string, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	//获取环境变量
	if val := os.Getenv(strings.ToUpper(key)); !validator.IsEmpty(val) {
		return val, nil
	} else {
		//获取配置文件中的值
		if val, ok := c.Conf[key]; ok {
			return val, nil
		} else {
			return "", errors.New("cant't find value of " + key)
		}
	}
}

func (c *config) splitEqual(str string) (string, string, error) {
	by := []byte(str)
	bgn := 0
	end := 0
	var key string
	var keyVal string

	for _, val := range by {
		if val == '"' && bgn == 0 {
			bgn = 1
		} else if val == '"' && bgn == 1 {
			bgn = 0
		}

		if bgn == 0 && val == '=' && end == 0 {
			end = 1
		} else if end == 0 {
			key += string(val)
		} else {
			keyVal += string(val)
		}
	}
	if keyVal == "" || key == "" {
		return "", "", errors.New("empty value")
	}
	return strings.Trim(key, "\""), strings.Trim(keyVal, "\""), nil
}

// GetDetails configuration infomation
func GetConfig(path string) (*config, error) {
	conf := createConf()
	conf.file = path
	conf.lock = new(sync.RWMutex)
	err := conf.getResource(path)
	if err != nil {
		return nil, err
	} else {
		return conf, nil
	}
}

// String 获取conf中的值
func String(key string) string {
	if confMap == nil {
		conf, err := GetConfig(path.Join(GetCurrentDirectory(), "conf", "app.conf"))
		if err != nil {
			log.Fatalln("init conf:%s", err)
		}
		confMap = conf
	}
	value, err := confMap.Get(key)
	if err != nil {
		return ""
	}
	return value
}

// GetCurrentDirectory 获取当前文件路径
func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0])) //返回绝对路径  filepath.Dir(os.Args[0])去除最后一个元素的路径
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1) //将\替换成/
}
