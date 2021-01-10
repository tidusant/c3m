# gRPC shop

### how to build test

   - env CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o orders .
    
### run in local:
cd colis/grpcs/lptpl
env CHADMIN_DB_HOST=127.0.0.1:27017 env CHADMIN_DB_NAME=cuahang env CHADMIN_DB_USER=cuahang env CHADMIN_DB_PASS=cuahang1234@ env PORT=32002 go run orders.go 

### run in docker:
docker build -t tidusant/colis-grpc-lptpl .  