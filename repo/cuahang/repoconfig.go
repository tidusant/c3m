package cuahang

import (
	"github.com/tidusant/c3m/common/c3mcommon"

	"context"
	"github.com/tidusant/c3m/repo/models"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"

	"gopkg.in/mgo.v2/bson"
)

func GetTemplateConfigs(shopid, templatecode string) []models.TemplateConfig {
	col := db.Collection("addons_configs")
	var rs []models.TemplateConfig
	cursor, err := col.Find(context.TODO(), bson.M{"shopid": shopid, "templatecode": templatecode})
	if err = cursor.All(context.TODO(), &rs); err != nil {
		log.Fatal(err)
	}
	c3mcommon.CheckError("get template configs", err)
	return rs
}
func GetTemplateConfigByKey(shopid, templatecode, key string) models.TemplateConfig {
	col := db.Collection("addons_configs")
	var rs models.TemplateConfig
	cursor, err := col.Find(context.TODO(), bson.M{"shopid": shopid, "templatecode": templatecode, "key": key})
	if err = cursor.All(context.TODO(), &rs); err != nil {
		log.Fatal(err)
	}
	c3mcommon.CheckError("get template configs", err)
	return rs
}

func GetTemplateLang(shopid, templatecode, lang string) []models.TemplateLang {
	col := db.Collection("addons_langs")
	var rs []models.TemplateLang
	cursor, err := col.Find(context.TODO(), bson.M{"shopid": shopid, "templatecode": templatecode, "lang": lang})
	if err = cursor.All(context.TODO(), &rs); err != nil {
		log.Fatal(err)
	}
	c3mcommon.CheckError("get template langs", err)
	return rs
}

func SaveConfigs(newitem models.TemplateConfig) {
	col := db.Collection("addons_configs")

	if newitem.ID == "" {
		newitem.ID = bson.NewObjectId()
	}
	opts := options.Update().SetUpsert(true)
	filter := bson.D{{"_id", newitem.ID}}
	update := bson.D{{"$set", newitem}}

	_, err := col.UpdateOne(context.TODO(), filter, update, opts)
	c3mcommon.CheckError("save template configs", err)

}
