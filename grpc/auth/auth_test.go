package main

import (
	"github.com/tidusant/c3m/common/c3mcommon"

	pb "github.com/tidusant/c3m-grpc-protoc/protoc"
	"github.com/tidusant/c3m/repo/models"

	"encoding/json"
	"strings"

	"fmt"

	"os"
	"testing"

	"context"
)

var testsession string = "random"
var ctx context.Context
var svc *service
var appname = "test-grpc-auth"

func setup() {
	// Set up a connection to the server.
	ctx = context.Background()
	svc = &service{}

}
func doCall(testname, action, params string, t *testing.T) models.RequestResult {
	//fmt.Println("\n\n==== " + testname + " ====")
	resp, err := svc.Call(ctx, &pb.RPCRequest{AppName: appname, Action: action, Params: params, Session: testsession, UserIP: "127.0.0.1"})
	if err != nil {
		t.Fatalf("Test fail: Service error: %s", err.Error())
	}
	fmt.Printf("response return: %+v\n", resp)
	//check test data
	var rs models.RequestResult
	json.Unmarshal([]byte(resp.Data), &rs)
	fmt.Printf("Data return: %+v\n", rs)
	return rs
}
func TestMain(m *testing.M) {
	setup()
	exitVal := m.Run()
	os.Exit(exitVal)
}

func TestUnknowAction(t *testing.T) {
	fmt.Println("==== test TestUnknowAction ====")
	rs, err := svc.Call(ctx, &pb.RPCRequest{AppName: appname, Action: "lasdf", Params: "abc,123", Session: testsession, UserIP: "127.0.0.1"})
	if err != nil {
		t.Fatalf("Test fail: Service error: %s", err.Error())
	}
	//check test data
	fmt.Printf("Data return: %+v\n", rs)
	if rs.Data != "Hello "+appname {
		t.Fatalf("Test fail: not correct return string")
	}

}

func TestLoginWithSpecialChar(t *testing.T) {
	rs := doCall("TestLoginWithSpecialChar", "l", c3mcommon.GetSpecialChar(), t)
	if rs.Status != 0 {
		t.Fatalf("Test fail: User logged in")
	}
}

func TestLoginFail(t *testing.T) {
	rs := doCall("TestLoginFail", "l", "abc,123", t)
	if rs.Status != 0 {
		t.Fatalf("Test fail: User logged in")
	}
}
func TestLoginCorrect(t *testing.T) {
	rs := doCall("TestLoginCorrect", "l", "demo,123", t)
	if rs.Status != 1 {
		t.Fatalf("Test fail: User cannot log in")
	}
}
func TestCheckLoginWithSession(t *testing.T) {
	rs := doCall("TestCheckLoginWithSession", "aut", "", t)
	if rs.Status != 1 {
		t.Fatalf("Test fail: User login but cannot test")
	}
}
func TestCheckLoginWithSession_Portal(t *testing.T) {
	rs := doCall("TestCheckLoginWithSession_Portal", "aut", "", t)
	if rs.Status != 1 || strings.Index(rs.Data, `"userid"`) < 1 {
		t.Fatalf("Test fail: User login but cannot auth")
	}
}
