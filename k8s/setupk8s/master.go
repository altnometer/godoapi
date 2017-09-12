package setupk8s

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/altnometer/godoapi/droplet"
	"github.com/altnometer/godoapi/lib/support"
	"github.com/briandowns/spinner"
)

// SetUpMaster would setup k8s master.
func SetUpMaster(env, reg string) (string, string) {
	userName := os.Getenv("DOHostUsername")
	if userName == "" {
		support.YellowLn("You can set env var DOHostUsername!")
		userName = support.GetUserInput("Type in a DOHostUsername: ")
	}
	userPassword := os.Getenv("DOHostUsernamePassword")
	if userPassword == "" {
		support.YellowLn("You can set env var DOHostUsernamePassword!")
		userPassword = support.GetUserInput("Type in a DOHostUsernamePassword: ")
	}
	sshKeyPath := os.Getenv("DOSSHKeyPath")
	if sshKeyPath == "" {
		support.YellowLn("You can set env var DOSSHKeyPath!")
		sshKeyPath = support.GetUserInput("Type in a DOSSHKeyPath: ")
	}
	sshCmdGetToken := []string{
		"sudo",
		"kubeadm", "token", "list", "|", "awk",
		"'NR == 2 { printf $1 }'",
	}
	master1Name := "master-1"
	runningMasters := *droplet.ReturnDropletsByTag("master")
	var token string
	var PublicIP string
	for _, d := range runningMasters {
		if d["Name"] == master1Name {
			PublicIP = d["PublicIP"]
			token = support.FetchSSHOutput("root", PublicIP, sshKeyPath, sshCmdGetToken)
			support.YellowLn("Set env var for k8s token.")
			os.Setenv("K8SToken", token)
			support.RedPf("Droplet with %s name already exist!", master1Name)
			// return d["PublicIP"], token
		}
	}
	if token == "" {
		reqDataPtr := droplet.GetDefaultDropCreateData()
		reqDataPtr.Size = "1gb"
		reqDataPtr.Region = reg
		reqDataPtr.Names = []string{"master-1"}
		reqDataPtr.Tags = []string{"master", env}
		drSpecs := droplet.CreateDroplet(reqDataPtr)
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
				s.Start()
				PublicIP = dData["PublicIP"]
				s.Stop()
				break
			}
		}
	}
	if PublicIP == "" {
		panic("No PublicIP for k8s master host, cannot continue!")
	}
	// execSSH(userName, PublicIP, sshKeyPath)
	// bash master.sh --TARGET_MACHINE_IP 165.227.134.109 --PATH_TO_SSH_PRIV_KEYS ~/.ssh/circleci --USERNAME sally --USER_PASSWORD las
	arg := []string{
		"/home/sam/redmoo/devops/k8s/setupcluster/docean/master.sh",
		"--TARGET_MACHINE_IP",
		PublicIP,
		"--PATH_TO_SSH_PRIV_KEYS",
		sshKeyPath,
		"--USERNAME",
		userName,
		"--USER_PASSWORD",
		userPassword,
	}
	argstr := strings.Join(arg, " ")
	support.YellowLn("Executing ssh from golang with following args: ")
	fmt.Printf("argstr = %+v\n", argstr)

	// cmdOut, err := exec.Command("ssh", arg...).Output()
	cmd := exec.Command("bash", arg...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
	token = support.FetchSSHOutput("root", PublicIP, sshKeyPath, sshCmdGetToken)
	support.YellowLn("Set env var for k8s token.")
	os.Setenv("K8SToken", token)
	return PublicIP, token
}
