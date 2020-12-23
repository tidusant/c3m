package builder

import (
	"github.com/tidusant/c3m/common/c3mcommon"
	"github.com/tidusant/c3m/common/log"
	"github.com/tidusant/c3m/repo/models"

	//	"c3m/log"

	//"strings"

	"gopkg.in/mgo.v2/bson"
)

//=========================cat function==================
func SaveResource(newitem models.Resource) string {

	col := db.C("resources")
	cond := bson.M{"shopid": newitem.ShopID, "templatecode": newitem.TemplateCode, "key": newitem.Key}
	change := bson.M{"$set": bson.M{"value": newitem.Value}}
	err := col.Update(cond, change)

	c3mcommon.CheckError("SaveResource Update", err)
	return ""
}
func GetAllResource(templatecode, shopid string) []models.Resource {
	col := db.C("resources")
	var rs []models.Resource
	err := col.Find(bson.M{"shopid": shopid, "templatecode": templatecode}).All(&rs)
	c3mcommon.CheckError("GetAllResource", err)
	return rs
}
func GetResourceByKey(templatecode, shopid, key string) models.Resource {
	col := db.C("resources")
	var rs models.Resource
	cond := bson.M{"shopid": shopid, "key": key, "templatecode": templatecode}

	err := col.Find(cond).One(&rs)
	c3mcommon.CheckError("GetResourceByKey", err)
	return rs
}
func InsertResource(rs models.Resource) {
	col := db.C("resources")
	//check  exist:
	cond := bson.M{"shopid": rs.ShopID, "templatecode": rs.TemplateCode, "key": rs.Key}
	var oldrs models.Resource
	col.Find(cond).One(&oldrs)
	if oldrs.ID.Hex() != "" {
		//skip if exist
		log.Debugf("exist, skip")
		return
	}
	err := col.Insert(rs)
	c3mcommon.CheckError("Insert Resource "+rs.Key+" template:"+rs.TemplateCode, err)
}

func RemoveOldTemplateResource(shop models.Shop, template models.Template) {
	//remove old config
	colcfg := db.C("resources")
	cond := bson.M{"shopid": shop.ID.Hex(), "templatecode": template.Code}
	_, err := colcfg.RemoveAll(cond)
	c3mcommon.CheckError("remove old template resources,shopid:"+shop.ID.Hex()+",templatecode:"+template.Code, err)
}
func RemoveUnusedTemplateResource(shopid string, template models.Template, isntalledResource []string) {
	//remove old config
	colcfg := db.C("resources")
	cond := bson.M{"shopid": shopid, "templatecode": template.Code, "key": bson.M{"$nin": isntalledResource}}
	_, err := colcfg.RemoveAll(cond)
	c3mcommon.CheckError("remove old template resources,shopid:"+shopid+",templatecode:"+template.Code, err)
}
