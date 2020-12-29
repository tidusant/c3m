package main

import (
	maingrpc "github.com/tidusant/c3m/grpc"
	pb "github.com/tidusant/c3m/grpc/protoc"

	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	name string = "page"
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
			rs = m.LoadAllPage()
		} else if m.Usex.Action == "l" {
			rs = m.LoadPage()
		} else if m.Usex.Action == "la" {
			rs = m.LoadAllPage()
		} else {
			//unknow action
			return m.ReturnNilRespone(), nil
		}
	}
	return m.ReturnRespone(rs), nil
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
func (m *myRPC) LoadPage() models.RequestResult {
	pageid, _ := primitive.ObjectIDFromHex(m.Usex.Params)
	if pageid == primitive.NilObjectID {
		return models.RequestResult{Error: "No page found"}
	}

	item := m.Rpch.GetPageById(pageid)
	b, _ := json.Marshal(item)

	return models.RequestResult{Status: 1, Data: string(b)}

}
func (m *myRPC) LoadAllPage() models.RequestResult {
	shop := m.Rpch.GetShopById(m.Usex.UserID, m.Usex.ShopID)

	items := m.Rpch.GetAllPage(m.Usex.ShopID, shop.TemplateID)
	if len(items) == 0 {
		return models.RequestResult{Error: "No page found"}
	}
	b, _ := json.Marshal(items)
	return models.RequestResult{Status: 1, Data: string(b)}

}

//func (m *myRPC) SavePage() models.RequestResult {
//	var newitem models.Page
//	log.Debugf("Unmarshal %s", m.Usex.Params)
//	err := json.Unmarshal([]byte(m.Usex.Params), &newitem)
//	if !c3mcommon.CheckError("json parse page", err) {
//		return c3mcommon.ReturnJsonMessage("0", "json parse fail", "", "")
//	}
//
//	//update
//	//check olditem
//	log.Debugf("json parse: %v", newitem)
//	olditem := m.Rpch.GetPageByCode(m.Usex.Shop.Theme, m.Usex.Shop.ID.Hex(), newitem.Code)
//	if olditem.Code == "" {
//		return c3mcommon.ReturnJsonMessage("0", "item not found", "", "")
//	}
//
//	var langlinks []models.LangLink
//	for langcode, _ := range newitem.Langs {
//		if newitem.Langs[langcode].Title == "" {
//			continue
//		}
//		var newslug models.Slug
//		newslug.ShopId = m.Usex.Shop.ID.Hex()
//		newslug.Object = "page"
//		newslug.ObjectId = olditem.ID.Hex()
//		newslug.Lang = langcode
//		newslug.TemplateCode = m.Usex.Shop.Theme
//
//		//newslug
//		pagelang := newitem.Langs[langcode]
//		log.Debugf("pagelang %v", pagelang)
//		newslug.Slug = inflect.ParameterizeJoin(newitem.Langs[langcode].Title, "_")
//		//check slug
//		if newitem.Code == "home" && m.Usex.Shop.Config.DefaultLang == langcode {
//			newslug.Slug = ""
//		}
//
//		pagelang.Slug = m.Rpch.SaveSlug(newslug)
//		newitem.Langs[langcode] = pagelang
//
//		if newitem.Langs[langcode].Slug != "" {
//			langlinks = append(langlinks, models.LangLink{Href: newitem.Langs[langcode].Slug + "/", Code: langcode, Name: c3mcommon.GetLangnameByCode(langcode)})
//		} else {
//			langlinks = append(langlinks, models.LangLink{Href: newitem.Langs[langcode].Slug, Code: langcode, Name: c3mcommon.GetLangnameByCode(langcode)})
//		}
//		//=====
//
//	}
//
//	//update
//	olditem.Seo = newitem.Seo
//	olditem.Langs = newitem.Langs
//	olditem.Blocks = newitem.Blocks
//	olditem.LangLinks = langlinks
//	strrt := m.Rpch.SavePage(olditem)
//	if strrt == "0" {
//		return c3mcommon.ReturnJsonMessage("0", "error", "error", "")
//	}
//
//	//rebuild page
//	b, err := json.Marshal(olditem)
//
//	//create build
//
//	errstr := m.Rpch.CreateBuild("page", olditem.ID.Hex(), string(b), usex)
//	if errstr != "" {
//		return c3mcommon.ReturnJsonMessage("0", errstr, "build error", "")
//	}
//	errstr = m.Rpch.CreateCommonDataBuild(usex)
//	if errstr != "" {
//		return c3mcommon.ReturnJsonMessage("0", errstr, "build error", "")
//	}
//
//	// //rebuild home
//	// var bs models.BuildScript
//	// shop := m.Rpch.GetShopById(m.Usex.UserID, m.Usex.Shop.ID.Hex())
//	// bs.ShopID = m.Usex.Shop.ID.Hex()
//	// bs.TemplateCode = shop.Theme
//	// bs.Domain = shop.Domain
//	//  bs.ObjectID = "home"
//	//  rpb.CreateBuild(bs)
//
//	// //rebuild cat
//	// bs.Collection = "page"
//	// bs.ObjectID = newitem.Code
//	// rpb.CreateBuild(bs)
//	return c3mcommon.ReturnJsonMessage("1", "", "success", "")
//}
func main() {
	//default port for service
	var port string
	port = os.Getenv("PORT")
	if port == "" {
		port = "8904"
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
