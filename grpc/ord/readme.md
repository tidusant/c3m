# gRPC shop

### how to build test

   - env CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o orders .
    
### run in local:
cd colis/grpcs/session
go run session.go 

### run in docker:
docker build -t tidusant/colis-grpc-ord .  