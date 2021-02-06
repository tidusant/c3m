package main

import (
	"context"
	"encoding/json"
	"github.com/tidusant/c3m/common/mycrypto"
	maingrpc "github.com/tidusant/c3m/grpc"
	pb "github.com/tidusant/c3m/grpc/protoc"
	"io/ioutil"
	"time"

	"google.golang.org/grpc"
	"os"

	"github.com/tidusant/c3m/common/log"

	"github.com/tidusant/c3m/repo/models"

	"fmt"
	"net"
)

const (
	name string = "lplead"
	ver  string = "1"
)

type service struct {
	pb.UnimplementedGRPCServicesServer
}

//extend class MainRPC
type myRPC struct {
	maingrpc.MainRPC
}

var (
	LPminserver string
)

func (s *service) Call(ctx context.Context, in *pb.RPCRequest) (rt *pb.RPCResponse, err error) {
	m := myRPC{}
	rs := models.RequestResult{Error: m.InitUsex(ctx, in, name, ver)}
	err = nil
	defer func() {
		if err := recover(); err != nil {
			ioutil.WriteFile("templates/"+name+".panic.log", []byte(fmt.Sprint(time.Now().Format("2006-01-02 15:04:05")+" >> panic occurred:", err)), 0644)
			rs.Error = "Something wrong"
			rt = m.ReturnRespone(rs)
		}

	}()
	//generate user information into usex by calling parent func (m *myRPC) InitUsex that return error string

	//if not error then continue call func
	if rs.Error == "" {
		if m.Usex.Action == "s" {
			rs = m.SubmitLead()
		} else if m.Usex.Action == "laus" {
			rs = m.LoadAllUnSync()
		} else {
			//unknow action
			rt = m.ReturnNilRespone()
			return
		}
	}
	rt = m.ReturnRespone(rs)
	return
}

func (m *myRPC) SubmitLead() models.RequestResult {
	args := mycrypto.Base64Decode(m.Usex.Params)
	var lead models.Lead
	log.Debugf(args)
	err := json.Unmarshal([]byte(args), &lead)

	if err != nil {
		log.Debugf("error json parse: " + err.Error())
		return models.RequestResult{Error: "Invalid Params"}
	}
	if !m.Rpch.SaveLead(lead) {
		return models.RequestResult{Error: "Could not save lead"}
	}
	return models.RequestResult{Status: 1, Data: ""}
}

func (m *myRPC) LoadAllUnSync() models.RequestResult {
	var rs []models.Lead
	rs = m.Rpch.GetAllLeadByOrgID(m.Usex.Params)
	rt := `[{}`
	if len(rs) > 0 {
		for _, v := range rs {
			rt += fmt.Sprintf(`,{
"Name":%s,
"Email":"%s",
"Phone":"%s",
"Description":"%s",
"LeadSource":"Web",
"Status":"Open - Not Contacted"}`,
				v.Name, v.Email, v.Phone, v.Message)
		}
	}
	rt += `]`

	return models.RequestResult{Status: 1, Data: rt}
}

func main() {
	//default port for service
	var port string
	port = os.Getenv("PORT")
	if port == "" {
		port = "8907"
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
