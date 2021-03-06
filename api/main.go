package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/tidusant/c3m/common/c3mcommon"
	"github.com/tidusant/c3m/common/log"
	"github.com/tidusant/c3m/common/mycrypto"
	pb "github.com/tidusant/c3m/grpc/protoc"
	pbses "github.com/tidusant/c3m/grpc/protoc/session"
	"github.com/tidusant/c3m/repo/models"
	"google.golang.org/grpc"
	"sync"

	"os"
	"time"

	//"io" repush
	"net"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func init() {
	log.Debug("main init")

}

var grpcConns map[string]pb.GRPCServicesClient
var sessConn pbses.SessionServicesClient
var grpcAddressMap map[string]string
var exposeport = "8081"
var grpcLocker sync.Mutex

//main function app run here
func main() {

	//init grpc
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	connectGrpcs(ctx)
	//initCheckSession()

	//start gin
	router := gin.Default()
	router.POST("/*name", postHandler)
	router.Run(":" + exposeport)

}

func connectGrpcs(ctx context.Context) {
	//check all grpc is running:
	grpcAddressMap = make(map[string]string)
	grpcConns = make(map[string]pb.GRPCServicesClient)
	grpcAddressMap["ses"] = os.Getenv("SESSION_IP")
	grpcAddressMap["aut"] = os.Getenv("AUTH_IP")
	grpcAddressMap["shop"] = os.Getenv("SHOP_IP")
	grpcAddressMap["ord"] = os.Getenv("ORD_IP")
	grpcAddressMap["page"] = os.Getenv("PAGE_IP")
	grpcAddressMap["lptpl"] = os.Getenv("LPTPL_IP")
	grpcAddressMap["lp"] = os.Getenv("LP_IP")
	grpcAddressMap["lpl"] = os.Getenv("LPL_IP")

	//init grpc ssession
	registerSessionGrpc(ctx, "ses")

	//implement concurrency
	for name, add := range grpcAddressMap {
		go func(name, add string) {
			if len(add) < 10 {
				fmt.Printf("%sIP is invalid:%s", name, add)
				return
			}
			registerGrpc(ctx, name)
		}(name, add)
	}

}

func registerGrpc(ctx context.Context, name string) {
	grpcLocker.Lock()
	fmt.Printf("Register grpc %s at %s\n", name, grpcAddressMap[name])
	conn, err := grpc.Dial(grpcAddressMap[name], grpc.WithInsecure())
	//defer conn.Close()

	if err != nil {
		fmt.Printf("Warning: can not call grpc auth %v", err)
	} else {
		grpcConns[name] = pb.NewGRPCServicesClient(conn)
	}
	grpcLocker.Unlock()
}

func registerSessionGrpc(ctx context.Context, name string) {
	fmt.Printf("Register grpc session at %s\n", grpcAddressMap[name])
	conn, err := grpc.Dial(grpcAddressMap[name], grpc.WithInsecure())
	//defer conn.Close()

	if err != nil {
		fmt.Printf("Warning: can not call grpc auth %v", err)
	} else {
		sessConn = pbses.NewSessionServicesClient(conn)
	}
}

func postHandler(c *gin.Context) {
	strrt := ""
	requestDomain := c.Request.Header.Get("Origin")
	c.Header("Access-Control-Allow-Origin", requestDomain)
	c.Header("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers,access-control-allow-credentials")
	c.Header("Access-Control-Allow-Credentials", "true")

	//check request url, only one unique url per second

	//if rpsex.CheckRequest(c.Request.URL.Path, c.Request.UserAgent(), c.Request.Referer(), c.Request.RemoteAddr, "POST") {

	if CheckRequest(c.Request.URL.Path, c.Request.RemoteAddr) {

		rs := myRoute(c)
		start2 := time.Now()

		b, _ := json.Marshal(rs)
		strrt = string(b)

		log.Debugf("marshall time:%s", time.Since(start2))
	} else {
		log.Debugf("request denied")
		strrt = `{"Status":0,"Error":"request denied"}`
	}

	if strrt == "" {
		strrt = c3mcommon.Fake64()
	} else {
		strrt = mycrypto.Encode(strrt, 8)
	}
	c.String(http.StatusOK, strrt)
}

func callgRPC(name string, rpcRequest pb.RPCRequest) models.RequestResult {
	start := time.Now()
	rs := models.RequestResult{Error: ""}
	if _, ok := grpcConns[name]; !ok {
		rs.Error = "Cannot found grpc " + name
		return rs
	}
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	r, err := grpcConns[name].Call(ctx, &rpcRequest)
	if err != nil {

		if strings.Index(err.Error(), "Error while dialing dial") > 0 {
			rs.Error = "service not run, please try again after 5 second"
			//delay time to reconnect grpc again
			go func(name string) {
				log.Debugf("wait to reconnect %s grpc", name)
				time.Sleep(time.Second * 5)
				registerGrpc(context.Background(), name)
			}(name)
		} else {
			rs.Error = err.Error()
		}
		return rs
	}

	err = json.Unmarshal([]byte(r.Data), &rs)
	if err != nil {
		rs.Error = "Parse data string from service Error."
	}
	duration, _ := time.ParseDuration(r.GetTime())

	log.Debugf("callgRPC %s query time:%s", name, duration)
	log.Debugf("callgRPC %s query count:%d", name, r.GetQuery())
	log.Debugf("total callgRPC time:%s", time.Since(start))
	return rs
}

