package main

import (
	"encoding/json"
	"fmt"
	"github.com/tidusant/c3m/common/log"
	maingrpc "github.com/tidusant/c3m/grpc"
	pb "github.com/tidusant/c3m/grpc/protoc"

	"github.com/tidusant/c3m/repo/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net"

	"context"
	"google.golang.org/grpc"
	"os"
)

const (
	name string = "shop"
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
	//generate user information into usex by calling parent func InitUsex that return error string
	rs := models.RequestResult{Error: m.InitUsex(ctx, in, name, ver)}

	//if not error then continue call func
	if rs.Error == "" {
		if m.Usex.Action == "cs" {
			rs = m.changeShop()
		} else if m.Usex.Action == "lsi" {
			rs = m.loadshopinfo()
		} else if m.Usex.Action == "ca" {
			rs = m.doCreateAlbum()
		} else if m.Usex.Action == "la" {
			rs = m.doLoadalbum()
		} else if m.Usex.Action == "ea" {
			rs = m.doEditAlbum()
		} else if m.Usex.Action == "cga" {
			rs = m.configGetAll()
		} else if m.Usex.Action == "cgs" {
			rs = m.configSave()
		} else if m.Usex.Action == "lims" {
			rs = m.getShopLimits()
		} else {
			//unknow action
			return m.ReturnNilRespone(), nil
		}
	}

	//if there is other repo than rpch then accummulate query count here
	//m.QueryCount+=m.Rpch.QueryCount

	return m.ReturnRespone(rs), nil

}

type ConfigViewData struct {
	ShopConfigs     models.ShopConfigs
	TemplateConfigs []ConfigItem
	BuildConfigs    models.BuildConfig
}
type ConfigItem struct {
	Key   string
	Type  string
	Value string
}

func (m *myRPC) loadshopinfo() models.RequestResult {

	strrt := `{"DefaultShopId":"` + m.Usex.ShopID.Hex() + `","Shops":`

	//maxfileupload
	//strrt += `,"MaxFileUpload":` + strconv.Itoa(cuahang.GetShopLimitbyKey(m.Usex.Shop.ID, "maxfileupload"))
	//strrt += `,"MaxSizeUpload":` + strconv.Itoa(cuahang.GetShopLimitbyKey(m.Usex.Shop.ID, "maxsizeupload"))

	//orther shop
	Shops := m.Rpch.GetShopsByUserId(m.Usex.UserID)

	strrt += `[`
	for _, shop := range Shops {
		strrt += `{"Name":"` + shop.Name + `","ID":"` + shop.ID.Hex() + `"},`
	}
	if len(Shops) > 0 {
		strrt = strrt[:len(strrt)-1] + `]`
	} else {
		strrt += `]`
	}
	strrt += "}"
	return models.RequestResult{Status: 1, Error: "", Message: "", Data: strrt}

}

func (m *myRPC) changeShop() models.RequestResult {
	shopidObj, _ := primitive.ObjectIDFromHex(m.Usex.Params)
	shopchange := m.Rpch.GetShopById(m.Usex.UserID, shopidObj)
	if shopchange.ID == primitive.NilObjectID {
		return models.RequestResult{Error: "Change shop fail"}
	}
	//change shopid
	m.Usex.ShopID = shopchange.ID
	return models.RequestResult{Status: 1, Data: shopchange.ID.Hex()}
}
func (m *myRPC) configSave() models.RequestResult {
	//var config ConfigViewData

	//err := json.Unmarshal([]byte(m.Usex.Params), &config)
	//if !c3mcommon.CheckError("json parse page", err) {
	//	return models.RequestResult{Error: "json parse fail"}
	//}
	//m.Usex.Shop.Config = config.ShopConfigs
	//cuahang.SaveShopConfig(m.Usex.Shop)
	//
	//// //save template config
	//str := `{"Code":"` + m.Usex.Shop.Theme + `","TemplateConfigs":[{}`
	//for _, conf := range config.TemplateConfigs {
	//	str += `,{"Key":"` + conf.Key + `","Value":"` + conf.Value + `"}`
	//}
	//str += `]`
	//b, _ := json.Marshal(config.BuildConfigs)
	//str += `,"BuildConfig":` + string(b) + `}`
	//
	//request := "savetemplateconfig|" + m.Usex.Session
	//resp := c3mcommon.RequestBuildService(request, "POST", str)
	//
	//if resp.Status != 1 {
	//	return resp
	//}
	//
	//// //save build config
	//
	//// var bcf models.BuildConfig
	//// bcf = config.BuildConfigs
	//// bcf.ShopId = m.Usex.Shop.ID
	//// rpb.SaveConfig(bcf)
	////rebuild config
	//cuahang.Rebuild(usex)
	//return c3mcommon.ReturnJsonMessage("1", "", "success", "")
	return models.RequestResult{Status: 1}

}
func (m *myRPC) configGetAll() models.RequestResult {
	//var config ConfigViewData
	//config.ShopConfigs = m.Usex.Shop.Config
	//log.Debugf("configGetAll")
	//request := "gettemplateconfig|" + m.Usex.Session
	//resp := c3mcommon.RequestBuildService(request, "POST", m.Usex.Shop.Theme)
	//log.Debugf("RequestBuildService call done")
	//if resp.Status != 1 {
	//	return resp
	//}
	//var confs struct {
	//	TemplateConfigs []ConfigItem
	//	BuildConfigs    models.BuildConfig
	//}
	//json.Unmarshal([]byte(resp.Data), &confs)
	//
	//config.TemplateConfigs = confs.TemplateConfigs
	//config.BuildConfigs = confs.BuildConfigs
	//config.BuildConfigs.ID = ""
	//config.BuildConfigs.ShopId = ""
	//b, _ := json.Marshal(config)
	//
	//return c3mcommon.ReturnJsonMessage("1", "", "success", string(b))
	return models.RequestResult{Status: 1}

}
func (m *myRPC) getShopLimits() models.RequestResult {

	limits := m.Rpch.GetShopLimits(m.Usex.ShopID)

	b, _ := json.Marshal(limits)
	return models.RequestResult{Status: 1, Error: "", Data: string(b)}

}

