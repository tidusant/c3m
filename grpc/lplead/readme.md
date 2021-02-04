# gRPC lplead
    
### run in local:
```
env CHADMIN_URI="mongodb://sellpos:sellpos1234%40@127.0.0.1:27017/sellpos?retryWrites=true&w=majority" CHADMIN_DB="sellpos" go run main.go 
```
### test in local
```
env CHADMIN_URI="mongodb://sellpos:sellpos1234%40@127.0.0.1:27017/sellpos?retryWrites=true&w=majority" CHADMIN_DB="sellpos" go test
```