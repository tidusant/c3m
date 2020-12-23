package builder

import (

	//	"c3m/log"

	//"strings"
	"github.com/tidusant/c3m/common/c3mcommon"
	"github.com/tidusant/c3m/repo/models"

	"gopkg.in/mgo.v2/bson"
)

func AuthByKey(key string) models.User {

	col := db.C("users")

	// if prod.Code {

	// 	err := col.Insert(prod)
	// 	c3mcommon.CheckError("product Insert", err)
	// } else {
	var rs models.User
	err := col.Find(bson.M{"keypair": key}).One(&rs)
	c3mcommon.CheckError("AuthByKey", err)
	return rs
}

func ActiveTemplate(userid, shopid string, template, oldtemplate models.Template) string {

	if userid == "" || shopid == "" {
		return ""
	}
	col := db.C("templates")

	cond := bson.M{"status": 1, "code": template.Code}
	change := bson.M{"$set": bson.M{"activedid": template.ActiveIDs}}
	col.Update(cond, change)

	cond = bson.M{"status": 1, "code": oldtemplate.Code}
	change = bson.M{"$set": bson.M{"activedid": oldtemplate.ActiveIDs}}
	col.Update(cond, change)

	return template.Code

}
func UpdateActiveID(templatecode string, ActiveIDs []string) string {

	col := db.C("templates")

	cond := bson.M{"status": 1, "code": templatecode}
	change := bson.M{"$set": bson.M{"activedid": ActiveIDs}}
	err := col.Update(cond, change)
	c3mcommon.CheckError("UpdateActiveID template", err)

	return templatecode
}
func UpdateInstallID(templatecode string, InstalledIDs []string) string {

	col := db.C("templates")

	cond := bson.M{"status": 1, "code": templatecode}
	change := bson.M{"$set": bson.M{"installedid": InstalledIDs}}
	err := col.Update(cond, change)
	c3mcommon.CheckError("install template", err)

	return templatecode
}

func GetAllTemplates() []models.Template {
	col := db.C("templates")
	var rs []models.Template
	err := col.Find(bson.M{"status": 1}).Select(bson.M{"code": 1, "title": 1, "screenshot": 1, "installedid": 1, "activedid": 1}).All(&rs)
	c3mcommon.CheckError("get all templates", err)
	return rs
}

func GetTemplateByCode(code string) models.Template {
	var rs models.Template

	col := db.C("templates")
	err := col.Find(bson.M{"status": 1, "code": code}).One(&rs)
	c3mcommon.CheckError("get template", err)
	return rs
}

func GetAllTemplatesInstalled(userid, shopid string) []models.Template {
	col := db.C("templates")
	var rs []models.Template
	err := col.Find(bson.M{"status": 1, "installedid": shopid}).Select(bson.M{"code": 1, "title": 1, "screenshot": 1}).All(&rs)
	c3mcommon.CheckError("get all templates", err)
	return rs
}

func GetTemplatesByUserId(userid string) []models.Template {
	var rt []models.Template
	col := db.C("templates")
	var cond bson.M
	if userid != "0" {
		cond = bson.M{"userid": userid}
	}

	err := col.Find(cond).All(&rt)
	c3mcommon.CheckError("GetTemplatesByUserId", err)
	return rt
}
func SaveTemplate(newtmpl models.Template) string {
	col := db.C("templates")
	_, err := col.UpsertId(newtmpl.ID, newtmpl)
	c3mcommon.CheckError("UpsertId template", err)
	return newtmpl.Code
}
func GetAllTemplatesCode() map[string]string {
	rt := make(map[string]string)
	col := db.C("templates")
	var cond bson.M
	var rs []models.Template
	err := col.Find(cond).All(&rs)
	c3mcommon.CheckError("GetAllTemplatesCode", err)
	for _, v := range rs {
		rt[v.Code] = v.Code
	}
	return rt
}

func CheckTemplateDup(TemplateTitle string) bool {
	count := 0
	col := db.C("templates")
	var cond bson.M
	cond = bson.M{"title": TemplateTitle}
	count, err := col.Find(cond).Count()
	if c3mcommon.CheckError("CheckTemplateDup", err) && count == 0 {
		return true
	}
	return false
}
