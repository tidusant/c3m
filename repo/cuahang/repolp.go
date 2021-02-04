package cuahang

import (
	"github.com/tidusant/c3m/common/c3mcommon"
	"github.com/tidusant/c3m/repo/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func (r *Repo) LPCheckDomainUrlExist(domainUrl string) bool {
	start := time.Now()
	col := db.Collection("landingpages")
	var rs models.LandingPage
	err := col.FindOne(ctx, bson.M{"customhost": true, "domainname": domainUrl}).Decode(&rs)
	c3mcommon.CheckError("LPCheckDomainUrlExist", err)
	r.QueryCount++
	r.QueryTime += time.Since(start)
	return !rs.ID.IsZero()
}
func (r *Repo) GetLPByCampID(campID string, userID primitive.ObjectID) models.LandingPage {
	start := time.Now()
	col := db.Collection("landingpages")
	var rs models.LandingPage
	err := col.FindOne(ctx, bson.M{"userid": userID, "campaignid": campID}).Decode(&rs)
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
func (r *Repo) GetAllLpByCampIds(campIds []string, userid primitive.ObjectID) []models.LPTemplate {
	col := db.Collection("landingpages")
	var rs []models.LPTemplate
	cond := bson.M{"userid": userid, "campaignid": bson.M{"$in": campIds}}
	cursor, err := col.Find(ctx, cond)
	r.QueryCount++
	if err = cursor.All(ctx, &rs); err != nil {
		c3mcommon.CheckError("AddSFUser ", err)
	}
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
			"path":           lp.Path,
			"favicon":        lp.Favicon,
			"successmessage": lp.SuccessMessage,
			"successtitle":   lp.SuccessTitle,
			"content":        lp.Content,
			"sfuserid":       lp.SFUserID,
			"modified":       lp.Modified,
			"lptemplateid":   lp.LPTemplateID,
			"customhost":     lp.CustomHost,
			"domainname":     lp.DomainName,
			"ftphost":        lp.FTPHost,
			"ftpuser":        lp.FTPUser,
			"ftppass":        lp.FTPPass,
		}}

		_, err := col.UpdateOne(ctx, filter, update)
		if err != nil {
			c3mcommon.CheckError("AddSFUser ", err)
			rs = false
		}
		r.QueryCount++
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
