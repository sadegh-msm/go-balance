# GO-balancer

`GO-balancer` is a layer 7 load balancer that supports http and https, and it is also a go library that implements `load balancing` algorithms.

It currently supports load balancing algorithms:
* `round-robin`
* `random`
* `ip-hash`
* `least-load`


## Run
`Balancer` needs to configure the `config.yaml` file, see [config.yaml](https://github.com/sadegh-msm/go-balancer/blob/main/config/config.yaml) :

and now, you can execute `balancer`, the balancer will print the configuration details:
```shell                                       
Schema: http
Port: 8089
Health Check: true
Location:
        Route: /
        Proxy Pass: [http://192.168.1.1 http://192.168.1.2:1015 https://192.168.1.2 http://my-server.com]
        Mode: round-robin
```
`balancer` will perform `health check` on all proxy hosts periodically. When the site is unreachable, it will be removed from the balancer automatically . However, `balancer` will still perform `health check` on unreachable sites. When the site is reachable, it will add it to the balancer automatically.

also, each load balancer implements the `balancer.Balancer` interface:
```go
type Balancer interface {
	Add(string)
	Remove(string)
	Balance(string) (string, error)
	Inc(string)
	Done(string)
}
```

