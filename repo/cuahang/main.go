package cuahang

import (
	"context"
	"github.com/tidusant/c3m/common/c3mcommon"
	"github.com/tidusant/c3m/common/log"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
	"time"
)

var (
	db  *mongo.Database
	ctx context.Context
)

type Repo struct {
	QueryCount int `default:0`
	QueryTime  time.Duration
}

func init() {
	log.Info("init repo cuahang...")
	strErr := ""
	ctx = context.Background()
	context.WithValue(ctx, 43, 44)
	u, ok := ctx.Value(43).(int64)
	log.Debugf("context:%v - %v", u, ok)
	db, strErr = c3mcommon.ConnectAtlasDB(ctx, "chadmin")
	if strErr != "" {
		log.Infof(strErr)
		os.Exit(1)
	}
	log.Info("done")
}

//==============slug=================
//func CreateBuild(object, objectid, data string, usex models.UserSession) string {
//	request := "createbuild|" + usex.Session
//	var bs models.BuildScript
//	bs.Object = object
//	bs.ObjectId = objectid
//	bs.Data = data
//	//bs.ShopConfigs = usex.Shop.Config
//	//bs.TemplateCode = usex.Shop.Theme
//	//bs.ShopId = usex.Shop.ID.Hex()
//	b, _ := json.Marshal(bs)
//	c3mcommon.RequestBuildServiceAsync(request, "POST", string(b))
//
//	return ""
//}

