# Service

TODO:
https://go-zero.dev/docs/advance/rpc-call

## Usage

```go
type Config struct {
service.ServerConfiguration
Mysql struct {
DataSource string
}
CacheRedis cache.CacheConf
}
```
```yaml
Name: user.rpc
ListenOn: 127.0.0.1:8080
Etcd:
  Hosts:
    - $etcdHost
  Key: user.rpc
Mysql:
  DataSource: $user:$password@tcp($url)/$db?charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai
CacheRedis:
  - Host: $host
    Pass: $pass
    Type: node 
```
