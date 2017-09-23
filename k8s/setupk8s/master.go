package setupk8s

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/altnometer/godoapi/droplet"
	"github.com/altnometer/godoapi/lib/support"
	"github.com/digitalocean/godo"
)

// SetUpMaster would setup k8s master.
func SetUpMaster(
	crData *godo.DropletMultiCreateRequest,
	userName,
	password,
	sshKeyPath string) (string, string, error) {
	sshCmdGetToken := []string{
		"sudo",
		"kubeadm", "token", "list", "|", "awk",
		"'NR == 2 { printf $1 }'",
	}
	runningMasters, err := droplet.ReturnDropletsByTag("master")
	if err != nil {
		return "", "", err
	}
	var token string
	var publicIP string
	var privateIP string
	for _, d := range runningMasters {
		if d.Name == support.Master1Name {
			if publicIP, err = d.PublicIPv4(); err != nil {
				return "", "", err
			}
			if privateIP, err = d.PrivateIPv4(); err != nil {
				return "", "", err
			}
			support.RedPf("Droplet with %s name already exist!\n", support.Master1Name)
			break
		}
	}
	if publicIP == "" || privateIP == "" {
		drSpecs, err := droplet.CreateDroplet(crData)
		if err != nil {
			return "", "", err
		}
		if drSpecs != nil {
			support.YellowLn("Initializing the droplet ...")
			time.Sleep(time.Second * 10)
		}
		for _, d := range drSpecs {
			if d.Name != "" {
				support.RedLn("Only SINGLE master is handled currently!!!")
				dData := droplet.ReturnDropletByID(d.ID)
				if publicIP, err = dData.PublicIPv4(); err != nil {
					return "", "", err
				}
				if privateIP, err = dData.PrivateIPv4(); err != nil {
					return "", "", err
				}
				break
			}
		}
	}
	if publicIP == "" || privateIP == "" {
		return "", "", fmt.Errorf("no publicIP or privateIP for k8s master host, cannot continue")
	}
	scriptPath := "/home/sam/redmoo/devops/k8s/setupcluster/docean/master-1.sh"
	cmdOpts := []string{
		"--TARGET_MACHINE_IP", publicIP,
		"--PATH_TO_SSH_PRIV_KEYS", sshKeyPath,
		"--USERNAME", userName,
		"--USER_PASSWORD", password,
	}
	args := append([]string{
		"scp", "-o", "UserKnownHostsFile=/dev/null", "-o", "StrictHostKeyChecking=no",
		"-i", sshKeyPath, scriptPath,
		"root" + "@" + publicIP + ":/root",
	})
	if err := support.ExecCmd(args); err != nil {
		return "", "", err
	}
	cmd := fmt.Sprintf(
		"/bin/bash ./%s %s",
		filepath.Base(scriptPath), strings.Join(cmdOpts, " "))
	sshRootClient, err := support.GetSSHClient("root", publicIP)
	if err != nil {
		return "", "", err
	}
	defer sshRootClient.Close()
	sshSession, err := support.GetSSHInterSession(sshRootClient)
	if err != nil {
		return "", "", err
	}
	if err = sshSession.Run(cmd); err != nil {
		return "", "", err
	}
	scriptPath = "/home/sam/redmoo/devops/k8s/setupcluster/docean/master-2.sh"
	args = append([]string{"sudo", "-E", "bash", scriptPath}, cmdOpts...)
	if err := support.ExecCmd(args); err != nil {
		return "", "", err
	}
	token = support.FetchSSHOutput("root", publicIP, sshKeyPath, sshCmdGetToken)
	// return publicIP, token, nil
	return privateIP, token, nil
}
