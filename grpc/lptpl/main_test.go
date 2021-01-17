package main

import (
	"encoding/json"
	pb "github.com/tidusant/c3m/grpc/protoc"
	"github.com/tidusant/c3m/repo/models"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"fmt"

	"os"
	"testing"

	"context"
)

var testsession string = "random"
var ctx context.Context
var svc *service
var appname = "test-grpc-lptemplate"
var userId = primitive.NilObjectID
var shopId = primitive.NilObjectID
var shopOriginId = "5955d130e761cf70ffb8e49b" //shopname: demo
var m myRPC

func setup() {
	// Set up a connection to the server.
	ctx = context.Background()
	svc = &service{}
	m = myRPC{}
	//get userid and shopid from session random
	//NOTE: must run auth_test before to have data in db
	userLogin := m.Rpch.GetLogin(testsession)

	userId = userLogin.UserId
	shopId = userLogin.ShopId
	//change to demoshop
	shopOriginIdObj, _ := primitive.ObjectIDFromHex(shopOriginId)
	shopchange := m.Rpch.UpdateShopLogin(testsession, shopOriginIdObj)
	if shopchange.ID.Hex() == "" {
		fmt.Println("Test fail: User can not change to origin shop in setup")
		os.Exit(0)
	}
}
func TestMain(m *testing.M) {
	setup()
	exitVal := m.Run()
	os.Exit(exitVal)
}
func doCall(testname, action, params string, t *testing.T) models.RequestResult {
	fmt.Println("\n\n==== " + testname + " ====")
	resp, err := svc.Call(ctx, &pb.RPCRequest{AppName: appname, Action: action, Params: params, Session: testsession, UserID: userId.Hex(), ShopID: shopId.Hex(), UserIP: "127.0.0.1"})
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
func TestUnknowAction(t *testing.T) {
	rs := doCall("TestUnknowAction", "lasdf", "abc,123", t)
	//check test data
	fmt.Printf("Data return: %+v\n", rs)
	if rs.Data != "Hello "+appname {
		t.Fatalf("Test fail: not correct return string")
	}

}
