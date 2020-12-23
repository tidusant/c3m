# gRPC shop

### how to build test
- create go.mod file
- update latest dependency:
    go clean --modcache
    go get github.com/tidusant/c3m-common/c3mcommon@master
    go get github.com/tidusant/c3m-common/log@master
    go get github.com/tidusant/c3m-common/mycrypto@master
    go get github.com/tidusant/c3m-common/mystring@master
    go get github.com/tidusant/chadmin-repo/models@master
    go get github.com/tidusant/chadmin-repo/session@master
    go get github.com/tidusant/chadmin-repo/cuahang@master
    
    go get github.com/tidusant/c3m-grpc-protoc/protoc
- compile code:
    - env CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o c3mgrpc_auth .
    
### run in local:
cd colis/grpcs/auth
env CHADMIN_DB_HOST=127.0.0.1:27017 env CHADMIN_DB_NAME=cuahang env CHADMIN_DB_USER=cuahang env CHADMIN_DB_PASS=cuahang1234@ env PORT=32002 go run shop.go 

### run in docker:
docker build -t tidusant/colis-grpc-shop .  