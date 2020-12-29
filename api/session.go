package main

import (
	"github.com/tidusant/c3m/common/log"
	pbses "github.com/tidusant/c3m/grpc/protoc/session"
	"github.com/tidusant/c3m/repo/models"
	"golang.org/x/net/context"
	"strings"
	"time"
)

var SessionPool map[string]*models.Session
var timeStoreRequest = 1 * time.Second * 3600 //time to store request for check duplicate, by hours
var request map[string]bool                   //check duplicate request during timeStoreRequest time
var nipcheck = 100                            //number of request per second per address
var ipcheck map[string]int                    //check how many request in 1 second
var sessionTime = 30 * 60 * time.Second       //time of a session to expire
//CreateSession: create session string and save into database
func CreateSession() string {
	//sex := mystring.RandString(20)
	//SessionPool[sex] = &models.Session{Time: time.Now()}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	r, err := sessConn.CreateSession(ctx, &pbses.Void{})
	if err != nil {
		CheckReDialSessionService(err)
		return ""
	}
	return r.Data
}
func CheckSession(s string) bool {
	if s == "" {
		return false
	}
	//old check in memory

	//
	//if _, ok := SessionPool[s]; !ok {
	//	return false
	//}
	////update time session
	//SessionPool[s].Time = time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	r, err := sessConn.CheckSession(ctx, &pbses.DataRequest{Data: s})
	if err != nil {
		CheckReDialSessionService(err)
		return false
	}
	return r.Data
}

func SaveSession(ses *pbses.SessionMessage) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	_, err := sessConn.SaveSession(ctx, ses)
	CheckReDialSessionService(err)

}
func GetSession(sex string) *pbses.SessionMessage {
	//if _, ok := SessionPool[sex]; ok {
	//	return SessionPool[sex]
	//}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	r, err := sessConn.GetSession(ctx, &pbses.DataRequest{Data: sex})
	CheckReDialSessionService(err)
	if err != nil {
		return &pbses.SessionMessage{}
	}
	return r
}

//CheckRequest: check request for anti ddos with request limit from env: REQUEST_LIMIT
func CheckRequest(uri, remoteAddress string) bool {
	if uri == "" {
		return false
	}
	////check duplicate request
	//if _, ok := request[uri]; ok {
	//	return false
	//}
	//request[uri] = true
	////if in s second, there are more than n request in one ipaddress -> return
	//if _, ok := ipcheck[remoteAddress]; !ok {
	//	ipcheck[remoteAddress] = 1
	//}
	//if ipcheck[remoteAddress] > nipcheck {
	//	return false
	//}
	//ipcheck[remoteAddress]++

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	r, err := sessConn.CheckRequest(ctx, &pbses.CheckURLRequest{URL: uri, Address: remoteAddress})

	if err != nil {
		CheckReDialSessionService(err)
		return false
	}
	return r.Data

}
func CheckReDialSessionService(err error) {
	if err != nil {
		if strings.Index(err.Error(), "Error while dialing dial") > 0 {
			//delay time to reconnect grpc again
			go func(name string) {
				log.Debugf("wait to reconnect %s grpc", name)
				time.Sleep(time.Second * 5)
				registerGrpc(context.Background(), name)
			}("ses")
		}
		log.Debugf("Error while call session rpc:%s", err.Error())
	}

}
func initCheckSession() {
	SessionPool = make(map[string]*models.Session)
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
				if time.Since(SessionPool[k].Time) > sessionTime {
					SessionPool[k] = nil
					delete(SessionPool, k)
				}
			}
			time.Sleep(timeStoreRequest)
		}
	}()
}
