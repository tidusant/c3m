### Test server for landing page template
###run at local
 env LPMIN_ADD=http://127.0.0.1:8090 API_ADD=http://127.0.0.1:8081 go run main.go localhandle.go serverhandle.go -debug=true -rootpath=/testtpl
 ###test at local
 env LPMIN_ADD=http://127.0.0.1:8090 API_ADD=http://127.0.0.1:8081 go test
