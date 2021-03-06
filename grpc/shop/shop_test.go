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
var appname = "test-grpc-shop"
var userId = primitive.NilObjectID
var shopId = primitive.NilObjectID
var shopChangeNoPermissionId = primitive.NilObjectID //shopname: Hoa Giay
var shopChangeId = primitive.NilObjectID
var shopOriginId = primitive.NilObjectID
var m myRPC

func setup() {
	// Set up a connection to the server.
	ctx = context.Background()
	svc = &service{}
	m = myRPC{}
	//get userid and shopid from session random
	//NOTE: must run auth_test before to have data in db
	userLogin := m.Rpch.GetLogin(testsession)
	shopChangeId, _ = primitive.ObjectIDFromHex("59a275355f4aec1b026b6f5e")             //shopname: Colis
	shopOriginId, _ = primitive.ObjectIDFromHex("5955d130e761cf70ffb8e49b")             //shopname: demo
	shopChangeNoPermissionId, _ = primitive.ObjectIDFromHex("5a52f89a7b4b30ed5ecfce9d") //shopname: Hoa Giay test
	userId = userLogin.UserId
	shopId = userLogin.ShopId
	//change to demoshop
	shopchange := m.Rpch.UpdateShopLogin(testsession, shopOriginId)
	if shopchange.ID.Hex() == "" {
		fmt.Println("Test fail: User can not change to origin shop in setup")
		os.Exit(0)
	}
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
func TestMain(m *testing.M) {
	setup()
	exitVal := m.Run()
	os.Exit(exitVal)
}

func TestUnknowAction(t *testing.T) {
	rs := doCall("TestUnknowAction", "lasdf", "abc,123", t)
	//check test data
	fmt.Printf("Data return: %+v\n", rs)
	if rs.Data != "Hello "+appname {
		t.Fatalf("Test fail: not correct return string")
	}

}

func TestChangeShopWithoutShopPermission(t *testing.T) {
	rs := doCall("TestChangeShopWithoutShopPermission", "cs", shopChangeNoPermissionId.Hex(), t)
	if rs.Status == 1 {
		t.Fatalf("Test fail: User can change shop whithout shop permission")
	}
}

func TestChangeShop(t *testing.T) {
	rs := doCall("TestChangeShop", "cs", shopChangeId.Hex(), t)
	if rs.Status != 1 {
		t.Fatalf("Test fail: User can not change shop")
	}
	//change to demoshop
	shopchange := m.Rpch.UpdateShopLogin(testsession, shopOriginId)
	if shopchange.ID.Hex() == "" {
		fmt.Println("Test fail: User can not change to origin shop after TestChangeShop")
		os.Exit(0)
	}
}