// func loadcat(usex models.UserSession) string {
// 	log.Debugf("loadcat begin")
// 	shop := cuahang.GetShopById(m.Usex.UserID, m.Usex.ShopID)

// 	strrt := "["
// 	log.Debugf("load cats:%v", shop.ShopCats)
// 	catinfstr := ""
// 	for _, cat := range shop.ShopCats {
// 		catlangs := ""
// 		for lang, catinf := range cat.Langs {
// 			catlangs += """ + lang + "":{"name":"" + catinf.Slug + "","slug":"" + catinf.Name + "","description":"" + catinf.Description + ""},"
// 		}
// 		catlangs = catlangs[:len(catlangs)-1]
// 		catinfstr += "{"code":"" + cat.Code + "","langs":{" + catlangs + "}},"
// 	}
// 	if catinfstr == "" {
// 		strrt += "{}]"
// 	} else {
// 		strrt += catinfstr[:len(catinfstr)-1] + "]"
// 	}

// 	return c3mcommon.ReturnJsonMessage("1", "", "success", strrt)

// }

func (m *myRPC) doCreateAlbum() models.RequestResult {
	//albumname := m.Usex.Params
	//if albumname == "" {
	//	return models.RequestResult{Error:"album's name empty"}
	//
	//}
	////get config
	//
	//if m.Usex.Shop.ID == "" {
	//	return models.RequestResult{Error:"shop not found"}
	//
	//}
	//
	//// if m.Usex.Shop.Config.Level == 0 {
	//// 	return c3mcommon.ReturnJsonMessage("0", "config error", "", "")
	//
	//// }
	//// if m.Usex.Shop.Config.MaxAlbum <= len(m.Usex.Shop.Albums) {
	//// 	return c3mcommon.ReturnJsonMessage("2", "album count limited", "", "")
	//// }
	//
	//slug := strings.ToLower(mystring.Camelize(mystring.Asciify(albumname)))
	//albumslug := slug
	//if slug == "" {
	//	return models.RequestResult{Error:"album's slug empty"}
	//
	//}
	//
	////save albumname
	//var album models.ShopAlbum
	//album.Slug = albumslug
	//album.Name = albumname
	//album.ShopID = m.Usex.Shop.ID
	//album.UserId = m.Usex.UserID
	//album = cuahang.SaveAlbum(album)
	//b, _ := json.Marshal(album)
	//
	//return c3mcommon.ReturnJsonMessage("1", "", "success", string(b))
	return models.RequestResult{Status: 1}

}
func (m *myRPC) doLoadalbum() models.RequestResult {

	////get albums
	//albums := cuahang.LoadAllShopAlbums(m.Usex.Shop.ID)
	//if len(albums) == 0 {
	//	//create
	//	var album models.ShopAlbum
	//	album.Slug = "default"
	//	album.Name = "Default"
	//	album.ShopID = m.Usex.Shop.ID
	//	album.UserId = m.Usex.UserID
	//	album = cuahang.SaveAlbum(album)
	//	albums = append(albums, album)
	//}
	//
	//b, err := json.Marshal(albums)
	//c3mcommon.CheckError("json parse doLoadalbum", err)
	//return c3mcommon.ReturnJsonMessage("1", "", "", string(b))

	return models.RequestResult{Status: 1}

}
func (m *myRPC) doEditAlbum() models.RequestResult {
	////log.Debugf("update album ")
	//var newitem models.ShopAlbum
	//log.Debugf("Unmarshal %s", m.Usex.Params)
	//err := json.Unmarshal([]byte(m.Usex.Params), &newitem)
	//if !c3mcommon.CheckError("json parse page", err) {
	//	return c3mcommon.ReturnJsonMessage("0", "json parse ShopAlbum fail", "", "")
	//}
	//newitem.ShopID = m.Usex.Shop.ID
	//newitem.UserId = m.Usex.UserID
	//cuahang.SaveAlbum(newitem)
	////log.Debugf("update album false %s", albumname)
	//return c3mcommon.ReturnJsonMessage("0", "album not found", "", "")
	return models.RequestResult{Status: 1}
}
func main() {
	//default port for service
	var port string
	port = os.Getenv("PORT")
	if port == "" {
		port = "8902"
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

//repush
