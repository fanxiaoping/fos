package db

import (
	"github.com/fanxiaoping/fos/util/config"
	"gopkg.in/mgo.v2"
	"log"
)

var(
	session *mgo.Session
	database *mgo.Database
)

func init() {
	var err error

	dialInfo := &mgo.DialInfo{
		Addrs: []string{config.String("mg_url")},
		Username:config.String("mg_uid"),
		Password:config.String("mg_pwd"),
		Source:"admin",
		Direct: false,
		PoolLimit:100,
	}
	session, err = mgo.DialWithInfo(dialInfo)
	if err != nil {
		log.Println(err.Error())
	}
	session.SetMode(mgo.Monotonic, true)

	database = session.DB(config.String("mg_database"))
}

/**
*	获取mongodb Database对象
 */
func MG() *mgo.Database {
	return database
}

func MSession() *mgo.Session  {
	return session
}
