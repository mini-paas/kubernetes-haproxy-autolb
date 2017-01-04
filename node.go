package main

import (
	"fmt"
	"github.com/byebye758/kubernetes-haproxy-autolb/backends/con"
	"github.com/byebye758/kubernetes-haproxy-autolb/backends/etcd3client"
	"github.com/byebye758/kubernetes-haproxy-autolb/backends/node"
	"github.com/byebye758/kubernetes-haproxy-autolb/backends/watch"
	//"time"
)

var (
	endpoints = []string{"10.1.10.201:2379"}
	serviceip = "192.168.110.0/24"
)

func main() {
	ch := make(chan string)
	nodeip := con.HostIP()
	dockerip, _ := con.Getdockerip()
	a := etcd3client.Autotable{endpoints, "/autohaproxy/autotable/"}

	f := etcd3client.Node{endpoints, "/autohaproxy/node/nodeip/"}

	g := etcd3client.NodeRegister{
		endpoints,
		"/autohaproxy/node/nodeip/" + nodeip,
		nodeip,
		dockerip,
		ch,
	}

	go node.Noderegister(g, ch)
	node.Iproute(f, endpoints)
	// fmt.Println("2")
	node.Noderoute(a, endpoints)
	// fmt.Println("3")
	node.Serviceiproute(serviceip)

	go watch.Nodeiproutewatch("/autohaproxy/node/nodeip/", endpoints, f, ch)
	go watch.Nodenoderoutewatch("/autohaproxy/autotable/", endpoints, a, ch)
	// go node.Iproute(f, endpoints)
	// go node.Noderoute(a, endpoints)
	for {

		fmt.Println(<-ch)

	}
}