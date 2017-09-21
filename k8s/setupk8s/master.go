package setupk8s

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/altnometer/godoapi/droplet"
	"github.com/altnometer/godoapi/lib/support"
	"github.com/briandowns/spinner"
	"github.com/digitalocean/godo"
)

// SetUpMaster would setup k8s master.
func SetUpMaster(
	crData *godo.DropletMultiCreateRequest,
	userName,
	password,
	sshKeyPath string) (string, string) {
	sshCmdGetToken := []string{
		"sudo",
		"kubeadm", "token", "list", "|", "awk",
		"'NR == 2 { printf $1 }'",
	}
	runningMasters, err := droplet.ReturnDropletsByTag("master")
	if err != nil {
		log.Fatal(err)
	}
	var token string
	var publicIP string
	for _, d := range runningMasters {
		if d.Name == support.Master1Name {
			publicIP, err := d.PublicIPv4()
			if err != nil {
				log.Fatal(err)
			}
			token = support.FetchSSHOutput("root", publicIP, sshKeyPath, sshCmdGetToken)
			support.YellowLn("Set env var for k8s token.")
			os.Setenv("K8SToken", token)
			support.RedPf("Droplet with %s name already exist!", support.Master1Name)
			// return d["publicIP"], token
		}
	}
	if token == "" {
		drSpecs := droplet.CreateDroplet(crData)
		s := spinner.New(spinner.CharSets[9], 150*time.Millisecond)
		s.Start()
		// Give it some time for IPs to be assigned to the droplets.
		support.YellowLn("Initializing the droplet ...")
		time.Sleep(time.Second * 10)
		s.Stop()
		for _, d := range drSpecs {
			if d.Name != "" {
				support.RedLn("Only SINGLE master is handled currently!!!")
				dData := droplet.ReturnDropletByID(d.ID)
				publicIP, err = dData.PublicIPv4()
				if err != nil {
					log.Fatal(err)
				}
				break
			}
		}
	}
	if publicIP == "" {
		panic("No publicIP for k8s master host, cannot continue!")
	}
	// execSSH(userName, publicIP, sshKeyPath)
	// bash master.sh --TARGET_MACHINE_IP 165.227.134.109 --PATH_TO_SSH_PRIV_KEYS ~/.ssh/circleci --USERNAME sally --USER_PASSWORD las
	arg := []string{
		"/home/sam/redmoo/devops/k8s/setupcluster/docean/master.sh",
		"--TARGET_MACHINE_IP",
		publicIP,
		"--PATH_TO_SSH_PRIV_KEYS",
		sshKeyPath,
		"--USERNAME",
		userName,
		"--USER_PASSWORD",
		password,
	}
	argstr := strings.Join(arg, " ")
	support.YellowLn("Executing ssh from golang with following args: ")
	fmt.Printf("argstr = %+v\n", argstr)

	// cmdOut, err := exec.Command("ssh", arg...).Output()
	cmd := exec.Command("bash", arg...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		panic(err)
	}
	token = support.FetchSSHOutput("root", publicIP, sshKeyPath, sshCmdGetToken)
	support.YellowLn("Set env var for k8s token.")
	os.Setenv("K8SToken", token)
	return publicIP, token
}
