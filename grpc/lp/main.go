package main

import (
	"bufio"
	"context"
	"encoding/json"
	"github.com/jlaffaye/ftp"
	"github.com/tidusant/c3m/common/c3mcommon"
	"github.com/tidusant/c3m/common/mycrypto"
	maingrpc "github.com/tidusant/c3m/grpc"
	pb "github.com/tidusant/c3m/grpc/protoc"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io/ioutil"
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

var (
	LPminserver string
)

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
		} else if m.Usex.Action == "p" {
			rs = m.Publish()
		} else if m.Usex.Action == "sc" {
			rs = m.SaveConfig()
		} else if m.Usex.Action == "d" {
			rs = m.Delete()
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
	lp := m.Rpch.GetLPByCampID(campID, m.Usex.UserID)
	tpl, err := m.Rpch.GetLpTemplateById(tplID)
	if err != nil {
		return models.RequestResult{Error: err.Error()}
	}
	if lp.ID.IsZero() {
		lp.UserID = m.Usex.UserID
		lp.Created = time.Now()
		lp.OrgID = orgID
		lp.CampaignID = campID
		lp.Path = mycrypto.StringRand(5) + mycrypto.StringRand(5)
		lp.Path = mycrypto.Encode4(tpl.Path + "/publish/" + lp.Path)
	} else {
		//if change tpl then remove publish folder
		if lp.LPTemplateID != tplID {
			publishFolder := "./templates/" + mycrypto.Decode4(lp.Path)
			os.RemoveAll(publishFolder)
		}
	}
	lp.Modified = time.Now()
	lp.LPTemplateID = tplID
	lp.SFUserID = orguserID
	lp.Content = content
	if !m.Rpch.SaveLP(lp) {
		return models.RequestResult{Error: "cannot save landing page"}
	}
	return models.RequestResult{Status: 1, Data: ""}
}

//load all landingpage for user
func (m *myRPC) SFLoadAll() models.RequestResult {
	if ok, _ := m.Usex.Modules["c3m-lptpl-user"]; !ok {
		return models.RequestResult{Error: "permission denied"}
	}
	campIdsWillRemove := strings.Split(m.Usex.Params, "-")
	//get all LP
	lps := m.Rpch.GetAllLP(m.Usex.UserID)
	rt := `{"":{}`
	if len(lps) > 0 {
		for _, v := range lps {
			rt += fmt.Sprintf(`,"%s":{"CustomHost":%t,
"DomainName":"%s",
"FTPUser":"%s",
"FTPPass":"%s",
"LPTemplateID":"%s",
"Created":"%s",
"LastBuild":"%s",
"Modified":"%s",
"Submitted":%d,
"Viewed":%d}`, v.CampaignID, v.CustomHost, v.DomainName, v.FTPUser, v.FTPPass, v.LPTemplateID.Hex(), v.Created, v.LastBuild, v.Modified, v.Submitted, v.Viewed)

			//remove campid in campIdsWillRemove list
			if ok, i := c3mcommon.InArray(v.CampaignID, campIdsWillRemove); ok {
				campIdsWillRemove = append(campIdsWillRemove[:i], campIdsWillRemove[i+1:]...)
			}
		}
	}
	rt += `}`
	//remove deleted landingpage's campaignid:
	m.Usex.Params = strings.Join(campIdsWillRemove, ",")
	m.Delete()

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
	lp := m.Rpch.GetLPByCampID(campID, m.Usex.UserID)
	if lp.ID.IsZero() {
		return models.RequestResult{Error: "cannot found landing page"}
	}

	return models.RequestResult{Status: 1, Data: lp.Content}
}

func (m *myRPC) SaveConfig() models.RequestResult {
	if ok, _ := m.Usex.Modules["c3m-lptpl-user"]; !ok {
		return models.RequestResult{Error: "permission denied"}
	}

	data, err := mycrypto.DecompressFromBase64(m.Usex.Params)
	if err != nil {
		return models.RequestResult{Error: err.Error()}
	}
	var lp models.LandingPage
	err = json.Unmarshal([]byte(data), &lp)
	if err != nil {
		return models.RequestResult{Error: "cannot parse json data to landing page object"}
	}
	oldlp := m.Rpch.GetLPByCampID(lp.CampaignID, m.Usex.UserID)
	if oldlp.ID.IsZero() {
		return models.RequestResult{Error: "cannot found landing page"}
	}
	//check domain name
	if !lp.CustomHost {
		if strings.Trim(lp.DomainName, " ") == "" {
			return models.RequestResult{Error: "Domain Name is empty"}
		}
		f := func(r rune) bool {
			return r < '0' || r > 'z'
		}
		if strings.IndexFunc(lp.DomainName, f) != -1 {
			return models.RequestResult{Error: "Found special char in Domain Name"}
		}
		if len(strings.Trim(lp.DomainName, " ")) < 4 {
			return models.RequestResult{Error: "Domain Name length must greater than 3"}
		}
	} else {
		if strings.Trim(lp.DomainName, " ") == "" {
			return models.RequestResult{Error: "Domain Name is empty"}
		}
		if strings.Trim(lp.FTPUser, " ") == "" {
			return models.RequestResult{Error: "FTPUser is empty"}
		}
		if strings.Trim(lp.FTPPass, " ") == "" {
			return models.RequestResult{Error: "FTPPass is empty"}
		}
	}

	oldlp.CustomHost = lp.CustomHost
	oldlp.DomainName = lp.DomainName
	oldlp.FTPPass = lp.FTPPass
	oldlp.FTPUser = lp.FTPUser
	if !m.Rpch.SaveLP(oldlp) {
		return models.RequestResult{Error: "cannot save config"}
	}
	return models.RequestResult{Status: 1, Data: ""}
}

