package cuahang

import (
	"github.com/tidusant/c3m/common/c3mcommon"
	"github.com/tidusant/c3m/repo/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func (r *Repo) GetLPByCampID(campID, orgID string, userID primitive.ObjectID) models.LandingPage {
	start := time.Now()
	col := db.Collection("landingpages")
	var rs models.LandingPage
	err := col.FindOne(ctx, bson.M{"userid": userID, "campaignid": campID, "orgid": orgID}).Decode(&rs)
	c3mcommon.CheckError("GetLPByCampID ", err)
	r.QueryCount++
	r.QueryTime += time.Since(start)
	return rs
}
func (r *Repo) GetAllLP(userID primitive.ObjectID) []models.LandingPage {
	start := time.Now()
	col := db.Collection("landingpages")
	var rs []models.LandingPage
	opts := options.Find().SetProjection(bson.M{"_id": 0, "content": 0})

	cursor, err := col.Find(ctx, bson.M{"userid": userID}, opts)
	r.QueryCount++
	if err = cursor.All(ctx, &rs); err != nil {
		c3mcommon.CheckError("GetAllLP ", err)
	}
	r.QueryTime += time.Since(start)
	return rs
}
func (r *Repo) SaveLP(lp models.LandingPage) bool {
	start := time.Now()
	col := db.Collection("landingpages")
	rs := true
	if lp.ID.IsZero() {
		//insert
		_, err := col.InsertOne(ctx, lp)
		r.QueryCount++
		if err != nil {
			rs = false
		}
	} else {
		//update
		filter := bson.M{"campaignid": lp.CampaignID, "orgid": lp.OrgID}
		update := bson.M{"$set": bson.M{
			"content":      lp.Content,
			"sfuserid":     lp.SFUserID,
			"modified":     lp.Modified,
			"lptemplateid": lp.LPTemplateID,
		}}

		_, err := col.UpdateOne(ctx, filter, update)
		if err != nil {
			c3mcommon.CheckError("AddSFUser ", err)
			rs = false
		}
		r.QueryCount++
		r.QueryTime += time.Since(start)
	}
	r.QueryTime += time.Since(start)
	return rs
}
func (r *Repo) DeleteLP(campids []string, userid primitive.ObjectID) bool {
	start := time.Now()
	col := db.Collection("landingpages")
	rs := true
	_, err := col.DeleteMany(ctx, bson.M{"userid": userid, "campaignid": bson.M{"$in": campids}})
	r.QueryCount++
	if err != nil {
		c3mcommon.CheckError("DeleteLP ", err)
		rs = false
	}
	r.QueryTime += time.Since(start)
	return rs
}
