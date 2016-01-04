package main

import (
	"encoding/json"
	"go-venus-plug/models"
)

func main() {
	venus := Conncect("10.32.172.124:16800")
	venus.AuthByDummy("venus")
	var userData models.FindNameData
	userData.UserName = "cscadmin"
	body, err := json.Marshal(userData)
	if err != nil {
		panic(err.Error())
	}
	venus.Request("permissionServiceOP.findUserByUserName", "1", string(body))
}
