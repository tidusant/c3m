package vrsgim

import (
	"os"

	"github.com/tidusant/c3m-common/c3mcommon"
	"github.com/tidusant/c3m-common/log"

	"gopkg.in/mgo.v2"
)

var (
	db *mgo.Database
)

func init() {
	log.Infof("init repo vrsgim...")
	strErr := ""
	db, strErr = c3mcommon.ConnectDB("vrsgim")
	if strErr != "" {
		log.Infof(strErr)
		os.Exit(1)
	}
	log.Info("done")
}
