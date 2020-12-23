# gRPC auth


### how to build test

- compile code:
    - env CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o c3mgrpc_auth .
 
### run test:
note: run auth_test before run this test to get session "random"
env CHADMIN_DB_HOST=127.0.0.1:27017 env CHADMIN_DB_NAME=cuahang env CHADMIN_DB_USER=cuahang env CHADMIN_DB_PASS=cuahang1234@ env PORT=32001 go test


### run in docker:
docker build -t tidusant/c3m-grpc-auth .   