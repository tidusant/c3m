package main

import (
	"github.com/tidusant/c3m/common/mystring"
	"time"
)

type Session struct {
	time     time.Time
	userid   string
	shopid   string
	username string
}

var SessionPool map[string]*Session
var timeStoreRequest = 1 * time.Second * 3600 //time to store request for check duplicate, by hours
var request map[string]bool                   //check duplicate request during timeStoreRequest time
var nipcheck = 100                            //number of request per second per address
var ipcheck map[string]int                    //check how many request in 1 second
var sessionTime = 30 * 60 * time.Second       //time of a session to expire
//CreateSession: create session string and save into database
func CreateSession() string {
	sex := mystring.RandString(20)
	SessionPool[sex] = &Session{time: time.Now()}
	return sex
}
func CheckSession(s string) bool {
	if s == "" {
		return false
	}

	if _, ok := SessionPool[s]; !ok {
		return false
	}
	//update time session
	SessionPool[s].time = time.Now()
	return true
}
func SaveSession(sex string, ses Session) {
	if _, ok := SessionPool[sex]; ok {
		SessionPool[sex] = &ses
	}
}
func GetSession(sex string) *Session {
	if _, ok := SessionPool[sex]; ok {
		return SessionPool[sex]
	}
	return nil
}

//CheckRequest: check request for anti ddos with request limit from env: REQUEST_LIMIT
func CheckRequest(uri, remoteAddress string) bool {
	if uri == "" {
		return false
	}
	//check duplicate request
	if _, ok := request[uri]; ok {
		return false
	}
	request[uri] = true
	//if in s second, there are more than n request in one ipaddress -> return
	if _, ok := ipcheck[remoteAddress]; !ok {
		ipcheck[remoteAddress] = 1
	}
	if ipcheck[remoteAddress] > nipcheck {
		return false
	}
	ipcheck[remoteAddress]++
	return true
}

func initCheckSession() {
	SessionPool = make(map[string]*Session)
	//service to remove all request in timeStoreRequest
	go func() {
		for {
			request = make(map[string]bool)
			time.Sleep(timeStoreRequest)
		}
	}()
	//service to remove all ipcheck in 1 second
	go func() {
		for {
			ipcheck = make(map[string]int)
			time.Sleep(time.Second)
		}
	}()
	//service to remove all expire session in every 5 second
	go func() {
		for {
			for k, _ := range SessionPool {
				if time.Since(SessionPool[k].time) > sessionTime {
					SessionPool[k] = nil
					delete(SessionPool, k)
				}
			}
			time.Sleep(timeStoreRequest)
		}
	}()
}
