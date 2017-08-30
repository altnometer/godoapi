package setupk8s

import "github.com/altnometer/godoapi/droplet"

// SSHCommander used to execute bash commands remotely via ssh.
type SSHCommander struct {
	User string
	IP   string
}

// SetUpMaster would setup k8s master.
func SetUpMaster(env, reg string) {
	reqDataPtr := droplet.CreateRequestData
	reqDataPtr.Size = "1gb"
	reqDataPtr.Region = reg
	reqDataPtr.Names = []string{"master-1"}
	reqDataPtr.Tags = []string{"master", env}
	droplet.CreateDroplet(reqDataPtr)
}
