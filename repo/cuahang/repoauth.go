package cuahang

import (
	"crypto/md5"
	"encoding/hex"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/tidusant/c3m/common/c3mcommon"
	"github.com/tidusant/c3m/common/log"
	"github.com/tidusant/c3m/repo/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2/bson"
	"time"
)

/*
Check login status by session and USERIP, return userId and shopid
*/

func (r *Repo) GetUserInfo(UserId primitive.ObjectID) models.User {
	start := time.Now()
	col := db.Collection("addons_users")
	var rs models.User
	col.FindOne(ctx, bson.M{"_id": UserId}).Decode(&rs)
	r.QueryCount++
	r.QueryTime += time.Since(start)
	return rs
}

//get user login by session and return current shop and user id
func (r *Repo) GetLogin(session string) models.UserLogin {
	start := time.Now()
	coluserlogin := db.Collection("addons_userlogin")
	var rs models.UserLogin
	cond := bson.M{"session": session}

	err := coluserlogin.FindOne(ctx, cond).Decode(&rs)
	r.QueryCount++
	if c3mcommon.CheckError("Error GetLogin", err) {
		if rs.ShopId == primitive.NilObjectID {

			shop := r.GetShopDefault(rs.UserId)
			rs.ShopId = shop.ID
			filter := bson.M{"userid": rs.UserId}
			update := bson.M{"$set": bson.M{"shopid": rs.ShopId}}
			_, err := coluserlogin.UpdateOne(ctx, filter, update)
			c3mcommon.CheckError("GetLogin UpdateOne", err)
			r.QueryCount++
		}
	}
	r.QueryTime += time.Since(start)
	return rs
}
func (r *Repo) UpdateShopLogin(session string, ShopChangeId primitive.ObjectID) (shopchange models.Shop) {
	start := time.Now()
	coluserlogin := db.Collection("addons_userlogin")
	var rs models.UserLogin
	err := coluserlogin.FindOne(ctx, bson.M{"session": session}).Decode(&rs)
	r.QueryCount++
	c3mcommon.CheckError("UpdateShopLogin FindOne", err)
	if rs.UserId.Hex() == "" {
		return shopchange
	}
	//get shop id

	shopchange = r.GetShopById(rs.UserId, ShopChangeId)
	if shopchange.ID == primitive.NilObjectID {
		return shopchange
	}
	rs.ShopId = shopchange.ID

	log.Debugf("shopid:%s", rs.ShopId)
	filter := bson.M{"_id": rs.ID}
	update := bson.M{"$set": bson.M{"shopid": rs.ShopId}}
	_, err = coluserlogin.UpdateOne(ctx, filter, update)
	r.QueryCount++
	c3mcommon.CheckError("Error update Session in UpdateShopLogin", err)
	r.QueryTime += time.Since(start)
	return shopchange
}

//Login user and update session
func (r *Repo) Login(user, pass, session, userIP string) models.User {
	start := time.Now()
	hash := md5.Sum([]byte(pass))
	passmd5 := hex.EncodeToString(hash[:])
	coluser := db.Collection("addons_users")

	log.Debugf("login:%s - %s", user, passmd5)
	var result models.User
	err := coluser.FindOne(ctx, bson.M{"user": user, "password": passmd5}).Decode(&result)
	r.QueryCount++
	c3mcommon.CheckError("error query user", err)
	log.Debugf("user result %v", result)
	if result.Name != "" {
		coluserlogin := db.Collection("addons_userlogin")
		var userlogin models.UserLogin
		err := coluserlogin.FindOne(ctx, bson.M{"userid": result.ID}).Decode(&userlogin)
		r.QueryCount++
		if c3mcommon.CheckError("Login FindOne", err) {
			userlogin.UserId = result.ID

			userlogin.LastLogin = time.Now().UTC()
			userlogin.LoginIP = userIP
			userlogin.Session = session

			opts := options.Update().SetUpsert(true)
			filter := bson.M{"userid": userlogin.UserId}
			update := bson.M{"$set": bson.M{
				"last":    userlogin.LastLogin,
				"ip":      userlogin.LoginIP,
				"session": userlogin.Session,
			}}

			_, err := coluserlogin.UpdateOne(ctx, filter, update, opts)
			r.QueryCount++
			c3mcommon.CheckError("Upsert login", err)
		}

	}
	r.QueryTime += time.Since(start)
	return result
}
func (r *Repo) Logout(session string) string {
	start := time.Now()
	col := db.Collection("addons_userlogin")
	col.DeleteOne(ctx, bson.M{"session": session})
	r.QueryCount++
	r.QueryTime += time.Since(start)
	return ""
}
