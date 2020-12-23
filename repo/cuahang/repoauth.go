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

func GetUserInfo(UserId primitive.ObjectID) models.User {
	col := db.Collection("addons_users")
	var rs models.User
	col.FindOne(ctx, bson.M{"_id": UserId}).Decode(&rs)
	return rs
}

//get user login by session and return current shop and user id
func GetLogin(session string) models.UserLogin {
	coluserlogin := db.Collection("addons_userlogin")
	var rs models.UserLogin
	cond := bson.M{"session": session}
	log.Debugf("cond getlogin:%+v", cond)
	err := coluserlogin.FindOne(ctx, cond).Decode(&rs)
	if c3mcommon.CheckError("Error GetLogin", err) {
		if rs.ShopId == primitive.NilObjectID {
			rs.ShopId = GetShopDefault(rs.UserId)
			log.Debugf("GetLogin Shopid", rs.ShopId)
			filter := bson.D{{"userid", rs.UserId}}
			update := bson.D{{"$set", bson.M{"shopid": rs.ShopId}}}
			coluserlogin.UpdateOne(ctx, filter, update)
		}
	}
	return rs
}
func UpdateShopLogin(session string, ShopChangeId primitive.ObjectID) (shopchange models.Shop) {
	coluserlogin := db.Collection("addons_userlogin")
	var rs models.UserLogin
	err := coluserlogin.FindOne(ctx, bson.M{"session": session}).Decode(&rs)
	log.Debugf("query user %s,%s, %v,%s", rs.Session, rs.ID, rs.ShopId, rs.UserId)
	c3mcommon.CheckError("Error query Session in UpdateShopLogin", err)
	if rs.UserId.Hex() == "" {
		return shopchange
	}
	//get shop id

	shopchange = GetShopById(rs.UserId, ShopChangeId)
	if shopchange.ID == primitive.NilObjectID {
		return shopchange
	}
	rs.ShopId = shopchange.ID

	log.Debugf("shopid:%s", rs.ShopId)
	filter := bson.M{"_id": rs.ID}
	update := bson.M{"$set": bson.M{"shopid": rs.ShopId}}
	_, err = coluserlogin.UpdateOne(ctx, filter, update)
	c3mcommon.CheckError("Error update Session in UpdateShopLogin", err)
	return shopchange
}

//Login user and update session
func Login(user, pass, session, userIP string) models.User {
	hash := md5.Sum([]byte(pass))
	passmd5 := hex.EncodeToString(hash[:])
	coluser := db.Collection("addons_users")

	log.Debugf("login:%s - %s", user, passmd5)
	var result models.User
	err := coluser.FindOne(ctx, bson.M{"user": user, "password": passmd5}).Decode(&result)
	c3mcommon.CheckError("error query user", err)
	log.Debugf("user result %v", result)
	if result.Name != "" {
		coluserlogin := db.Collection("addons_userlogin")
		var userlogin models.UserLogin
		coluserlogin.FindOne(ctx, bson.M{"userid": result.ID}).Decode(&userlogin)
		if userlogin.UserId.Hex() == "" {
			userlogin.UserId = result.ID
		}
		userlogin.LastLogin = time.Now().UTC()
		userlogin.LoginIP = userIP
		userlogin.Session = session

		opts := options.Update().SetUpsert(true)
		filter := bson.M{"userid": userlogin.UserId}
		update := bson.M{"$set": bson.M{
			"last":    userlogin.LastLogin,
			"id":      userlogin.LoginIP,
			"session": userlogin.Session,
		}}

		_, err := coluserlogin.UpdateOne(ctx, filter, update, opts)
		c3mcommon.CheckError("Upsert login", err)

	}
	return result
}
func Logout(session string) string {

	col := db.Collection("addons_userlogin")
	col.DeleteOne(ctx, bson.M{"session": session})

	return ""
}
