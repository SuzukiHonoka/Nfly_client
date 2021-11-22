package test

import (
	"encoding/json"
	nfly "nfly_client"
	"testing"
)

const (
	pushPayload = "{\"Title\":\"test\",\"MessageChain\":[{\"Type\":\"text\",\"Data\":\"test\"},{\"Type\":\"binary\",\"Data\":\"AQIBAgMEBQY=\"}]}"
)

func realTest(api nfly.Methods, user string, pass string) {
	// Register test
	if !api.Register(user, pass) {
		panic("cannot register")
	}
	// Login test
	if !api.Login(user, pass) {
		panic("cannot login")
	}
	// Logout test
	if !api.Logout() {
		panic("cannot logout")
	}
	if !api.Login(user, pass) {
		panic("cannot login")
	}
	// Push test
	if !api.Push(pushPayload) {
		panic("cannot push")
	}
	// Feeds test
	var uid string
	if f := api.Feeds(); len(f) == 0 {
		panic("cannot fetch feeds")
	} else {
		//t.Log(f)
		var data []interface{}
		_ = json.Unmarshal([]byte(f), &data)
		//t.Logf("%v+",data)
		uid = data[0].(map[string]interface{})["Header"].(map[string]interface{})["UUID"].(string)
		//t.Log(uid)
	}
	// Collect test
	if !api.Collect(uid, true) {
		panic("cannot collect")
	}
	// Delete test
	if !api.Delete(user) {
		panic("cannot delete")
	}
}

func TestAPI(t *testing.T) {
	api := nfly.NewAPI("127.0.0.1:7700")
	realTest(api, "test99", "test99")
}
