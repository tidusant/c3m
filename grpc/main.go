package grpc

import (
	"context"
	"encoding/json"
	"github.com/tidusant/c3m/common/c3mcommon"
	"github.com/tidusant/c3m/common/log"
	pb "github.com/tidusant/c3m/grpc/protoc"
	"github.com/tidusant/c3m/repo/cuahang"
	"github.com/tidusant/c3m/repo/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
	"time"
)

type MainRPC struct {
	Usex       models.UserSession
	Rpch       cuahang.Repo
	QueryCount int
	ctx        context.Context
	resp       *pb.RPCResponse
	start      time.Time
}

func (m *MainRPC) InitUsex(ctx context.Context, in *pb.RPCRequest, name, ver string) string {

	m.start = time.Now()
	rt := ""
	m.resp = &pb.RPCResponse{Data: `{"Status":1,"Data":"Hello ` + in.GetAppName() + `"}`, RPCName: name, Version: ver}
	//get input data into user session
	m.ctx = ctx
	m.Usex.Session = in.Session
	m.Usex.Action = in.Action
	m.Usex.AppName = in.AppName
	m.Usex.UserIP = in.UserIP
	m.Usex.Group = in.Group
	m.Usex.Params = in.Params
	m.Usex.Username = in.Username
	m.Usex.Modules = make(map[string]bool)
	for _, v := range strings.Split(in.Modules, ",") {
		m.Usex.Modules[v] = true
	}
	m.Usex.UserID, _ = primitive.ObjectIDFromHex(in.UserID)
	m.Usex.ShopID, _ = primitive.ObjectIDFromHex(in.ShopID)
	log.Debugf("action:%s - userid:%s - group:%s - modules:%+v", in.Action, in.UserID, in.Group, in.Modules)
	//check shop permission
	//if in.ShopID != "" {
	//	shopidObj, _ := primitive.ObjectIDFromHex(in.ShopID)
	//	shop := m.Rpch.GetShopById(m.Usex.UserID, shopidObj)
	//	if shop.Status == 0 {
	//		if m.Usex.Action != "lsi" {
	//			rt = "Site is disable"
	//		}
	//	}
	//	m.Usex.Shop = shop
	//}

	return rt
}

func (m *MainRPC) ReturnRespone(rs models.RequestResult) *pb.RPCResponse {
	//convert RequestResult into json
	b, err := json.Marshal(rs)
	c3mcommon.CheckError("ReturnRespone Parse JSON in "+m.resp.RPCName, err)
	m.resp.Data = string(b)
	m.resp.Query = int32(m.Rpch.QueryCount + m.QueryCount)
	m.resp.Time = time.Since(m.start).String()
	log.Debugf("query count :%d", m.resp.Query)
	log.Debugf("query time :%s", m.resp.Time)

	return m.resp
}

func (m *MainRPC) ReturnNilRespone() *pb.RPCResponse {
	log.Debugf("response Nil:%+v", m.resp.Data)
	m.resp.Query = int32(m.Rpch.QueryCount + m.QueryCount)
	m.resp.Time = time.Since(m.start).String()
	return m.resp
}
