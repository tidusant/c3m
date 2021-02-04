#api for admin
### run in local:
```
env SESSION_IP=127.0.0.1:8864 AUTH_IP=127.0.0.1:8901 SHOP_IP=127.0.0.1:8902 ORD_IP=127.0.0.1:8903 PAGE_IP=127.0.0.1:8904 LPTPL_IP=127.0.0.1:8905 LP_IP=127.0.0.1:8906 go run main.go session.go 
```
### test in local
```
env SESSION_IP=127.0.0.1:8864 AUTH_IP=127.0.0.1:8901 SHOP_IP=127.0.0.1:8902 ORD_IP=127.0.0.1:8903 PAGE_IP=127.0.0.1:8904 LPTPL_IP=127.0.0.1:8905 LP_IP=127.0.0.1:8906 go test
```
### build & run in docker:
```#bin/bash
docker build -t tidusant/c3madmin-portal .
```
### deploy mongodb replicate :
```#bin/bash
 kubectl apply -f k8s_mongo_devdeploy.yml
 ```
 then use dns-debugging to find out the mongo pod ip
 ```#bin/bash
 kubectl exec -i -t dnsutils -- ping mongo-0.mongo.colis-dev.svc.cluster.local
 kubectl exec -i -t dnsutils -- ping mongo-1.mongo.colis-dev.svc.cluster.local
 kubectl exec -i -t dnsutils -- ping mongo-2.mongo.colis-dev.svc.cluster.local
 ```
 access to mongo service:
 ```#bin/bash
 mongo --host mongodb://172.17.0.13,172.17.0.14,172.17.0.15/ 
 ```
 init the replicate set
 ```#bin/bash
rs.initiate(
    {
       _id: "rs1414",
       version: 1,
       members: [
          { _id: 0, host : "mongo-0.mongo" },
          { _id: 1, host : "mongo-1.mongo" },
          { _id: 2, host : "mongo-2.mongo" }
       ]
    }
 )
```
 Done. Now you can create user or restore your mongo data 



### reference:
https://kubernetes.io/docs/tasks/administer-cluster/dns-debugging-resolution/

### status of service is "terminating" (cannot delete)
You can get rid of this issue by manually editing the pv and then removing the finalizers which looked something like this:

kubectl get pv | tail -n+2 | awk '{print $1}' | xargs -I{} kubectl patch pv {} -p '{"metadata":{"finalizers": null}}'