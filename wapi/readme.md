#api for web
test

### run in local:
```
env SESSION_IP=127.0.0.1:8865 AUTH_IP=127.0.0.1:8901 SHOP_IP=127.0.0.1:8902 ORD_IP=127.0.0.1:8903 LPL_IP=127.0.0.1:8907 go run main.go session.go 
```
### test in local
```
env SESSION_IP=127.0.0.1:8865 AUTH_IP=127.0.0.1:8901 SHOP_IP=127.0.0.1:8902 ORD_IP=127.0.0.1:8903 LPL_IP=127.0.0.1:8907 go test
```