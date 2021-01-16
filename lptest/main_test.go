package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/tidusant/c3m/common/mycrypto"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

var r *gin.Engine
var testsession string

func doPOST(requeststring string, data string) (rs string, err error) {
	//encode data
	//requeststring = mycrypto.EncDat2(requeststring)
	//data = "data=" + mycrypto.EncDat2(data)

	//add body into request
	body := bytes.NewReader([]byte(data))
	req, err := http.NewRequest(http.MethodPost, "/"+requeststring, body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {

		return
	}
	// Create a response recorder so you can inspect the response test
	w := httptest.NewRecorder()

	// Perform the request

	r.ServeHTTP(w, req)

	// Check to see if the response was what you expected
	if w.Code != http.StatusOK {
		err = errors.New(fmt.Sprintf("Expected to get status %d but instead got %d\n", http.StatusOK, w.Code))
		return

	}

	//check data
	//get response body
	bodyresp, err := ioutil.ReadAll(w.Body)
	rs = string(bodyresp)
	//decode data

	//rtstr = mycrypto.DecodeOld(rtstr, 8)
	//json.Unmarshal([]byte(rtstr), &rs)
	return
}

func doCall(testname, requesturl, queryData string, t *testing.T) string {
	fmt.Println("\n\n==== " + testname + " ====")
	fmt.Printf("Data: url: %s - data:%s\n", requesturl, queryData)
	rs, err := doPOST(requesturl, queryData)
	if err != nil {
		t.Fatalf("Test fail: request error: %s", err.Error())
	}
	fmt.Printf("Request return: %+v\n", rs)
	return rs
}

func setup() {
	// Switch to test mode so you don't get such noisy output
	gin.SetMode(gin.TestMode)

	r = gin.Default()
	r.GET("/test/:action/:params", HandleTestRoute)
	r.POST("/test/:action/:params", HandleTestRoute)
}
func TestMain(m *testing.M) {
	setup()
	exitVal := m.Run()
	os.Exit(exitVal)
}

//test special char
func TestNoEncryptUrl(t *testing.T) {
	rs := doCall("TestNoEncryptUrl", "test/edit/templatename", "", t)
	//check test data
	if rs != "invalid url" {
		t.Fatalf("Test fail")
	}
}
func TestNoSession(t *testing.T) {
	rs := doCall("TestNoSession", "test/edit/"+mycrypto.EncDat2("|template"), "", t)
	//check test data
	if rs != "Please login." {
		t.Fatalf("Test fail")
	}
}
