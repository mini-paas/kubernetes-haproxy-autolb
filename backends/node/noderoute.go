package node

import (
	"fmt"
	"github.com/byebye758/kubernetes-haproxy-autolb/backends/con"
	"github.com/byebye758/kubernetes-haproxy-autolb/backends/etcd3client"
	"github.com/byebye758/kubernetes-haproxy-autolb/backends/node/cmd"
	"strings"
)

// var (
// 	endpoints = []string{"10.1.10.201:2379"}
// )

func Noderoute(a etcd3client.AGetr, endpoints []string) {
	nodemap := make(map[string]map[string]string) //etcd get  autotable   format map
	//a := etcd3client.Autotable{endpoints, "/autohaproxy/autotable/"}
	b, _ := cmd.NoderuleGet()
	autotable := a.AGet()
	nodeip := con.HostIP()

	fmt.Println(autotable, b)
	for _, v := range autotable {
		v := v.(map[string]interface{})
		etcdnodeip := v["Nodeip"].(string)
		etcdhaproxyip := v["Haproxyip"].(string)
		etcdhaproxytable := v["Haproxytable"].(string)
		etcdpodip := v["Podip"].(string)
		if strings.EqualFold(nodeip, etcdnodeip) {
			test := map[string]string{
				"Nodeip":       etcdnodeip,
				"Haproxyip":    etcdhaproxyip,
				"Haproxytable": etcdhaproxytable,
				"Podip":        etcdpodip,
			}
			nodemap[etcdpodip] = test

		}

	}
	fmt.Println(nodemap)

	for _, v := range b {

		if _, ok := nodemap[v["Podip"]]; ok {

		} else {
			cmd.Routetablecmd("ip rule del from "+v["Nodeip"], "")

		}

	}
	b, _ = cmd.NoderuleGet()
	noderulemap := make(map[string]map[string]string)
	for _, v := range b {

		noderulemap[v["Podip"]] = v

	}

	for _, v := range nodemap {
		podip := v["Podip"]
		haproxyip := v["Haproxyip"]
		haproxytable := v["Haproxytable"]

		if _, ok := noderulemap[podip]; ok {

		} else {
			cmd.Routetablecmd("ip rule add from "+podip+"/32 pref 30000 table ", haproxytable)
			cmd.Routetablecmd("ip route replace default via "+haproxyip+" table ", haproxytable)
		}
	}

}