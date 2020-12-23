package builder

import (
	"time"

	"github.com/tidusant/c3m-common/c3mcommon"
	"github.com/tidusant/c3m-common/log"
	"github.com/tidusant/chadmin-repo/models"
	"os"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	db *mgo.Database
)

func init() {
	log.Infof("init repo build...")
	strErr := ""
	db, strErr = c3mcommon.ConnectDB("chbuild")
	if strErr != "" {
		log.Infof(strErr)
		os.Exit(1)
	}
	log.Info("done")
}

//query and update https://stackoverflow.com/questions/11417784/mongodb-in-go-golang-with-mgo-how-do-i-update-a-record-find-out-if-update-wa
func CreateBuild(shopid string, bs models.BuildScript) string {
	col := db.C("builds")
	//remove old build
	cond := bson.M{"object": bs.Object, "shopid": shopid, "objectid": bs.ObjectId}
	//"objectid": buildscript.ObjectID, "collection": buildscript.Collection}

	_, err := col.RemoveAll(cond)

	bs.Status = 0
	bs.Created = time.Now().Unix()
	bs.Modified = time.Now().Unix()
	err = col.Insert(bs)
	c3mcommon.CheckError("insert build script", err)

	return ""
}

func GetBuild() models.BuildScript {
	col := db.C("builds")
	var bs models.BuildScript
	change := mgo.Change{
		Update:    bson.M{"$set": bson.M{"status": 1, "modified": time.Now().Unix()}},
		ReturnNew: true,
	}
	col.Find(bson.M{"status": 0}).Apply(change, &bs)

	return bs

}

func RemoveAllBuild(shopID string) string {
	col := db.C("builds")
	cond := bson.M{"shopid": shopID}

	_, err := col.RemoveAll(cond)
	c3mcommon.CheckError("RemoveAllBuild script", err)
	return ""
}
