package main

import (
	"github.com/gin-gonic/gin"
	"github.com/tidusant/c3m/common/c3mcommon"
	"io/ioutil"
	"os"
)

func SubmitTest(sex, tplname string, c *gin.Context) string {
	content := c.PostForm("data")

	//build content for test
	buildFolder := templateFolder + "/" + tplname + "/build"
	os.RemoveAll(buildFolder)
	err := os.Mkdir(buildFolder, 0775)
	if err != nil {
		return err.Error()
	}
	//copy all css file
	input, err := ioutil.ReadFile(schemeFolder + "/tailwind.css")
	if err != nil {
		return err.Error()
	}
	err = ioutil.WriteFile(buildFolder+"/tailwind.css", input, 0644)
	if err != nil {
		return err.Error()
	}
	//loop css folder and copy
	if _, err := os.Stat(buildFolder + "/../css"); !os.IsNotExist(err) {
		items, _ := ioutil.ReadDir(buildFolder + "/../css")
		for _, item := range items {
			if !item.IsDir() {
				input, err := ioutil.ReadFile(buildFolder + "/../css/" + item.Name())
				if err != nil {
					return err.Error()
				}
				err = ioutil.WriteFile(buildFolder+"/"+item.Name(), input, 0644)
				if err != nil {
					return err.Error()
				}
			}
		}
	}
	//copy all js file
	if _, err := os.Stat(buildFolder + "/../js"); !os.IsNotExist(err) {
		items, _ := ioutil.ReadDir(buildFolder + "/../js")
		for _, item := range items {
			if !item.IsDir() {
				input, err := ioutil.ReadFile(buildFolder + "/../js/" + item.Name())
				if err != nil {
					return err.Error()
				}
				err = ioutil.WriteFile(buildFolder+"/"+item.Name(), input, 0644)
				if err != nil {
					return err.Error()
				}
			}
		}
	}
	//create content
	err = ioutil.WriteFile(buildFolder+"/content.html", []byte(content), 0644)
	if err != nil {
		return err.Error()
	}

	//call test server to purgecss and minify
	bodystr := c3mcommon.RequestAPI2(lpminserver+"/purge", tplname, sex)
	return bodystr

}