//
//func CreateCommonDataBuild(usex models.UserSession) string {
//
//	var bs models.BuildScript
//	bs.ShopConfigs = usex.Shop.Config
//	bs.TemplateCode = usex.Shop.Theme
//	bs.ShopId = usex.Shop.ID.Hex()
//	bs.Object = "common"
//	var common models.CommonData
//	//get allpage
//	common.Pages = GetAllPage(usex.Shop.Theme, usex.Shop.ID.Hex())
//	//remove page content
//	for i, _ := range common.Pages {
//		common.Pages[i].Blocks = nil
//		common.Pages[i].Seo = ""
//	}
//	//get all product
//
//	//get all productcat
//
//	//get all news
//
//	//get all newscat
//
//	//get all albums
//
//	//get all images
//	b, _ := json.Marshal(common)
//	bs.Data = string(b)
//	request := "createbuild|" + usex.Session
//	b, _ = json.Marshal(bs)
//	c3mcommon.RequestBuildServiceAsync(request, "POST", string(b))
//
//	return ""
//}
//
//func Rebuild(usex models.UserSession) {
//
//	CreateBuild("script", "", "", usex)
//	//CreateBuild("image", "", "", usex)
//
//	//page build:
//	pages := GetAllPage(usex.Shop.Theme, usex.Shop.ID.Hex())
//
//	for _, page := range pages {
//		var langlinks []models.LangLink
//		for langcode, _ := range page.Langs {
//			if page.Langs[langcode].Title == "" {
//				continue
//			}
//			var newslug models.Slug
//			newslug.ShopId = usex.Shop.ID.Hex()
//			newslug.Object = "page"
//			newslug.ObjectId = page.ID.Hex()
//			newslug.Lang = langcode
//			newslug.TemplateCode = usex.Shop.Theme
//
//			//newslug
//			newslug.Slug = mystring.ParameterizeJoin(page.Langs[langcode].Title, "_")
//			//check slug
//			if page.Code == "home" && usex.Shop.Config.DefaultLang == langcode {
//				newslug.Slug = ""
//			}
//			pagelang := page.Langs[langcode]
//			pagelang.Slug = SaveSlug(newslug)
//			page.Langs[langcode] = pagelang
//
//			if page.Langs[langcode].Slug != "" {
//				langlinks = append(langlinks, models.LangLink{Href: page.Langs[langcode].Slug + "/", Code: langcode, Name: c3mcommon.GetLangnameByCode(langcode)})
//			} else {
//				langlinks = append(langlinks, models.LangLink{Href: page.Langs[langcode].Slug, Code: langcode, Name: c3mcommon.GetLangnameByCode(langcode)})
//			}
//			//=====
//		}
//
//		//update
//		page.LangLinks = langlinks
//		SavePage(page)
//		//create build
//		b, _ := json.Marshal(page)
//
//		errstr := CreateBuild("page", page.ID.Hex(), string(b), usex)
//		if errstr != "" {
//			log.Debugf("Rebuild Error :" + errstr)
//		}
//
//	}
//
//	// //productcat build:
//	// prodcats := rpch.GetAllCats(usex.UserID, usex.Shop.ID.Hex())
//	// for _, item := range prodcats {
//	// 	langlinks := make(map[string]string)
//	// 	for lang, _ := range item.Langs {
//	// 		var newslug models.Slug
//	// 		newslug.ShopId = usex.Shop.ID.Hex()
//	// 		newslug.Object = "prodcat"
//	// 		newslug.ObjectId = item.ID.Hex()
//	// newslug.TemplateCode = tmpl.Code
//	// 		newslug.Lang = lang
//	// 		newslug.Domain = usex.ShopConfigs.Domain
//	// 		//newslug
//	// 		tb, _ := lzjs.DecompressFromBase64(item.Langs[lang].Title)
//	// 		newslug.Slug = inflect.Parameterize(string(tb))
//	// 		//save slug
//	// 		b, err := json.Marshal(item)
//	// 		if c3mcommon.CheckError("json encode of item "+item.Code, err) {
//	// 			item.Langs[lang].Slug = rpch.SaveSlug(newslug, string(b))
//	// 			langlinks[lang] = item.Langs[lang].Slug
//	//
//	// 	}
//
//	// 	//check to save langlinks
//	// 	if item.LangLinks != nil {
//	// 		item.LangLinks = langlinks
//	// 		rpch.SaveCat(&item)
//	// 	}
//	// }
//
//	// //product build:
//	// prods := rpch.GetAllProds(usex.UserID, usex.Shop.ID.Hex())
//	// for _, item := range prods {
//	// 	langlinks := make(map[string]string)
//	// 	for lang, _ := range item.Langs {
//	// 		var newslug models.Slug
//	// newslug.TemplateCode = tmpl.Code
//	// 		newslug.ShopId = usex.Shop.ID.Hex()
//	// 		newslug.Object = "product"
//	// 		newslug.ObjectId = item.ID.Hex()
//	// 		newslug.Lang = lang
//	// 		newslug.Domain = usex.ShopConfigs.Domain
//	// 		//newslug
//	// 		log.Debugf("title %s - lang: %s", item.Langs[lang].PageInfo, lang)
//	// 		tb, _ := lzjs.DecompressFromBase64(item.Langs[lang].PageInfo.Title)
//	// 		newslug.Slug = inflect.Parameterize(string(tb))
//	// 		//save slug
//	// 		b, err := json.Marshal(item)
//	// 		if c3mcommon.CheckError("json encode of item "+item.Code, err) {
//	// 			item.Langs[lang].PageInfo.Slug = rpch.SaveSlug(newslug, string(b))
//	// 			langlinks[lang] = item.Langs[lang].PageInfo.Slug
//	// 		}
//	// 	}
//
//	// 	//check to save langlinks
//	// 	if item.LangLinks != nil {
//	// 		item.LangLinks = langlinks
//	// 		rpch.SaveProd(&item)
//	// 	}
//	// }
//
//	// //newscat build:
//	// newscats := rpch.GetAllNewsCats(usex.UserID, usex.Shop.ID.Hex())
//
//	// for _, item := range newscats {
//	// 	langlinks := make(map[string]string)
//	// 	for lang, _ := range item.Langs {
//	// 		var newslug models.Slug
//	//newslug.TemplateCode = tmpl.Code
//	// 		newslug.ShopId = usex.Shop.ID.Hex()
//	// 		newslug.Object = "newscat"
//	// 		newslug.ObjectId = item.ID.Hex()
//	// 		newslug.Lang = lang
//	// 		newslug.Domain = usex.ShopConfigs.Domain
//	// 		//newslug
//	// 		tb, _ := lzjs.DecompressFromBase64(item.Langs[lang].Title)
//	// 		newslug.Slug = inflect.Parameterize(string(tb))
//	// 		//save slug
//	// 		b, err := json.Marshal(item)
//	// 		if c3mcommon.CheckError("json encode of item "+item.Code, err) {
//	// 			item.Langs[lang].Slug = rpch.SaveSlug(newslug, string(b))
//	// 			langlinks[lang] = item.Langs[lang].Slug
//	// 		}
//	// 	}
//
//	// 	//check to save langlinks
//	// 	if item.LangLinks != nil {
//	// 		item.LangLinks = langlinks
//	// 		rpch.SaveNewsCat(&item)
//	// 	}
//	// }
//
//	// //news build:
//	// news := rpch.GetAllNews(usex.UserID, usex.Shop.ID.Hex())
//	// for _, item := range news {
//	// 	langlinks := make(map[string]string)
//	// 	for lang, _ := range item.Langs {
//	// 		var newslug models.Slug
//	// newslug.TemplateCode = tmpl.Code
//	// 		newslug.ShopId = usex.Shop.ID.Hex()
//	// 		newslug.Object = "news"
//	// 		newslug.ObjectId = item.ID.Hex()
//	// 		newslug.Lang = lang
//	// 		newslug.Domain = usex.ShopConfigs.Domain
//	// 		//newslug
//	// 		tb, _ := lzjs.DecompressFromBase64(item.Langs[lang].Title)
//	// 		newslug.Slug = inflect.Parameterize(string(tb))
//	// 		//save slug
//	// 		b, err := json.Marshal(item)
//	// 		if c3mcommon.CheckError("json encode of item "+item.Code, err) {
//	// 			item.Langs[lang].Slug = rpch.SaveSlug(newslug, string(b))
//	// 			langlinks[lang] = item.Langs[lang].Slug
//	// 		}
//	// 	}
//
//	// 	//check to save langlinks
//	// 	if item.LangLinks != nil {
//	// 		item.LangLinks = langlinks
//	// 		rpch.SaveNews(&item)
//	// 	}
//	// }
//	CreateCommonDataBuild(usex)
//
//}
