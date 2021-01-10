package cuahang

import (
	"fmt"

	"github.com/tidusant/c3m/repo/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

func (r *Repo) UpdateLpTemplate(template models.LPTemplate) error {
	col := db.Collection("lptemplates")
	//update
	cond := bson.D{{"_id", template.ID}}
	update := bson.D{{"$set", bson.D{
		{"status", template.Status},
		{"viewed", template.Viewed},
		{"installed", template.Installed},
	}}}
	_, err := col.UpdateOne(ctx, cond, update)
	r.QueryCount++
	if err != nil {
		return err
	}
	return nil
}
func (r *Repo) CreateLpTemplate(userid primitive.ObjectID, templatename string) error {

	col := db.Collection("lptemplates")
	cond := bson.M{"userid": userid, "name": templatename}
	count, err := col.CountDocuments(ctx, cond)
	if err != nil {
		return err
	}
	r.QueryCount++
	if count > 0 {
		return fmt.Errorf("Duplicate template name: %s", templatename)
	}
	//insert
	_, err = col.InsertOne(ctx, models.LPTemplate{UserID: userid, Name: templatename, Status: 2, Created: time.Now()})
	r.QueryCount++
	if err != nil {
		return err
	}
	return nil
}
func (r *Repo) GetLpTemplate(userid primitive.ObjectID, templatename string) (models.LPTemplate, error) {
	col := db.Collection("lptemplates")
	var rs models.LPTemplate
	cond := bson.M{"userid": userid, "name": templatename}
	//query
	err := col.FindOne(ctx, cond).Decode(&rs)
	r.QueryCount++
	if err != nil {
		return rs, err
	}
	return rs, nil
}

func (r *Repo) GetAllLpTemplate(userid primitive.ObjectID) ([]models.LPTemplate, error) {
	col := db.Collection("lptemplates")
	var rs []models.LPTemplate
	cond := bson.M{"userid": userid}
	//query
	cursor, err := col.Find(ctx, cond)
	r.QueryCount++
	if err = cursor.All(ctx, &rs); err != nil {
		return rs, err
	}
	return rs, nil
}
