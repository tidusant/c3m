package builder

import (
	"github.com/tidusant/c3m/common/c3mcommon"
	"github.com/tidusant/c3m/common/log"
	"github.com/tidusant/c3m/repo/models"

	"gopkg.in/mgo.v2/bson"
)

func GetBuildConfig(shopid string) models.BuildConfig {
	col := db.C("configs")
	var bs models.BuildConfig

	err := col.Find(bson.M{"shopid": shopid}).One(&bs)
	c3mcommon.CheckError("GetConfigs", err)
	return bs

}

func GetTemplateConfigs(shopid, templatecode string) []models.TemplateConfig {
	col := db.C("templateconfigs")
	var rs []models.TemplateConfig
	err := col.Find(bson.M{"shopid": shopid, "templatecode": templatecode}).All(&rs)

	c3mcommon.CheckError("get template configs", err)
	return rs
}
func GetTemplateConfigByKey(shopid, templatecode, key string) models.TemplateConfig {
	col := db.C("templateconfigs")
	var rs models.TemplateConfig
	err := col.Find(bson.M{"shopid": shopid, "templatecode": templatecode, "key": key}).One(&rs)
	c3mcommon.CheckError("get template "+templatecode+" configs "+key, err)
	return rs
}

func InsertTemplateConfig(config models.TemplateConfig) {
	col := db.C("templateconfigs")
	//check  exist:
	cond := bson.M{"shopid": config.ShopID, "templatecode": config.TemplateCode, "key": config.Key}
	var oldrs models.Resource
	col.Find(cond).One(&oldrs)
	if oldrs.ID.Hex() != "" {
		//skip if exist
		log.Debugf("exist, skip")
		return
	}
	err := col.Insert(config)
	c3mcommon.CheckError("InsertConfig "+config.Key+" template:"+config.TemplateCode, err)

}
func SaveConfigByShopId(config models.BuildConfig) {
	col := db.C("configs")
	//check  exist:

	var oldcf models.BuildConfig
	err := col.Find(bson.M{"shopid": config.ShopId}).One(&oldcf)
	if oldcf.ID.Hex() == "" {
		return
	}
	config.ID = oldcf.ID
	_, err = col.UpsertId(config.ID, config)
	c3mcommon.CheckError("UpsertId SaveConfig", err)
}
func SaveTemplateConfig(config models.TemplateConfig) {
	col := db.C("templateconfigs")
	//check  exist:
	cond := bson.M{"shopid": config.ShopID, "templatecode": config.TemplateCode, "key": config.Key}
	change := bson.M{"$set": bson.M{"value": config.Value}}
	err := col.Update(cond, change)
	c3mcommon.CheckError("SaveConfig "+config.Key+" template:"+config.TemplateCode, err)

}

func RemoveOldTemplateConfig(shop models.Shop, template models.Template) {
	//remove old config
	colcfg := db.C("templateconfigs")
	cond := bson.M{"shopid": shop.ID.Hex(), "templatecode": template.Code}
	_, err := colcfg.RemoveAll(cond)
	c3mcommon.CheckError("remove old template config,shopid:"+shop.ID.Hex()+",templatecode:"+template.Code, err)
}
func RemoveUnusedTemplateConfig(shopid string, template models.Template, installedConfigs []string) {
	//remove old config
	colcfg := db.C("templateconfigs")

	cond := bson.M{"shopid": shopid, "templatecode": template.Code, "key": bson.M{"$nin": installedConfigs}}
	_, err := colcfg.RemoveAll(cond)
	c3mcommon.CheckError("remove old template config,shopid:"+shopid+",templatecode:"+template.Code, err)
}
