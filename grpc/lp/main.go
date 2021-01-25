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
		} else if m.Usex.Action == "la" {
			rs = m.SFLoadAll()
		} else if m.Usex.Action == "llc" {
			rs = m.LoadLPContent()
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

	tplID, err := primitive.ObjectIDFromHex(args[3])
	if err != nil {
		return models.RequestResult{Error: "template id is invalid"}
	}
	orguserID := args[4]
	orgID := args[2]
	if orguserID == "" || orgID == "" || campID == "" || content == "" {
		return models.RequestResult{Error: "params is invalid"}
	}
	lp := m.Rpch.GetLPByCampID(campID, orgID)
	if lp.ID.IsZero() {
		lp.UserID = m.Usex.UserID
		lp.Created = time.Now()
		lp.Modified = time.Now()
		lp.LPTemplateID = tplID
		lp.OrgID = orgID
		lp.CampaignID = campID
		lp.SFUserID = orguserID
		lp.Content = content
	}
	if !m.Rpch.SaveLP(lp) {
		return models.RequestResult{Error: "cannot save template"}
	}
	return models.RequestResult{Status: 1, Data: ""}
}

//load all landingpage for user
func (m *myRPC) SFLoadAll() models.RequestResult {
	if ok, _ := m.Usex.Modules["c3m-lptpl-user"]; !ok {
		return models.RequestResult{Error: "permission denied"}
	}
	lps := m.Rpch.GetAllLP(m.Usex.UserID)
	rt := `{"":{}`
	if len(lps) > 0 {
		for _, v := range lps {
			rt += fmt.Sprintf(`,"%s":{"CustomDomain":%t,
"DomainName":"%s",
"FTPUser":"%s",
"FTPPass":"%s",
"LPTemplateID":"%s",
"Created":"%s",
"LastBuild":"%s",
"Modified":"%s",
"Submitted":%d,
"Viewed":%d}`, v.CampaignID, v.CustomDomain, v.DomainName, v.FTPUser, v.FTPPass, v.LPTemplateID.Hex(), v.Created, v.LastBuild, v.Modified, v.Submitted, v.Viewed)
		}
	}
	rt += `}`
	return models.RequestResult{Status: 1, Data: rt}
}

func (m *myRPC) LoadLPContent() models.RequestResult {
	if ok, _ := m.Usex.Modules["c3m-lptpl-user"]; !ok {
		return models.RequestResult{Error: "permission denied"}
	}
	args := strings.Split(m.Usex.Params, ",")

	if len(args) < 2 {
		return models.RequestResult{Error: "not enough params"}
	}
	campID := args[0]
	orgID := args[1]
	if orgID == "" || campID == "" {
		return models.RequestResult{Error: "params is invalid"}
	}
	lp := m.Rpch.GetLPByCampID(campID, orgID)
	if lp.ID.IsZero() {
		return models.RequestResult{Error: "cannot found landing page"}
	}

	return models.RequestResult{Status: 1, Data: lp.Content}
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
