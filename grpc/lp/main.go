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
	"regexp"
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
			rt = m.ReturnNilRespone()
			return
		}
	}
	rt = m.ReturnRespone(rs)
	return
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
"FTPHost":"%s",
"FTPUser":"%s",
"LPTemplateID":"%s",
"Created":"%s",
"LastBuild":"%s",
"Modified":"%s",
"Submitted":%d,
"Viewed":%d}`, v.CampaignID, v.CustomHost, v.DomainName, v.FTPHost, v.FTPUser, v.LPTemplateID.Hex(), v.Created, v.LastBuild, v.Modified, v.Submitted, v.Viewed)

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
			return (r < '0' || r > 'z') && r != '-' && r != '_'
		}
		if strings.IndexFunc(lp.DomainName, f) != -1 {
			return models.RequestResult{Error: "Found special char in URL"}
		}
		if len(strings.Trim(lp.DomainName, " ")) < 4 {
			return models.RequestResult{Error: "URL length must greater than 3"}
		}
	} else {
		if strings.Trim(lp.DomainName, " ") == "" {
			return models.RequestResult{Error: "URL is empty"}
		}
		regex := `(http(s)?):\/\/[(www\.)?a-zA-Z0-9@:%._\+~#=]{2,256}\.[a-z]{2,6}\b([-a-zA-Z0-9@:%_\+.~#?&//=]*)`
		if ok, _ := regexp.MatchString(regex, lp.DomainName); !ok {
			return models.RequestResult{Error: "URL is invalid"}
		}
		if strings.Trim(lp.FTPHost, " ") == "" {
			return models.RequestResult{Error: "FTPHost is empty"}
		}
		if strings.Trim(lp.FTPUser, " ") == "" {
			return models.RequestResult{Error: "FTPUser is empty"}
		}
		if strings.Trim(lp.FTPPass, " ") == "" {
			return models.RequestResult{Error: "FTPPass is empty"}
		}

		//check duplicate url
		if m.Rpch.LPCheckDomainUrlExist(lp.DomainName) {
			return models.RequestResult{Error: "URL is exist. Please choose another one."}
		}

	}

	oldlp.CustomHost = lp.CustomHost
	oldlp.DomainName = lp.DomainName
	oldlp.FTPHost = lp.FTPHost
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
	if lp.DomainName == "" {
		return models.RequestResult{Error: "Domain Name is empty"}
	}

	//build content for publish
	lppath := mycrypto.Decode4(lp.Path)
	argspath := strings.Split(lppath, "/")
	if len(argspath) < 3 {
		return models.RequestResult{Error: "landing page path invalid"}
	}

	publishFolder := "./templates/" + lppath
	os.RemoveAll(publishFolder)
	err := os.MkdirAll(publishFolder, 0755)
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

	bodystr := c3mcommon.RequestAPI2(LPminserver+"/publish", argspath[0], m.Usex.Session+","+argspath[2])
	log.Debugf("bodystr rt:%s", bodystr)
	var rs models.RequestResult
	err = json.Unmarshal([]byte(bodystr), &rs)

	if err != nil {
		return models.RequestResult{Error: err.Error()}
	}
	if rs.Status != 1 {
		return models.RequestResult{Error: rs.Error}
	}
	errStr := ""
	if lp.CustomHost {
		//connect ftp
		var ftpclient *ftp.ServerConn
		log.Debugf("connecting to %s", lp.FTPHost)
		host := strings.Replace(lp.FTPHost, `http://`, "", -1)
		host = strings.Replace(host, `https://`, "", -1)
		host = strings.Replace(host, `ftp://`, "", -1)
		hostsplit := strings.Split(host, ":")
		port := "21"
		if len(hostsplit) > 1 {
			port = hostsplit[len(hostsplit)-1]
			host = hostsplit[len(hostsplit)-2]
		}

		ftpclient, err := ftp.Dial(host + ":" + port)
		if err != nil {
			return models.RequestResult{Error: err.Error()}
		}
		if ftpclient == nil {
			errStr = fmt.Sprintf("Cannot connect to %s" + lp.FTPHost)
		}

		err = ftpclient.Login(lp.FTPUser, lp.FTPPass)

		// perform copy
		if err == nil {
			file, err := os.Open(publishFolder + "/index.html")
			if c3mcommon.CheckError("Read file index.html", err) {
				err = ftpclient.Stor("./index.html", bufio.NewReader(file))
				c3mcommon.CheckError("Store file index.html to host "+lp.FTPHost, err)
				file.Close()
				if err == nil {
					file, err = os.Open(publishFolder + "/style.css")
					if c3mcommon.CheckError("Read file style.css", err) {
						err = ftpclient.Stor("./style.css", bufio.NewReader(file))
						c3mcommon.CheckError("Store file style.css to host "+lp.FTPHost, err)
						file.Close()
						if err != nil {
							errStr = fmt.Sprintf("Error when storing file to host: Cannot write file to host %s", lp.FTPHost)
						}
					} else {
						errStr = fmt.Sprintf("Error when storing file to host: Cannot read file style.css of landing page")
					}
				} else {
					errStr = fmt.Sprintf("Error when storing file to host: Cannot write file to host %s", lp.FTPHost)
				}
			} else {
				errStr = fmt.Sprintf("Error when storing file to host: Cannot read file index.html of landing page")
			}
		} else {
			errStr = fmt.Sprintf("Cannot logged in to %s" + lp.FTPHost)
		}
		ftpclient.Quit()
	} else {
		//using .c3m.site domain
		log.Debugf("using c3m site")
		lpPath := "./lp/" + lp.DomainName
		os.MkdirAll(lpPath, 0755)
		fileb, err := ioutil.ReadFile(publishFolder + "/index.html")
		if err != nil {
			return models.RequestResult{Error: "error reading index.html"}
		}
		//copy file
		err = ioutil.WriteFile(lpPath+"/index.html", fileb, 0644)
		if err != nil {
			return models.RequestResult{Error: "error creating index.html"}
		}
		fileb, err = ioutil.ReadFile(publishFolder + "/style.css")
		if err != nil {
			return models.RequestResult{Error: "error reading style.css"}
		}
		//copy file
		err = ioutil.WriteFile(lpPath+"/style.css", fileb, 0644)
		if err != nil {
			return models.RequestResult{Error: "error creating style.css"}
		}
	}

	if errStr != "" {
		return models.RequestResult{Error: errStr}
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
