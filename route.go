package main

import (
	"flag"
	"fmt"
	"kubernetes-haproxy-autolb/backends/etcd3client"
	"kubernetes-haproxy-autolb/backends/node"
	"kubernetes-haproxy-autolb/backends/watch"
	"strings"
)

// var (
// 	endpoints = []string{"10.1.10.201:2379"}
// 	serviceip = "192.168.110.0/24"
// )

func main() {
	enps := flag.String("Endpoints", "10.1.10.201:2379", "etcdserverip eg:--endpoints=10.1.10.201,10.1.10.202:2379 ")
	// sip := flag.String("Serviceip", "192.168.110.0/24", "eg:  --Serviceip=192.168.110.0/24")
	flag.Parse()
	endpoints := strings.SplitN(*enps, ",", -1)
	// serviceip := *sip
	ch := make(chan string)

	f := etcd3client.Node{endpoints, "/autohaproxy/node/nodeip/"}

	go watch.Nodeiproutewatch("/autohaproxy/node/nodeip/", endpoints, f, ch)

	node.Iproute(f, endpoints)
	for {

		fmt.Println(<-ch)

	}
}
