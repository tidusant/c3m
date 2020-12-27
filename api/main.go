package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/tidusant/c3m/common/c3mcommon"
	"github.com/tidusant/c3m/common/log"
	"github.com/tidusant/c3m/common/mycrypto"
	pb "github.com/tidusant/c3m/grpc/protoc"
	"github.com/tidusant/c3m/repo/models"
	"google.golang.org/grpc"

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
var grpcAddressMap map[string]string
var exposeport = "8081"

//main function app run here
func main() {

	//init grpc
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	connectGrpcs(ctx)
	initCheckSession()

	//start gin
	router := gin.Default()
	router.POST("/*name", postHandler)
	router.Run(":" + exposeport)

}

func connectGrpcs(ctx context.Context) {
	//check all grpc is running:
	grpcAddressMap = make(map[string]string)
	grpcConns = make(map[string]pb.GRPCServicesClient)
	grpcAddressMap["aut"] = os.Getenv("AUTH_IP")
	grpcAddressMap["shop"] = os.Getenv("SHOP_IP")
	grpcAddressMap["ord"] = os.Getenv("ORD_IP")

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
	fmt.Printf("Register grpc %s at %s\n", name, grpcAddressMap[name])
	conn, err := grpc.Dial(grpcAddressMap[name], grpc.WithInsecure())
	//defer conn.Close()

	if err != nil {
		fmt.Printf("Warning: can not call grpc auth %v", err)
	} else {
		grpcConns[name] = pb.NewGRPCServicesClient(conn)
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
	}

	if strrt == "" {
		strrt = c3mcommon.Fake64()
	} else {
		start := time.Now()
		strrt = mycrypto.Encode(strrt, 8)
		log.Debugf("encode time:%s", time.Since(start))
	}
	c.String(http.StatusOK, strrt)
}

func callgRPC(name string, rpcRequest pb.RPCRequest) models.RequestResult {
	start := time.Now()
	rs := models.RequestResult{Error: "service not run, please try again after 5 second"}
	if _, ok := grpcConns[name]; !ok {
		return rs
	}
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()
	r, err := grpcConns[name].Call(ctx, &rpcRequest)
	if err != nil {

		if strings.Index(err.Error(), "Error while dialing dial") > 0 {
			//delay time to reconnect grpc again
			go func(name string) {
				log.Debugf("wait to reconnect %s grpc", name)
				time.Sleep(time.Second * 5)
				registerGrpc(context.Background(), name)
			}(name)
		} else {
			rs.Error = "Error while parsing response data:" + err.Error()
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
	//decode request data from Form and get array of args

	datargs := strings.Split(mycrypto.Decode(data), "|")
	session := mycrypto.Decode(datargs[0])
	log.Debugf("total decode time:%s", time.Since(start))

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
	log.Debugf("RPCname:%s, action:%s", RPCname, requestAction)
	if RPCname == "CreateSex" {
		//create session string and save it into db
		//data = rpsex.CreateSession()
		return models.RequestResult{Status: 1, Error: "", Data: CreateSession()}
	}

	//check session

	//if !rpsex.CheckSession(session) {
	if !CheckSession(session) {
		return models.RequestResult{Status: -1, Error: "Session not found"}
	}
	log.Debugf("session found: %+v", SessionPool[session])

	if RPCname == "aut" && requestAction == "l" {
		reply := callgRPC("aut", pb.RPCRequest{AppName: "admin-portal", Action: requestAction, Params: requestParams, Session: session, UserIP: userIP})
		if reply.Status == 1 {
			var data map[string]string
			c3mcommon.CheckError("parse login response", json.Unmarshal([]byte(reply.Data), &data))
			//save userinfo int session
			SessionPool[session].shopid = data["shopid"]
			SessionPool[session].username = data["username"]
			SessionPool[session].userid = data["userid"]
			//remove userid & shopid in reply
			reply.Data = fmt.Sprintf(`{"username":"%s"}`, data["username"])
		}
		return reply
	}

	//always check login if RPCname not aut and create session
	//reply := callgRPC("aut", pb.RPCRequest{AppName: "admin-portal", Action: "aut", Params: requestParams, Session: session, UserIP: userIP})
	//if reply.Status != 1 {
	//	return reply
	//}
	reply := models.RequestResult{Error: ""}
	sex := GetSession(session)
	if sex == nil {
		reply.Error = "No session found"
		return reply
	}

	//some specific function relate to Session
	if RPCname == "shop" && requestAction == "cs" {
		//change shop
		reply = callgRPC(RPCname, pb.RPCRequest{AppName: "admin-portal", Action: requestAction, Params: requestParams, Session: session, UserID: sex.userid, UserIP: userIP, ShopID: sex.shopid})
		if reply.Status == 1 {
			//update session
			SessionPool[session].shopid = reply.Data
		}
		log.Debugf("cs return:%+v", reply)
	} else if RPCname == "aut" && requestAction == "t" {
		reply = models.RequestResult{Status: 1, Error: "", Data: `{"sex":"` + session + `","username":"` + sex.username + `","shop":"` + sex.shopid + `"}`}
		log.Debugf("t return:%+v", reply)
	} else {

		//normal gRPC call
		log.Debugf("RPCname: %s", RPCname)
		//time.Sleep(0 * time.Second)

		reply = callgRPC(RPCname, pb.RPCRequest{AppName: "admin-portal", Action: requestAction, Params: requestParams, Session: session, UserID: sex.userid, UserIP: userIP, ShopID: sex.shopid})
		log.Debugf("total time:%s", time.Since(start))
	}

	return reply
}