func myRoute(c *gin.Context) models.RequestResult {
	//get request name
	start := time.Now()
	name := c.Param("name")
	name = name[1:] //remove  slash

	//get request data from Form
	data := c.PostForm("data")

	//get userip for check on 1 ip login
	userIP, _, _ := net.SplitHostPort(c.Request.RemoteAddr)
	//log.Debugf("decode name:%s", mycrypto.Decode(name))
	//decode request name and get array of args
	args := strings.Split(mycrypto.Decode(name), "|")
	RPCname := args[0]
	AppName := "admin-portal"
	if len(args) > 1 {
		AppName = args[1]
	}
	//decode request data from Form and get array of args

	datargs := strings.Split(mycrypto.Decode(data), "|")
	session := mycrypto.Decode(datargs[0])
	log.Debugf("datargs:%+v", datargs)

	requestAction := ""
	requestParams := ""
	if len(datargs) > 1 {
		requestAction = datargs[1]
	}
	if len(datargs) > 2 {
		requestParams = datargs[2]
	}

	//get rpc call name from first arg
	log.Debugf("session: %+v", session)
	log.Debugf("RPCname:%s, action:%s, app:%s, params:%s", RPCname, requestAction, AppName, requestParams)
	if RPCname == "CreateSex" {
		//create session string and save it into db
		//data = rpsex.CreateSession()
		sex := CreateSession()
		if sex == "" {
			return models.RequestResult{Error: "Cannot create session"}
		}
		return models.RequestResult{Status: 1, Error: "", Data: sex}
	}

	//check session

	//if !rpsex.CheckSession(session) {
	sex := GetSession(session)

	if sex.Session == "" {
		return models.RequestResult{Status: -1, Error: "Session not found"}
	}
	log.Debugf("session found: %+v", sex)
	reply := models.RequestResult{Error: ""}
	if RPCname == "aut" {
		if requestAction == "t" {
			reply = models.RequestResult{Status: 1, Error: "", Data: `{"sex":"` + session + `","username":"` + sex.UserName + `","shop":"` + sex.ShopID + `","group":"` + sex.Group + `","modules":"` + sex.Modules + `"}`}
		} else {
			reply = callgRPC("aut", pb.RPCRequest{AppName: AppName, Action: requestAction, Params: requestParams, Session: session, UserIP: userIP})
			if reply.Status == 1 {
				if requestAction == "l" {
					var data map[string]string
					c3mcommon.CheckError("parse login response", json.Unmarshal([]byte(reply.Data), &data))
					//save userinfo int session
					SaveSession(&pbses.SessionMessage{Session: session, UserID: data["userid"], UserName: data["username"], ShopID: data["shopid"], Group: data["group"], Modules: data["modules"]})
					//remove userid & shopid in reply
					if AppName == `sf` {
						reply.Data = fmt.Sprintf(`{"sex":"%s","username":"%s","group":"%s","modules":"%s","shop":"%s"}`, session, data["username"], data["group"], data["modules"], data["shopid"])
					} else {
						reply.Data = fmt.Sprintf(`{"username":"%s","group":"%s","modules":"%s","shop":"%s"}`, data["username"], data["group"], data["modules"], data["shopid"])
					}
				}
			}
		}
		return reply
	}

	//check auth
	if sex.UserName == "" {
		reply.Error = "Not authorize"
		return reply
	}

	if RPCname == "shop" && requestAction == "cs" {
		//change shop
		reply = callgRPC(RPCname, pb.RPCRequest{AppName: AppName, Action: requestAction, Params: requestParams, Session: session, UserID: sex.UserID, UserIP: userIP, ShopID: sex.ShopID, Group: sex.Group, Username: sex.UserName, Modules: sex.Modules})
		if reply.Status == 1 {
			//update session
			sex.ShopID = requestParams
			SaveSession(sex)
		}
		return reply

	}

	//normal gRPC call

	log.Debugf("RPCname: %s", RPCname)
	//time.Sleep(0 * time.Second)

	reply = callgRPC(RPCname, pb.RPCRequest{AppName: AppName, Action: requestAction, Params: requestParams, Session: session, UserID: sex.UserID, UserIP: userIP, ShopID: sex.ShopID, Group: sex.Group, Username: sex.UserName, Modules: sex.Modules})
	log.Debugf("total time:%s", time.Since(start))

	return reply
}
