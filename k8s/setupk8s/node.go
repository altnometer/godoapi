package setupk8s

import (
	"fmt"
	"os"

	"github.com/altnometer/godoapi/droplet"
)

var userDataPart1 = `#! /bin/bash

apt-get update && apt-get upgrade -y
curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key add -
cat <<EOF2 > /etc/apt/sources.list.d/kubernetes.list
deb http://apt.kubernetes.io/ kubernetes-xenial main
EOF2
apt-get update -y
apt-get install -y docker.io
apt-get install -y --allow-unauthenticated kubelet kubeadm=1.7.0-00 kubectl kubernetes-cni`

// SetUpNode would setup k8s master.
func SetUpNode(env, reg, ip, token string) {
	userDataPart2 := fmt.Sprintf("\nkubeadm join --token %s %s:6443", token, ip)
	userData := fmt.Sprintln(userDataPart1, userDataPart2)
	reqDataPtr := droplet.GetDefaultDropCreateData()
	reqDataPtr.Size = "1gb"
	reqDataPtr.Region = reg
	reqDataPtr.Names = []string{"node-1"}
	reqDataPtr.Tags = []string{"node", env}
	reqDataPtr.UserData = userData
	fmt.Printf("reqDataPtr = %+v\n", reqDataPtr)
	drSpecs := droplet.CreateDroplet(reqDataPtr)
	fmt.Printf("drSpecs = %+v\n", drSpecs)
	os.Exit(0)
}
