package main

import (
	"context"
	maingrpc "github.com/tidusant/c3m/grpc"
	pb "github.com/tidusant/c3m/grpc/protoc"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
	"time"

	"google.golang.org/grpc"
	"os"

	"github.com/tidusant/c3m/common/log"

	"github.com/tidusant/c3m/repo/models"

	"fmt"
	"net"
)

const (
	name string = "landingpage"
	ver  string = "1"
)

type service struct {
	pb.UnimplementedGRPCServicesServer
}

//extend class MainRPC
type myRPC struct {
	maingrpc.MainRPC
}

func (s *service) Call(ctx context.Context, in *pb.RPCRequest) (*pb.RPCResponse, error) {
	m := myRPC{}
	//generate user information into usex by calling parent func (m *myRPC) InitUsex that return error string
	rs := models.RequestResult{Error: m.InitUsex(ctx, in, name, ver)}
	//if not error then continue call func
	if rs.Error == "" {
		if m.Usex.Action == "s" {
			rs = m.Save()
			//} else if m.Usex.Action == "l" {
			//	rs = m.Load()
			//} else if m.Usex.Action == "p" {
			//	rs = m.Publish()
			//} else if m.Usex.Action == "sc" {
			//	rs = m.SaveConfig()
			//} else if m.Usex.Action == "lc" {
			//	rs = m.LoadConfig()
			//} else if m.Usex.Action == "d" {
			//	rs = m.Delete()
		} else {
			//unknow action
			return m.ReturnNilRespone(), nil
		}
	}
	return m.ReturnRespone(rs), nil
}

//load all template for test and approve
func (m *myRPC) Save() models.RequestResult {
	if ok, _ := m.Usex.Modules["c3m-lptpl-user"]; !ok {
		return models.RequestResult{Error: "permission denied"}
	}
	args := strings.Split(m.Usex.Params, ",")
	if len(args) < 5 {
		return models.RequestResult{Error: "not enough params"}
	}
	campID := args[1]
	content := args[0]
	tplID, err := primitive.ObjectIDFromHex(args[2])
	if err != nil {
		return models.RequestResult{Error: err.Error()}
	}
	orguserID := args[3]
	orgID := args[4]
	lp := m.Rpch.GetLPByCampID(campID, orgID)
	if lp.ID.IsZero() {
		lp.Created = time.Now()
		lp.Modified = time.Now()
		lp.LPTemplateID = tplID
		lp.OrgID = orgID
		lp.SFUserID = orguserID
		lp.Content = content
	}
	if !m.Rpch.SaveLP(lp) {
		return models.RequestResult{Error: "cannot save template"}
	}
	return models.RequestResult{Status: 1, Data: ""}
}

func main() {
	//default port for service
	var port string
	port = os.Getenv("PORT")
	if port == "" {
		port = "8906"
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
