package main

import (
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net"
	"os"
	"strings"
	"time"

	"github.com/tidusant/c3m/common/log"

	pb "github.com/tidusant/c3m/grpc/protoc"
	"github.com/tidusant/c3m/repo/cuahang"
	"github.com/tidusant/c3m/repo/models"

	"context"
	"google.golang.org/grpc"
)

const (
	name string = "auth"
	ver  string = "1"
)

type service struct {
	pb.UnimplementedGRPCServicesServer
}
type Auth struct {
	rpch cuahang.Repo
}

func (s *service) Call(ctx context.Context, in *pb.RPCRequest) (*pb.RPCResponse, error) {
	start := time.Now()
	resp := &pb.RPCResponse{Data: `{"Status":1,"Data":"Hello ` + in.GetAppName() + `"}`, RPCName: name, Version: ver}
	rs := models.RequestResult{}
	a := Auth{}
	//get input data into user session
	var usex models.UserSession
	usex.Session = in.Session
	usex.Action = in.Action
	usex.AppName = in.AppName
	usex.UserIP = in.UserIP
	usex.Params = in.Params
	usex.UserID, _ = primitive.ObjectIDFromHex(in.UserID)
	if usex.Action == "l" {
		rs = a.login(usex)
	} else if usex.Action == "i" {
		rs = a.install(usex)
		//} else if usex.Action == "aut" {
		//	rs = a.auth(usex)
	} else if usex.Action == "lo" {
		rs = a.logout(usex.Session)
		//} else if usex.Action == "aut" {
		//	rs = a.auth(usex)
	} else {
		//unknow action
		return resp, nil
	}
	//convert RequestResult into json
	b, _ := json.Marshal(rs)
	resp.Data = string(b)
	resp.Query = int32(a.rpch.QueryCount)
	resp.Time = time.Since(start).String()
	return resp, nil
}

//auth: to authenticate user already login, for portal, return userid[+]shopid
//func (a *Auth) auth(usex models.UserSession) models.RequestResult {
//	rs := a.rpch.GetLogin(usex.Session)
//	if rs.UserId == primitive.NilObjectID {
//		return models.RequestResult{Status: -1, Error: "user not logged in"}
//	} else {
//		user := a.rpch.GetUserInfo(rs.UserId)
//		return models.RequestResult{Error: "", Status: 1, Data: `{"userid":"` + rs.UserId.Hex() + `","name":"` + user.Name + `","sex":"` + usex.Session + `","shop":"` + rs.ShopId.Hex() + `"}`}
//	}
//}

//login user and update Session and IP in user_login. then return auth call to get userid[+]shopid
func (a *Auth) login(usex models.UserSession) models.RequestResult {

	args := strings.Split(usex.Params, ",")
	if len(args) < 2 {
		return models.RequestResult{Error: "empty username or pass"}
	}
	username := args[0]
	pass := args[1]
	user := a.rpch.Login(username, pass, usex.Session, usex.UserIP)
	if user.Name != "" {
		//get shop default
		if user.ShopId == primitive.NilObjectID {
			shop := a.rpch.GetShopDefault(user.ID)
			user.ShopId = shop.ID
		}
		return models.RequestResult{
			Error:   "",
			Status:  1,
			Message: "logged in",
			Data:    `{"username":"` + user.Name + `","userid":"` + user.ID.Hex() + `","shopid":"` + user.ShopId.Hex() + `","group":"` + user.Group + `","modules":"` + user.Modules + `"}`}

	}
	return models.RequestResult{Error: "Login failed"}

}

func (a *Auth) logout(session string) models.RequestResult {
	a.rpch.Logout(session)
	return models.RequestResult{Error: "", Status: 1, Message: "Logout success"}

}

func (a *Auth) install(usex models.UserSession) models.RequestResult {
	args := strings.Split(usex.Params, ",")
	if usex.AppName == "sf" {
		if len(args) < 4 {
			return models.RequestResult{Error: "invalid install"}
		}
		orgId := args[0]
		orgName := args[1]
		userId := args[2]
		userEmail := args[3]
		if orgId != "" && userId != "" {
			a.rpch.InstallSaleForce(orgId, orgName, userEmail)
			return models.RequestResult{
				Error:   "",
				Status:  1,
				Message: "",
				Data:    ``}

		}
	}
	return models.RequestResult{Error: "invalid app"}
}
func main() {
	//default port for service
	var port string
	port = os.Getenv("PORT")
	if port == "" {
		port = "8901"
	}
	//open service and listen
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Errorf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	fmt.Printf("listening on %s\n", port)
	pb.RegisterGRPCServicesServer(s, &service{})
	if err := s.Serve(lis); err != nil {
		log.Errorf("failed to serve : %v", err)
	}
	fmt.Print("exit")

}
