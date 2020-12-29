package main

import (
	"github.com/tidusant/c3m/common/mystring"
	pb "github.com/tidusant/c3m/grpc/protoc/session"
	"sync"
	"time"

	"context"
	"google.golang.org/grpc"
	"os"

	"github.com/tidusant/c3m/common/log"

	"github.com/tidusant/c3m/repo/models"

	"fmt"
	"net"
)

var (
	SessionPool map[string]*models.Session

	timeStoreRequest = 1 * time.Second * 3600 //time to store request for check duplicate, by hours
	request          map[string]bool          //check duplicate request during timeStoreRequest time
	nipcheck         = 100                    //number of request per 10 second per address
	ipcheck          map[string]int           //check how many request in 1 second
	blockip          map[string]bool          //block ip multi request in 180 s
	blockTime        = 3 * 60 * time.Second
	sessionTime      = 30 * 60 * time.Second //time of a session to expire
	l                = sync.Mutex{}
)

type service struct {
	pb.UnimplementedSessionServicesServer
}

func (s *service) CreateSession(ctx context.Context, in *pb.Void) (*pb.StringResponse, error) {
	rs := &pb.StringResponse{}
	sex := mystring.RandString(20)
	SessionPool[sex] = &models.Session{Time: time.Now()}
	rs.Data = sex
	return rs, nil
}
func (s *service) CheckSession(ctx context.Context, in *pb.DataRequest) (*pb.BoolResponse, error) {
	rs := &pb.BoolResponse{}
	rs.Data = false
	if _, ok := SessionPool[in.Data]; ok && in.Data != "" {
		rs.Data = true
		//update time session
		SessionPool[in.Data].Time = time.Now()
	}
	return rs, nil
}
func (s *service) CheckRequest(ctx context.Context, in *pb.CheckURLRequest) (*pb.BoolResponse, error) {
	rs := &pb.BoolResponse{}
	rs.Data = true
	//check duplicate request
	l.Lock()
	if _, ok := blockip[in.Address]; ok {
		rs.Data = false
	} else if _, ok := request[in.URL]; ok || in.URL == "" {
		rs.Data = false
	} else {
		request[in.URL] = true
		//if in s second, there are more than n request in one ipaddress -> return

		if _, ok := ipcheck[in.Address]; !ok {
			ipcheck[in.Address] = 1
		}

		if ipcheck[in.Address] > nipcheck {
			//block ip
			blockip[in.Address] = true
			rs.Data = false
		} else {
			ipcheck[in.Address]++
		}

	}
	l.Unlock()
	return rs, nil
}
func (s *service) GetSession(ctx context.Context, in *pb.DataRequest) (*pb.SessionMessage, error) {
	rs := &pb.SessionMessage{}
	if _, ok := SessionPool[in.Data]; ok {
		rs.Session = in.Data
		rs.UserID = SessionPool[in.Data].UserID
		rs.UserName = SessionPool[in.Data].UserName
		rs.ShopID = SessionPool[in.Data].ShopID
		//update session time expire
		SessionPool[in.Data].Time = time.Now()
	}
	return rs, nil
}
func (s *service) SaveSession(ctx context.Context, in *pb.SessionMessage) (*pb.BoolResponse, error) {
	rs := &pb.BoolResponse{}
	if _, ok := SessionPool[in.Session]; ok {
		SessionPool[in.Session].UserName = in.UserName
		SessionPool[in.Session].UserID = in.UserID
		SessionPool[in.Session].ShopID = in.ShopID

		rs.Data = true
	}
	return rs, nil
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
	//service to remove all ipcheck in 10 second
	go func() {
		for {
			ipcheck = make(map[string]int)
			time.Sleep(time.Second * 10)
		}
	}()
	//service to remove all blockip in 180 second
	go func() {
		for {
			blockip = make(map[string]bool)
			time.Sleep(blockTime)
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
func main() {
	//default port for service
	var port string
	port = os.Getenv("PORT")
	if port == "" {
		port = "8864"
	}
	initCheckSession()
	//open service and listen
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Errorf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	fmt.Printf("listening on %s\n", port)
	pb.RegisterSessionServicesServer(s, &service{})
	if err := s.Serve(lis); err != nil {
		log.Errorf("failed to serve : %v", err)
	}
	fmt.Print("exit")

}
