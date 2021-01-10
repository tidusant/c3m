package main

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
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
	templateFolder        = "lptemplates"
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
		} else if m.Usex.Action == "la" {
			rs = m.LoadAll()
		} else {
			//unknow action
			return m.ReturnNilRespone(), nil
		}
	}
	return m.ReturnRespone(rs), nil
}
func (m *myRPC) LoadAll() models.RequestResult {
	templates, err := m.Rpch.GetAllLpTemplate(m.Usex.UserID)
	if err != nil {
		return models.RequestResult{Error: err.Error()}
	}
	b, _ := json.Marshal(templates)
	return models.RequestResult{Status: 1, Data: string(b)}
}
func (m *myRPC) Submit(resubmit bool) models.RequestResult {
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
	tplUserFolder := templateFolder + "/" + tplname + "_" + m.Usex.UserID.Hex()
	//delete old content
	os.RemoveAll(tplUserFolder)
	os.Mkdir(tplUserFolder, 0775)
	for k, v := range mfile {
		//check file folder
		if strings.Index(k, "/") > 0 {
			fpath := k[0:strings.LastIndex(k, "/")]
			os.MkdirAll(tplUserFolder+"/"+fpath, 0775)
		}
		err := ioutil.WriteFile(tplUserFolder+"/"+k, v, 0644)
		if err != nil {
			return models.RequestResult{Error: err.Error()}
		}
	}

	//update database

	if !resubmit {
		err := m.Rpch.CreateLpTemplate(m.Usex.UserID, tplname)
		if err != nil {
			os.RemoveAll(tplUserFolder)
			return models.RequestResult{Error: err.Error()}
		}
	} else {
		//get template
		tpl, err := m.Rpch.GetLpTemplate(m.Usex.UserID, tplname)
		if err != nil {
			return models.RequestResult{Error: err.Error()}
		}
		//reset to waiting approve
		tpl.Status = 2
		err = m.Rpch.UpdateLpTemplate(tpl)
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
