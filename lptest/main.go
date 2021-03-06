package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tidusant/c3m/common/c3mcommon"
	"github.com/tidusant/c3m/common/log"
	"github.com/tidusant/c3m/common/mycrypto"
	"github.com/tidusant/c3m/repo/models"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	//"github.com/gin-gonic/contrib/static"
)

var (
	loaddatadone bool
	layoutPath   = "./template/out"

	schemeFolder   = "./scheme"
	templateFolder = "./templates"
	rootPath       = ""
	apiserver      string
	lpminserver    string
)

func main() {
	initdata()
	if !loaddatadone {
		log.Errorf("Load data fail.")
		return
	}
	var port int
	var debug bool
	//check port
	rand.Seed(time.Now().Unix())
	port = 8082

	//fmt.Println(mycrypto.Encode("abc,efc", 5))
	flag.BoolVar(&debug, "debug", true, "Indicates if debug messages should be printed in log files")
	rootPath = os.Getenv("ROOTPATH")
	flag.Parse()

	logLevel := log.DebugLevel
	if !debug {
		layoutPath = "./layout"
		logLevel = log.InfoLevel
		gin.SetMode(gin.ReleaseMode)
		log.SetOutputFile(fmt.Sprintf("portal-"+strconv.Itoa(port)), logLevel)
		defer log.CloseOutputFile()
		log.RedirectStdOut()
	}
	log.Infof("debug %v", debug)

	//init config
	router := gin.Default()

	//router.Use(static.Serve("/", static.LocalFile("static", false)))
	router.StaticFile("/", layoutPath+"/index.html")
	//nextjs request File
	router.Static("/templates", "./templates")
	router.Static("/scheme", "./scheme")
	//router.StaticFile("/edit", layoutPath+"/edit.html")
	//router.LoadHTMLGlob(layoutPath+"/edit.html")

	router.GET("/test/:action/:params", HandleTestRoute)
	router.POST("/test/:action/:params", HandleTestRoute)

	log.Infof("running with port:" + strconv.Itoa(port))
	router.Run(":" + strconv.Itoa(port))

}

func HandleTestRoute(c *gin.Context) {
	//get cookie
	var rs models.RequestResult
	args := strings.Split(mycrypto.Decode(c.Param("params")), "|")
	if len(args) < 2 {
		rs.Error = "invalid url"
		dataReturn(c, rs)
		return
	}
	sex := args[0]
	tplname := args[1]
	log.Debugf("cookies: %+v", sex)
	c.Writer.WriteHeader(http.StatusOK)
	if sex == "" {
		rs.Error = "Please login."
		dataReturn(c, rs)
		return
	}
	//get session to auth
	bodystr := c3mcommon.RequestAPI(apiserver, "aut", sex+"|t")

	err := json.Unmarshal([]byte(bodystr), &rs)

	if err != nil {
		rs.Error = err.Error()
		dataReturn(c, rs)
		return
	}
	if rs.Status != 1 {
		dataReturn(c, rs)
		return
	}
	var rt map[string]string
	err = json.Unmarshal([]byte(rs.Data), &rt)
	if err != nil {
		rs.Error = err.Error()
		dataReturn(c, rs)
		return
	}
	if v, ok := rt["username"]; !ok || v == "" {

		rs.Error = "Please login again."
		dataReturn(c, rs)
		return
	}
	//get modules permission from session
	modules := make(map[string]bool)
	for _, v := range strings.Split(rt["modules"], ",") {
		modules[v] = true
	}
	//check module permission
	if ok, _ := modules["c3m-lptpl-admin"]; !ok {

		rs.Error = "Permission denied."
		dataReturn(c, rs)
		return
	}

	action := c.Param("action")
	switch action {
	case "edit":

		c.Writer.WriteString(GetTest(sex, tplname, c))
	case "submit":
		c.Writer.WriteString(SubmitTest(sex, tplname, c))
	}
}

func dataReturn(c *gin.Context, rs models.RequestResult) {
	if c.Request.Method == "GET" {
		if rs.Status == 1 {
			c.Writer.WriteString(rs.Data)
		} else {
			c.Writer.WriteString(rs.Error)
		}
	} else {
		b, _ := json.Marshal(rs)
		c.Writer.WriteString(string(b))
	}
}

func initdata() {
	apiserver = os.Getenv("API_ADD")
	lpminserver = os.Getenv("LPMIN_ADD")
	if len(apiserver) < 10 {
		log.Error("Api ip INVALID")
		os.Exit(0)
	}
	loaddatadone = true
}
