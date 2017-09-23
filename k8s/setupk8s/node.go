package setupk8s

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/altnometer/godoapi/droplet"
	"github.com/altnometer/godoapi/lib/support"
)

// SetUpNode would setup k8s node.
func SetUpNode(env, reg, ip, token string) error {
	runningNodes, err := droplet.ReturnDropletsByTag("node")
	var publicIP string
	for _, d := range runningNodes {
		if publicIP, err = d.PublicIPv4(); err != nil {
			return err
		}
		fmt.Printf("publicIP = %+v\n", publicIP)
	}
	if publicIP == "" {
		reqDataPtr := droplet.GetDefaultDropCreateData()
		reqDataPtr.Size = "1gb"
		reqDataPtr.Region = reg
		reqDataPtr.Names = []string{"node-1"}
		reqDataPtr.Tags = []string{"node", env}
		drSpecs, err := droplet.CreateDroplet(reqDataPtr)
		if err != nil {
			return err
		}
		fmt.Println("Wait for the droplet to initialize ...")
		time.Sleep(10 * time.Second)
		dData := droplet.ReturnDropletByID(drSpecs[0].ID)
		if publicIP, err = dData.PublicIPv4(); err != nil {
			return err
		}
	}
	scriptPath := "/home/sam/redmoo/devops/k8s/setupcluster/docean/node-1.sh"
	args := append([]string{
		"scp", "-o", "UserKnownHostsFile=/dev/null", "-o", "StrictHostKeyChecking=no",
		scriptPath,
		"root" + "@" + publicIP + ":/root",
	})
	if err := support.ExecCmd(args); err != nil {
		return err
	}
	cmd := fmt.Sprintf(
		"/bin/bash ./%s",
		filepath.Base(scriptPath))
	sshRootClient, err := support.GetSSHClient("root", publicIP)
	if err != nil {
		return err
	}
	defer sshRootClient.Close()
	sshSession, err := support.GetSSHInterSession(sshRootClient)
	if err != nil {
		return err
	}
	if err = sshSession.Run(cmd); err != nil {
		return err
	}
	sshSession.Close()
	time.Sleep(2 * time.Second)
	sshSession, err = support.GetSSHInterSession(sshRootClient)
	userDataPart2 := fmt.Sprintf("\nkubeadm join --token %s %s:6443", token, ip)
	if err = sshSession.Run(userDataPart2); err != nil {
		return err
	}
	return nil
}
