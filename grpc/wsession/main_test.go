package main

import (
	"fmt"
	"github.com/tidusant/c3m/common/mystring"
	pb "github.com/tidusant/c3m/grpc/protoc/wsession"
	"sync"

	"os"
	"testing"

	"context"
)

var (
	testsession string = "random"

	ctx            context.Context
	svc            *service
	shopTestId     = "5955d130e761cf70ffb8e49b"
	UserTestId     = "aaaaaaa"
	UserNameTest   = "test"
	URLTest        = "testurl"
	URLAddressTest = "localhost"
)

func setup() {
	// Set up a connection to the server.
	initCheckSession()
	ctx = context.Background()
	svc = &service{}
	//get userid and shopid from session random
	//NOTE: must run auth_test before to have data in db

}
func TestMain(m *testing.M) {
	setup()
	exitVal := m.Run()
	os.Exit(exitVal)
}

func TestCreateSex(t *testing.T) {
	fmt.Println("\n\n==== test TestCreateSex ====")
	rs, err := svc.CreateSession(ctx, nil)
	if err != nil {
		t.Fatalf("Test fail: %s", err.Error())
	}
	//check test data
	fmt.Printf("Data return: %+v\n", rs)
	if rs.Data == "" {
		t.Fatalf("Test fail: not correct return string")
	}
	testsession = rs.Data
}

func TestCheckSessionWithWrongSession(t *testing.T) {
	fmt.Println("\n\n==== test TestCheckSessionWithWrongSession ====")
	rs, err := svc.CheckSession(ctx, &pb.DataRequest{Data: "randomsession"})
	if err != nil {
		t.Fatalf("Test fail: %s", err.Error())
	}
	//check test data
	fmt.Printf("Data return: %+v\n", rs.Data)
	if rs.Data {
		t.Fatalf("Test Fail")
	}
}
func TestCheckSession(t *testing.T) {
	fmt.Println("\n\n==== test TestCheckSession ====")
	rs, err := svc.CheckSession(ctx, &pb.DataRequest{Data: testsession})
	if err != nil {
		t.Fatalf("Test fail: %s", err.Error())
	}
	//check test data
	fmt.Printf("Data return: %+v\n", rs)
	if !rs.Data {
		t.Fatalf("Test Fail")
	}
}
func TestSaveSession(t *testing.T) {
	fmt.Println("\n\n==== test TestSaveSession ====")

	rs, err := svc.SaveSession(ctx, &pb.SessionMessage{Session: testsession, UserID: UserTestId, UserName: UserNameTest})
	if err != nil {
		t.Fatalf("Test fail: %s", err.Error())
	}
	//check test data
	fmt.Printf("Data return: %+v\n", rs)
	if !rs.Data {
		t.Fatalf("Test Fail")
	}
}

func TestGetSession(t *testing.T) {
	fmt.Println("\n\n==== test TestGetSession ====")

	rs, err := svc.GetSession(ctx, &pb.DataRequest{Data: testsession})
	if err != nil {
		t.Fatalf("Test fail: %s", err.Error())
	}
	//check test data
	fmt.Printf("Data return: %+v\n", rs)
	if rs.UserID != UserTestId || rs.UserName != UserNameTest {
		t.Fatalf("Test Fail")
	}
}

func TestCheckRequest(t *testing.T) {
	fmt.Println("\n\n==== test TestCheckRequest ====")
	rs, err := svc.CheckRequest(ctx, &pb.CheckURLRequest{URL: URLTest, Address: URLAddressTest})
	if err != nil {
		t.Fatalf("Test fail: %s", err.Error())
	}
	//check test data
	fmt.Printf("Data return: %+v\n", rs)
	if !rs.Data {
		t.Fatalf("Test Fail")
	}
}

func TestCheckRequestExist(t *testing.T) {
	fmt.Println("\n\n==== test TestCheckRequestExist ====")
	rs, err := svc.CheckRequest(ctx, &pb.CheckURLRequest{URL: URLTest, Address: URLAddressTest})
	if err != nil {
		t.Fatalf("Test fail: %s", err.Error())
	}
	//check test data
	fmt.Printf("Data return: %+v\n", rs)
	if rs.Data {
		t.Fatalf("Test Fail")
	}
}

func TestCheck200Request(t *testing.T) {
	fmt.Println("\n\n==== test TestCheck1000Request ====")
	w := sync.WaitGroup{}
	for i := 0; i < 200; i++ {
		w.Add(1)
		go func(w *sync.WaitGroup) {
			rand := mystring.RandString(10)
			svc.CheckRequest(ctx, &pb.CheckURLRequest{URL: rand, Address: URLAddressTest})
			w.Done()
		}(&w)
	}
	w.Done()
	rs, err := svc.CheckRequest(ctx, &pb.CheckURLRequest{URL: URLTest, Address: URLAddressTest})
	if err != nil {
		t.Fatalf("Test fail: %s", err.Error())
	}
	//check test data
	fmt.Printf("Data return: %+v\n", rs)
	if rs.Data {
		t.Fatalf("Test Fail")
	}
}