func (m *myRPC) Delete() models.RequestResult {
	if ok, _ := m.Usex.Modules["c3m-lptpl-user"]; !ok {
		return models.RequestResult{Error: "permission denied"}
	}
	campidsstr := strings.Split(m.Usex.Params, ",")

	//detele
	camps := m.Rpch.GetAllLpByCampIds(campidsstr, m.Usex.UserID)
	for _, v := range camps {
		os.RemoveAll("./templates/" + mycrypto.Decode4(v.Path))
	}
	if !m.Rpch.DeleteLP(campidsstr, m.Usex.UserID) {
		return models.RequestResult{Error: "Could not delele landing page of campaign ID: " + m.Usex.Params}
	}
	return models.RequestResult{Status: 1, Data: ""}
}

func (m *myRPC) Publish() models.RequestResult {
	if ok, _ := m.Usex.Modules["c3m-lptpl-user"]; !ok {
		return models.RequestResult{Error: "permission denied"}
	}

	lp := m.Rpch.GetLPByCampID(m.Usex.Params, m.Usex.UserID)
	if lp.ID.IsZero() {
		return models.RequestResult{Error: "Landing page not found"}
	}

	//build content for publish
	lppath := mycrypto.Decode4(lp.Path)
	publishFolder := "./templates/" + lppath
	os.RemoveAll(publishFolder)
	err := os.MkdirAll(publishFolder, 0775)
	if err != nil {
		return models.RequestResult{Error: err.Error()}
	}
	//create content
	content, err := mycrypto.DecompressFromBase64(lp.Content)
	if err != nil {
		return models.RequestResult{Error: err.Error()}
	}
	err = ioutil.WriteFile(publishFolder+"/content.html", []byte(content), 0644)
	if err != nil {
		return models.RequestResult{Error: err.Error()}
	}
	//call service publish
	argspath := strings.Split(lppath, "/")
	if len(argspath) < 3 {
		return models.RequestResult{Error: "landing page path invalid"}
	}
	bodystr := c3mcommon.RequestAPI2(LPminserver+"/publish", argspath[0], m.Usex.Session+","+argspath[2])
	log.Debug(bodystr)
	var rs models.RequestResult
	err = json.Unmarshal([]byte(bodystr), &rs)

	if err != nil {
		return models.RequestResult{Error: err.Error()}
	}
	if rs.Status != 1 {
		return models.RequestResult{Error: rs.Error}
	}
	if lp.DomainName == "" {
		return models.RequestResult{Error: "Domain Name is empty"}
	}

	if lp.CustomHost {
		//connect ftp
		ftpclient, err := ftp.Dial(lp.DomainName)
		if err == nil {
			err = ftpclient.Login(lp.FTPUser, lp.FTPPass)
		}
		// perform copy
		if err == nil {
			file, err := os.Open(publishFolder + "/index.html")
			if c3mcommon.CheckError("Read file index.html", err) {
				err = ftpclient.Stor("./index.html", bufio.NewReader(file))
				file.Close()
				if err == nil {
					file, err = os.Open(publishFolder + "/style.css")
					if c3mcommon.CheckError("Read file style.css", err) {
						err = ftpclient.Stor("./style.css", bufio.NewReader(file))
						file.Close()
					}
				}
			}
		}
		ftpclient.Quit()
	} else {
		//using .c3m.site domain
	}

	//update lp last publish
	lp.LastBuild = time.Now()
	if !m.Rpch.SaveLP(lp) {
		return models.RequestResult{Error: "Can not update Last Build time after publish"}
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
	LPminserver = os.Getenv("LPMIN_ADD")
	s := grpc.NewServer()
	fmt.Printf("listening on %s\n", port)
	pb.RegisterGRPCServicesServer(s, &service{})

	if err := s.Serve(lis); err != nil {
		log.Errorf("failed to serve : %v", err)
	}
	fmt.Print("exit")

}
