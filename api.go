package nfly

import (
	"encoding/json"
	"log"
	"net/url"
	"nfly_client/helper"
	"path"
	"strings"
	"time"
)

type User struct {
	Email    string
	Password string
}

type API struct {
	Server     *url.URL
	Request    *helper.Request
	User       *User
	SessionEXP time.Time
}

type Methods interface {
	Login(usr string, pass string) bool
	Register(usr string, pass string) bool
	Delete(usr string) bool
	Logout() bool
	Feeds() string
	Collect(uid string, ok bool) bool
	Push(formatted string) bool
}

var (
	apiURL = map[string]*url.URL{}
)

func NewAPI(host string) *API { // http will only available in test env
	var server *url.URL
	if testing {
		server, _ = url.Parse("http://" + host)
	} else {
		server, _ = url.Parse("https://" + host)
	}
	registerAPI(server)
	return &API{server, helper.NewRequest(), nil, time.Now()}
}

func registerAPI(server *url.URL) {
	// register base api url
	for _, apiPath := range apiPaths {
		tmp := *server
		tmp.Path = path.Join(tmp.Path, apiPath)
		apiURL[apiPath] = &tmp
	}
}

func (x *API) Login(usr string, pass string) bool {
	data, err := x.Request.PostForm(apiURL["login"].String(), url.Values{
		"email":    []string{usr},
		"password": []string{pass},
	})
	if err != nil {
		log.Printf("%s error: %s\n", "login", err.Error())
		return false
	}
	if data.Status {
		// save credit in case session exp
		x.User = &User{usr, pass}
		x.SessionEXP = time.Now().AddDate(0, 0, 7)
	}
	return data.Status
}

func (x *API) Register(usr string, pass string) bool {
	data, err := x.Request.PostForm(apiURL["register"].String(), url.Values{
		"email":    []string{usr},
		"password": []string{pass},
	})
	if err != nil {
		log.Printf("%s error: %s\n", "register", err.Error())
		return false
	}
	return data.Status
}

func (x *API) Logout() bool {
	if !x.CheckSession() {
		return true
	}
	data, err := x.Request.Get(apiURL["logout"].String())
	if err != nil {
		log.Printf("%s error: %s\n", "logout", err.Error())
		return false
	}
	return data.Status
}

func (x *API) Delete(usr string) bool {
	x.CheckAndUpdateSession()
	target := *apiURL["delete"]
	target.Path = path.Join(target.Path, usr)
	data, err := x.Request.Delete(target.String())
	if err != nil {
		log.Printf("%s error: %s\n", "delete", err.Error())
		return false
	}
	return data.Status
}

func (x *API) Feeds() string {
	x.CheckAndUpdateSession()
	data, err := x.Request.Get(apiURL["feeds"].String())
	if err != nil || !data.Status {
		log.Println("warning: fetch feeds failed")
		return ""
	}
	f, _ := json.Marshal(data.Data)
	return string(f)
}

func (x *API) Push(formatted string) bool {
	x.CheckAndUpdateSession()
	data, err := x.Request.Post(apiURL["push"].String(), "json", strings.NewReader(formatted))
	if err != nil {
		log.Printf("%s error: %s\n", "push", err.Error())
		return false
	}
	return data.Status
}

func (x *API) Collect(uid string, ok bool) bool {
	x.CheckAndUpdateSession()
	target := *apiURL["collect"]
	target.Path = path.Join(target.Path, uid)
	var status []byte
	if ok {
		status = []byte{0x01}
	} else {
		status = []byte{0x00}
	}
	data, err := x.Request.Put(target.String(), status)
	if err != nil {
		log.Printf("%s error: %s\n", "collect", err.Error())
		return false
	}
	return data.Status
}

func (x *API) CheckSession() bool {
	return x.SessionEXP.Sub(time.Now()) > 0
}

func (x *API) CheckAndUpdateSession() {
	if !x.CheckSession() {
		log.Println("session updating..")
		result := x.Login(x.User.Email, x.User.Password)
		log.Printf("session update status: %t", result)
		//return
	}
	//log.Println("session valid")
}
