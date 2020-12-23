#admin portal

### how to build test

   - env CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

### run in local:
cd colis/portals/admin

env AUTH_IP=127.0.0.1:8901 SHOP_IP=127.0.0.1:8902 SESSION_URI="mongo_server_uri" SESSION_DB="mongodb_name" CHADMIN_URI="mongo_server_url" CHADMIN_DB="mongodb_name"  go run main.go 

### build & run in docker:

docker build -t tidusant/c3madmin-portal . 
 

### reference:
https://kubernetes.io/docs/tasks/administer-cluster/dns-debugging-resolution/
