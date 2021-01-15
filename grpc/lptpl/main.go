package main

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"github.com/tidusant/c3m/common/mycrypto"
	maingrpc "github.com/tidusant/c3m/grpc"
	pb "github.com/tidusant/c3m/grpc/protoc"
	"io/ioutil"
	"strings"

	"context"

	"google.golang.org/grpc"
	"os"

	"github.com/tidusant/c3m/common/log"

	"github.com/tidusant/c3m/repo/models"

	//	"c3m/common/inflect"
	//	"c3m/log"
	"encoding/json"

	"fmt"
	"net"
)

const (
	name           string = "lptemplate"
	ver            string = "1"
	templateFolder        = "templates"
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
			rs = m.Submit(false)
		} else if m.Usex.Action == "rs" {
			rs = m.Submit(true)
		} else if m.Usex.Action == "lat" {
			rs = m.LoadAllTest()
		} else if m.Usex.Action == "la" {
			rs = m.LoadAll()
		} else {
			//unknow action
			return m.ReturnNilRespone(), nil
		}
	}
	return m.ReturnRespone(rs), nil
}
func (m *myRPC) LoadAllTest() models.RequestResult {
	if ok, _ := m.Usex.Modules["c3m-lptpl-admin"]; !ok {
		return models.RequestResult{Error: "permission denied"}
	}
	templates, err := m.Rpch.GetAllLpTemplate(m.Usex.UserID, true)
	if err != nil {
		return models.RequestResult{Error: err.Error()}
	}
	var ptemplates []models.LPTemplate
	var atemplates []models.LPTemplate
	for _, v := range templates {
		if v.Status == 2 {
			ptemplates = append(ptemplates, v)
		} else if v.Status == 1 {
			atemplates = append(atemplates, v)
		}
	}
	//b, _ := json.Marshal(ptemplates)
	type RT struct {
		Ptemplates []models.LPTemplate
		Atemplates []models.LPTemplate
	}

	b, err := json.Marshal(RT{Ptemplates: ptemplates, Atemplates: atemplates})
	if err != nil {
		log.Debugf("error:%s", err.Error())
	}
	return models.RequestResult{Status: 1, Data: string(b)}
}
func (m *myRPC) LoadAll() models.RequestResult {
	if ok, _ := m.Usex.Modules["c3m-lptpl-builder"]; !ok {
		return models.RequestResult{Error: "permission denied"}
	}
	templates, err := m.Rpch.GetAllLpTemplate(m.Usex.UserID, false)
	if err != nil {
		return models.RequestResult{Error: err.Error()}
	}
	b, _ := json.Marshal(templates)
	return models.RequestResult{Status: 1, Data: string(b)}
}
func (m *myRPC) Submit(resubmit bool) models.RequestResult {
	if ok, _ := m.Usex.Modules["c3m-lptpl-builder"]; !ok {
		return models.RequestResult{Error: "permission denied"}
	}
	args := strings.Split(m.Usex.Params, ",")
	if len(args) < 2 {
		return models.RequestResult{Error: "something wrong"}
	}
	tplname := args[0]
	b64str := args[1]
	gzipbyte, err := base64.StdEncoding.DecodeString(b64str)
	if err != nil {
		return models.RequestResult{Error: err.Error()}
	}
	//unzip
	var bb bytes.Buffer
	bb.Write(gzipbyte)

	r, err := gzip.NewReader(&bb)
	if err != nil {
		return models.RequestResult{Error: err.Error()}
	}
	r.Close()
	s, err := ioutil.ReadAll(r)
	if err != nil {
		return models.RequestResult{Error: err.Error()}
	}
	mfile := make(map[string][]byte)
	json.Unmarshal(s, &mfile)
	tplpath := mycrypto.EncodeA(tplname + "_" + m.Usex.Username + mycrypto.StringRand(5))
	//delete old content if resubmit
	var oldtpl models.LPTemplate
	if resubmit {
		oldtpl, err = m.Rpch.GetLpTemplate(m.Usex.UserID, tplname)
		if err != nil {
			return models.RequestResult{Error: err.Error()}
		}
		tplpath = oldtpl.Path
		os.RemoveAll(templateFolder + "/" + tplpath)
	}

	os.Mkdir(templateFolder+"/"+tplpath, 0775)
	for k, v := range mfile {
		//check file folder
		if strings.Index(k, "/") > 0 {
			fpath := k[0:strings.LastIndex(k, "/")]
			os.MkdirAll(templateFolder+"/"+tplpath+"/"+fpath, 0775)
		}
		err := ioutil.WriteFile(templateFolder+"/"+tplpath+"/"+k, v, 0644)
		if err != nil {
			return models.RequestResult{Error: err.Error()}
		}
	}

	//update database

	if !resubmit {
		err := m.Rpch.CreateLpTemplate(m.Usex.UserID, tplname, tplpath)
		if err != nil {
			os.RemoveAll(templateFolder + "/" + tplpath)
			return models.RequestResult{Error: err.Error()}
		}
	} else {

		//reset to waiting approve
		oldtpl.Status = 2
		err = m.Rpch.UpdateLpTemplate(oldtpl)
		if err != nil {
			return models.RequestResult{Error: err.Error()}
		}
	}

	return models.RequestResult{Status: 1, Data: ""}

}

// func (m *myRPC) Remove(usex models.UserSession) string {
// 	log.Debugf("remove  %s", m.Usex.Params)
// 	args := strings.Split(m.Usex.Params, ",")
// 	if len(args) < 2 {
// 		return c3mcommon.ReturnJsonMessage("0", "error submit fields", "", "")
// 	}
// 	log.Debugf("save prod %s", args)
// 	code := args[0]
// 	lang := args[1]
// 	itemremove := m.Rpch.GetPageByCode(m.Usex.UserID, m.Usex.ShopID, code)
// 	if itemremove.Langs[lang] != nil {
// 		//remove slug
// 		m.Rpch.RemoveSlug(itemremove.Langs[lang].Slug, m.Usex.ShopID)
// 		delete(itemremove.Langs, lang)
// 		m.Rpch.SavePage(itemremove)
// 	}

// 	//build home
// 	var bs models.BuildScript
// 	shop := m.Rpch.GetShopById(m.Usex.UserID, m.Usex.ShopID)
// 	bs.ShopID = m.Usex.ShopID
// 	bs.TemplateCode = shop.Theme
// 	bs.Domain = shop.Domain
// 	bs.ObjectID = "home"
// 	rpb.CreateBuild(bs)

// 	//build cat
// 	bs.Collection = "page"
// 	bs.ObjectID = itemremove.Code
// 	rpb.CreateBuild(bs)
// 	return c3mcommon.ReturnJsonMessage("1", "", "success", "")

// }

func main() {
	//default port for service
	var port string
	port = os.Getenv("PORT")
	if port == "" {
		port = "8905"
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
