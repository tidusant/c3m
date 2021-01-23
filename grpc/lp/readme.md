# gRPC landing page
    
### run in local:
cd colis/grpcs/lp
env CHADMIN_URI="mongodb://sellpos:sellpos1234%40@127.0.0.1:27017/sellpos?retryWrites=true&w=majority" CHADMIN_DB="sellpos" go run main.go
### run in docker:
docker build -t tidusant/colis-grpc-lp .  

###reference
