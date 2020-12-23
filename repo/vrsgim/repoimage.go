package vrsgim

import (
	"github.com/tidusant/c3m-common/c3mcommon"

	"github.com/spf13/viper"
	"github.com/tidusant/chadmin-repo/models"

	"gopkg.in/mgo.v2/bson"
)

/*for upload
=============================================================================
*/
func SaveImages(images []models.CHImage) string {
	col := db.C("files")

	if len(images) == 0 {
		return "0"
	}
	imgfiles := make([]interface{}, len(images))
	for i, image := range images {
		imgfiles[i] = image
	}
	bulk := col.Bulk()
	bulk.Unordered()
	bulk.Insert(imgfiles...)
	_, bulkErr := bulk.Run()
	c3mcommon.CheckError("insert bulk", bulkErr)

	return "1"
}
func SaveImage(image models.CHImage) string {
	col := db.C("files")
	err := col.Insert(image)
	c3mcommon.CheckError("insert SaveImage", err)

	return "1"
}
func ImageCount(shopid string) int {
	col := db.C("files")
	count := -1
	count, err := col.Find(bson.M{"shopid": shopid, "appname": viper.GetString("config.appname")}).Count()
	c3mcommon.CheckError("image count error", err)

	return count
}

func RemoveImage(shopid, filename string) bool {
	col := db.C("files")
	cond := bson.M{"filename": filename}
	var image models.CHImage
	err := col.Find(cond).One(&image)
	if image.Shopid == shopid && image.Filename == filename {

		err = col.Remove(bson.M{"filename": filename, "shopid": shopid})
		if c3mcommon.CheckError("remove image", err) {
			return true
		}

	}

	return false
}

func GetImages(shopid, albumid string) []models.CHImage {
	col := db.C("files")
	var rs []models.CHImage
	cond := bson.M{}

	cond = bson.M{"shopid": shopid, "albumid": albumid, "appname": "chadmin"}

	err := col.Find(cond).All(&rs)
	c3mcommon.CheckError("error get images shopid:"+shopid+" albumid:"+albumid, err)
	return rs
}
