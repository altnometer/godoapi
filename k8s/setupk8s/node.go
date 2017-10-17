package setupk8s

import (
	"fmt"
	"path/filepath"
	"strconv"
	"time"

	"github.com/altnometer/godoapi/droplet"
	"github.com/altnometer/godoapi/lib/support"
	"github.com/digitalocean/godo"
)

func addNode(env, reg, size string) error {
	runningMasters, err := droplet.ReturnDropletsByTag("master")
	if err != nil {
		return err
	}
	var token string
	var masPubIP string
	var masPrivIP string
	var nodePubIP string
	sshCmdGetToken := []string{
		"sudo",
		"kubeadm", "token", "list", "|", "awk",
		"'NR == 2 { printf $1 }'",
	}
	for _, d := range runningMasters {
		if d.Name == support.Master1Name &&
			d.Region.Slug == reg && support.StringInSlice(env, d.Tags) {
			if masPubIP, err = d.PublicIPv4(); err != nil {
				return err
			}
			if masPrivIP, err = d.PrivateIPv4(); err != nil {
				return err
			}
			support.RedLn("Only SINGLE master is implemented currently!!!")
			token = support.GetSSHOutput("root", masPubIP, sshCmdGetToken)
			if token == "" {
				return fmt.Errorf("no token has been fetched")
			}
			break
		}
	}
	if token == "" {
		return fmt.Errorf("no masters in region %s for env %s", reg, env)
	}
	runningNodes, err := droplet.ReturnDropletsByTag("node")
	if err != nil {
		return err
	}
	if len(runningNodes) > 0 {
		fmt.Println("Running nodes:")
		for i, d := range runningNodes {
			fmt.Printf("    %d. %+v\n", i, support.YellowSp(d.Name))
		}
		confirmed, err := support.UserConfirmDefaultN("Redeploy any nodes?")
		if err != nil {
			return err
		}
		if confirmed {
			drNum, err := strconv.ParseInt(
				support.GetUserInput("Select host # if you wish to redeploy node: "), 10, 0)
			if err != nil {
				return err
			}
			if drNum > 0 || drNum < int64(len(runningNodes)) {
				nodeRedeploy := runningNodes[drNum]
				if nodePubIP, err = nodeRedeploy.PublicIPv4(); err != nil {
					return err
				}
				if err := execNodeSetupCmds(nodePubIP, masPrivIP, token); err != nil {
					return err
				}
				return nil
			}
		}
	}
	name := fmt.Sprintf("node-%02d", len(runningNodes)+1)
	node, err := createNode(name, env, reg, size)
	if err != nil {
		return err
	}
	if nodePubIP, err = node.PublicIPv4(); err != nil {
		return err
	}
	if err := execNodeSetupCmds(nodePubIP, masPrivIP, token); err != nil {
		return err
	}
	return nil
}

// SetUpNode would setup k8s node.
func SetUpNode(env, reg, masterIP, token string) error {
	runningNodes, err := droplet.ReturnDropletsByTag("node")
	if err != nil {
		return err
	}
	var publicIP string
	for _, d := range runningNodes {
		if publicIP, err = d.PublicIPv4(); err != nil {
			return err
		}
	}
	if publicIP == "" {
		dData, err := createNode("node-01", env, reg, support.NodeMemSize)
		if err != nil {
			return err
		}
		if publicIP, err = dData.PublicIPv4(); err != nil {
			return err
		}
	} else {
		prompt := "Node(s) exist. Do you wish to continue?"
		confirmed, err := support.UserConfirmDefaultN(prompt)
		if err != nil {
			return err
		}
		if !confirmed {
			return fmt.Errorf("execution is not confirmed by user")
		}
	}
	if err := execNodeSetupCmds(publicIP, masterIP, token); err != nil {
		return err
	}
	return nil
}

func createNode(name, env, reg, size string) (*godo.Droplet, error) {
	reqDataPtr := droplet.GetDefaultDropCreateData()
	reqDataPtr.Size = size
	reqDataPtr.Region = reg
	reqDataPtr.Names = []string{name}
	reqDataPtr.Tags = []string{"k8s", "node", env}
	drops, err := droplet.CreateDroplet(reqDataPtr)
	if err != nil {
		return nil, err
	}
	if drops == nil {
		return nil, fmt.Errorf("no droplet(s) created")
	}
	fmt.Println("Waiting for the droplet to initialize ...")
	time.Sleep(10 * time.Second)
	dData := droplet.ReturnDropletByID(drops[0].ID)
	return &dData, nil
}

func execNodeSetupCmds(pubIP, masterIP, token string) error {
	scriptPath := "/home/sam/redmoo/devops/k8s/setupcluster/docean/node-1.sh"
	args := append([]string{
		"scp", "-o", "UserKnownHostsFile=/dev/null", "-o", "StrictHostKeyChecking=no",
		scriptPath,
		"root" + "@" + pubIP + ":/root",
	})
	if err := support.ExecCmd(args); err != nil {
		return err
	}
	cmd := fmt.Sprintf(
		"/bin/bash ./%s",
		filepath.Base(scriptPath))
	sshRootClient, err := support.GetSSHClient("root", pubIP)
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
	userDataPart2 := fmt.Sprintf(
		"\nkubeadm join --token %s %s:6443", token, masterIP)
	if err = sshSession.Run(userDataPart2); err != nil {
		return err
	}
	return nil

}
