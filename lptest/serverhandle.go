package main

import (
	"github.com/gin-gonic/gin"
	"github.com/tidusant/c3m/common/c3mcommon"
	"github.com/tidusant/c3m/common/mycrypto"
	"io/ioutil"
	"os"
)

func SubmitTest(sex, tplname string, c *gin.Context) string {
	content := c.PostForm("data")

	//build content for test
	buildFolder := templateFolder + "/" + tplname
	os.RemoveAll(buildFolder + "/out")
	err := os.Mkdir(buildFolder+"/out", 0755)
	if err != nil {
		return err.Error()
	}
	//copy tailwind css file
	input, err := ioutil.ReadFile(schemeFolder + "/tailwind.css")
	if err != nil {
		return err.Error()
	}
	err = ioutil.WriteFile(buildFolder+"/css/tailwind.css", input, 0644)
	if err != nil {
		return err.Error()
	}

	//create content
	htmlcontent, err := mycrypto.DecompressFromBase64(content)

	if err != nil {
		return err.Error()
	}
	err = ioutil.WriteFile(buildFolder+"/out/content.html", []byte(htmlcontent), 0644)
	if err != nil {
		return err.Error()
	}

	//call test server to purgecss and minify
	bodystr := c3mcommon.RequestAPI2(lpminserver+"/purge", tplname, sex)

	return bodystr

}
