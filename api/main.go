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
	rpsex "github.com/tidusant/c3m/repo/session"
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

var authIP = "127.0.0.1:8901"

var exposeport = "8081"

//main function app run here
func main() {

	//get address of grpc auth from ENV
	authIP = os.Getenv("AUTH_IP")

	//show info to console
	fmt.Println("auth address: " + authIP)
	fmt.Println("\n portal admin running with port " + exposeport)

	//start gin
	router := gin.Default()
	router.POST("/*name", postHandler)
	router.Run(":" + exposeport)

}

func postHandler(c *gin.Context) {
	strrt := ""
	requestDomain := c.Request.Header.Get("Origin")
	c.Header("Access-Control-Allow-Origin", requestDomain)
	c.Header("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers,access-control-allow-credentials")
	c.Header("Access-Control-Allow-Credentials", "true")

	//check request url, only one unique url per second
	if rpsex.CheckRequest(c.Request.URL.Path, c.Request.UserAgent(), c.Request.Referer(), c.Request.RemoteAddr, "POST") {
		rs := myRoute(c)
		b, _ := json.Marshal(rs)
		strrt = string(b)
	} else {
		log.Debugf("request denied")
	}

	if strrt == "" {
		strrt = c3mcommon.Fake64()
	} else {
		strrt = mycrypto.Encode(strrt, 8)
		log.Debug(mycrypto.DecodeOld(strrt, 8))
	}
	c.String(http.StatusOK, strrt)
}

func callgRPC(address string, rpcRequest pb.RPCRequest) models.RequestResult {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	defer conn.Close()

	rs := models.RequestResult{Error: "service not run"}
	if err == nil && len(address) > 10 {
		rpc := pb.NewGRPCServicesClient(conn)
		// Contact the server and print out its response.
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
		defer cancel()
		r, err := rpc.Call(ctx, &rpcRequest)
		if err != nil {
			rs.Error = err.Error()
			return rs
		}

		err = json.Unmarshal([]byte(r.Data), &rs)
		if err != nil {
			rs.Error = r.Data
		}
	}

	return rs
}

func myRoute(c *gin.Context) models.RequestResult {
	//get request name
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
		data = rpsex.CreateSession()
		return models.RequestResult{Status: 1, Error: "", Data: data}
	}

	//check session

	if !rpsex.CheckSession(session) {
		return models.RequestResult{Status: -1, Error: "Session not found"}
	}
	if RPCname == "aut" && requestAction == "l" {
		return callgRPC(authIP, pb.RPCRequest{AppName: "admin-portal", Action: requestAction, Params: requestParams, Session: session, UserIP: userIP})
	}

	//always check login if RPCname not aut and create session
	reply := callgRPC(authIP, pb.RPCRequest{AppName: "admin-portal", Action: "aut", Params: requestParams, Session: session, UserIP: userIP})
	if reply.Status != 1 {
		return reply
	}
	log.Debugf("authentication: %+v", reply)
	//get logininfo: from check login in format: userid[+]shopid
	var rs map[string]string
	json.Unmarshal([]byte(reply.Data), &rs)

	ShopId := rs["shop"]
	UserId := rs["userid"]

	//test function
	if requestAction == "t" {
		return models.RequestResult{Status: 1, Error: "", Data: `{"sex":"` + session + `","name":"` + rs["name"] + `","shop":"` + ShopId + `"}`}

	}

	//begin gRPC call
	log.Debugf("RPCname: %s", RPCname)
	//time.Sleep(0 * time.Second)
	RPCname = strings.ToUpper(RPCname)
	grpcIP := os.Getenv(RPCname + "_IP")
	return callgRPC(grpcIP, pb.RPCRequest{AppName: "admin-portal", Action: requestAction, Params: requestParams, Session: session, UserID: UserId, UserIP: userIP, ShopID: ShopId})

}
