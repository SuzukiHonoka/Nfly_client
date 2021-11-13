package test

import (
	"encoding/json"
	nfly "nfly_client"
	"testing"
)

const (
	pushPayload = "{\"Title\":\"test\",\"MessageChain\":[{\"Type\":\"text\",\"Data\":\"test\"},{\"Type\":\"binary\",\"Data\":\"AQIBAgMEBQY=\"}]}"
)

func TestAPI(t *testing.T) {
	user,pass := "test5","1"
	api := nfly.NewAPI("127.0.0.1:7700")
	// Register test
	if !api.Register(user,pass){
		t.Error("cannot register")
	}
	t.Log("register test ok")
	// Login test
	if !api.Login(user,pass){
		t.Error("cannot login")
	}
	t.Log("login test ok")
	// Push test
	if !api.Push(pushPayload) {
		t.Error("cannot push")
	}
	t.Log("push test ok")
	// Feed test
	var uid string
	if f:=api.Feed();len(f) == 0 {
		t.Error("cannot fetch feeds")
	}else {
		//t.Log(f)
		var data []interface{}
		_ = json.Unmarshal([]byte(f),&data)
		//t.Logf("%v+",data)
		uid = data[0].(map[string]interface{})["Header"].(map[string]interface{})["UUID"].(string)
		//t.Log(uid)
	}
	t.Log("feed test ok")
	// Collect test
	if !api.Collect(uid,true) {
		t.Error("cannot collect")
	}
	t.Log("collect test ok")
	// Delete test
	if !api.Delete(user){
		t.Error("cannot delete")
	}
	t.Log("delete test ok")
}