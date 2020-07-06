package convert

import (
	"fmt"
	"testing"
)

func TestOfficeToPNG(t *testing.T){
	res,err := OfficeToPNG("","")
	if err != nil{
		fmt.Println(err)
	}
	fmt.Println(res)
}