package cuahang

import (
	"github.com/tidusant/c3m/common/c3mcommon"
	"github.com/tidusant/c3m/repo/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func (r *Repo) GetLPByCampID(campID, orgID string) models.LandingPage {
	start := time.Now()
	col := db.Collection("landingpages")
	var rs models.LandingPage
	col.FindOne(ctx, bson.M{"campid": campID, "orgid": orgID}).Decode(&rs)
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
	col := db.Collection("landingpage")
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
		filter := bson.M{"campid": lp.CampaignID, "orgid": lp.OrgID}
		update := bson.M{"$set": bson.M{
			"content":  lp.Content,
			"sfuserid": lp.SFUserID,
			"modified": lp.Modified,
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
