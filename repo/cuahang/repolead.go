package cuahang

import (
	"github.com/tidusant/c3m/common/c3mcommon"
	"github.com/tidusant/c3m/repo/models"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

func (r *Repo) SaveLead(ld models.Lead) bool {
	start := time.Now()
	col := db.Collection("lpleads")
	rs := true
	if ld.ID.IsZero() {
		//insert
		ld.Created = time.Now()
		_, err := col.InsertOne(ctx, ld)
		r.QueryCount++
		if err != nil {
			c3mcommon.CheckError("SaveLead", err)
			rs = false
		}
	} else {
		//update
		filter := bson.M{"_id": ld.ID}
		update := bson.M{"$set": bson.M{
			"status":   ld.Status,
			"lastsync": ld.Lastsync,
		}}

		_, err := col.UpdateOne(ctx, filter, update)
		if err != nil {
			c3mcommon.CheckError("UpdateLead ", err)
			rs = false
		}
		r.QueryCount++
	}
	r.QueryTime += time.Since(start)
	return rs
}
