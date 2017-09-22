package admin

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/altnometer/godoapi/droplet"
	"github.com/altnometer/godoapi/ip"
	"github.com/altnometer/godoapi/lib/support"
	"github.com/digitalocean/godo"
)

var argAdminFailMsg = fmt.Sprintf("Provide <%s|%s|%s> subcommand, please.",
	support.YellowSp("list"), support.YellowSp("create"), support.YellowSp("delete"))

// ParseArgs handles os.Args and calls relevant functions in the package.
func ParseArgs(args []string) error {
	if len(args) < 1 {
		fmt.Println(argAdminFailMsg)
		return support.ErrBadArgs
	}
	switch args[0] {
	case "list":
		if err := listAdmin(); err != nil {
			return err
		}
	case "create":
		if err := parseArgsCreateAdmin(args[1:]); err != nil {
			return err
		}
	case "delete":
		if err := parseArgsDeleteAdmin(args[1:]); err != nil {
			return err
		}

	default:
		fmt.Print("Incorrect arg: ")
		support.RedBold.Println(args[0])
		fmt.Println(argAdminFailMsg)
		return support.ErrBadArgs
	}
	return nil
}

func listAdmin() error {
	droplets, err := droplet.ReturnDropletsByTag("admin")
	if err != nil {
		return err
	}
	support.PrintDropData(droplets)
	return nil
}

func parseArgsCreateAdmin(args []string) error {
	subCmd := flag.NewFlagSet("admin", flag.ExitOnError)
	regPtr := subCmd.String("region", "fra1", "-region=fra1")
	envPtr := subCmd.String("env", "dev", "-env=<prod|test|stage|dev>")
	sizePtr := subCmd.String("size", "512mb", "-size=<512mb|1gb|2gb...>")
	userNamePtr := subCmd.String("username", "", "-username=<somename>")
	passwordPtr := subCmd.String("password", "", "-password=<mypassword>")
	sshKeyPathPtr := subCmd.String("sshkeypath", "", "-sshkeypath=</ssh/path/myprivkey>")
	subCmd.Parse(args)
	if subCmd.Parsed() {
		if err := support.ValidateRegions(regPtr); err != nil {
			fmt.Println(err)
			subCmd.PrintDefaults()
			return err
		}
		if *userNamePtr == "" {
			err := errors.New(support.RedSp("no username arg"))
			fmt.Println(err)
			subCmd.PrintDefaults()
			return err
		}
		if *passwordPtr == "" {
			err := errors.New(support.RedSp("no password arg"))
			fmt.Println(err)
			subCmd.PrintDefaults()
			return err
		}
		if *sshKeyPathPtr == "" {
			err := errors.New(support.RedSp("no sshkeypath arg"))
			fmt.Println(err)
			subCmd.PrintDefaults()
			return err

		}
	}
	if len(args) < 1 {
		err := support.ErrBadArgs
		support.Red.Println(err)
		subCmd.PrintDefaults()
		return err
	}
	fmt.Printf("*regPtr = %+v\n", *regPtr)
	createDropData := droplet.GetDefaultDropCreateData()
	createDropData.Size = *sizePtr
	createDropData.Region = *regPtr
	createDropData.Names = []string{"admin"}
	createDropData.Tags = []string{"admin", *envPtr}
	// fmt.Printf("createDropData = %+v\n", createDropData)
	// fmt.Printf("(&multiName).String() = %+v\n", (&multiName).String())
	// fmt.Printf("multiName[0] = %+v\n", multiName[0])
	// fmt.Printf("(&multiTag).String() = %+v\n", (&multiTag).String())
	// // fmt.Printf("*namePtr = %+v\n", *namePtr)
	// fmt.Printf("*regPtr = %+v\n", *regPtr)
	// fmt.Printf("*sizePtr = %+v\n", *sizePtr)
	if err := setupAdmin(createDropData,
		*userNamePtr, *passwordPtr, *sshKeyPathPtr); err != nil {
		return err
	}
	return nil
}
func parseArgsDeleteAdmin(args []string) error {
	log.Println("called delete admin")
	return nil
}

func setupAdmin(
	crData *godo.DropletMultiCreateRequest,
	userName,
	password,
	sshKeyPath string) error {

	// Check if admin server exist.
	droplets, err := droplet.ReturnDropletsByTag("admin")
	var dr godo.Droplet
	if len(droplets) > 0 {
		dr = droplets[0]
	} else {
		drops, err := droplet.CreateDroplet(crData)
		if err != nil {
			return err
		}
		dr := drops[0]
		fmt.Println("Wait for the droplet to boot up...")
		time.Sleep(10 * time.Second)
		dr = droplet.ReturnDropletByID(dr.ID)
	}
	publicIP, err := dr.PublicIPv4()
	if err != nil {
		return err
	}
	cmdOpts := []string{
		"--TARGET_MACHINE_IP",
		publicIP,
		"--PATH_TO_SSH_PRIV_KEYS",
		sshKeyPath,
		"--USERNAME",
		userName,
		"--USER_PASSWORD",
		password,
	}
	if sshRootClient, err := support.GetSSHClient("root", publicIP); err == nil {
		defer sshRootClient.Close()
		sshSession, err := support.GetSSHInterSession(sshRootClient)
		if err != nil {
			return err
		}
		scriptPath := "/home/sam/redmoo/devops/k8s/setupcluster/docean/admin-1.sh"
		support.YellowPf("executing %s\n", scriptPath)
		// args := append([]string{"bash", scriptPath}, cmdOpts...)
		// if err := support.ExecCmd(args); err != nil {
		// 	return err
		// }
		args := append([]string{
			"scp", "-o", "UserKnownHostsFile=/dev/null", "-o", "StrictHostKeyChecking=no",
			"-i", sshKeyPath, scriptPath,
			"root" + "@" + publicIP + ":/root",
			// " sudo -E bash -c \"tar -C /etc -xzf -\"",
		})
		if err := support.ExecCmd(args); err != nil {
			return err
		}
		scriptName := filepath.Base(scriptPath)
		cmd := fmt.Sprintf("/bin/bash ./%s %s", scriptName, strings.Join(cmdOpts, " "))
		if err = sshSession.Run(cmd); err != nil {
			return err
		}
		// sshSession.Close()
	}
	scriptPath := "/home/sam/redmoo/devops/k8s/setupcluster/docean/admin-2.sh"
	args := append([]string{
		"scp", "-o", "UserKnownHostsFile=/dev/null", "-o", "StrictHostKeyChecking=no",
		"-i", sshKeyPath, scriptPath,
		userName + "@" + publicIP + ":/home/" + userName,
	})
	if err := support.ExecCmd(args); err != nil {
		return err
	}
	scriptName := filepath.Base(scriptPath)
	cmd := fmt.Sprintf("sudo -E bash ./%s %s", scriptName, strings.Join(cmdOpts, " "))
	support.YellowPf("executing %s\n", cmd)
	sshClient, err := support.GetSSHClient(userName, publicIP)
	if err != nil {
		return err
	}
	defer sshClient.Close()
	sshSession, err := support.GetSSHInterSession(sshClient)
	if err != nil {
		return err
	}
	if err = sshSession.Run(cmd); err != nil {
		return err
	}
	sshSession.Close()
	// args = append([]string{"bash", scriptPath}, cmdOpts...)
	// if err := support.ExecCmd(args); err != nil {
	// 	return err
	// }
	// Enable serving 'site unavailable page'
	// func UserConfirmDefaultN(prompt string) (conf bool, err error) {
	promtMsg := "Set up nginx serving 'site unavailable page'?"
	confirmed, err := support.UserConfirmDefaultN(promtMsg)
	if confirmed {
		support.YellowLn("Copy letsencrypt files ...")
		args = append([]string{
			"scp", "-o", "UserKnownHostsFile=/dev/null", "-o", "StrictHostKeyChecking=no",
			"-i", sshKeyPath, support.TSLArchSource,
			userName + "@" + publicIP + ":/home/" + userName,
			// " sudo -E bash -c \"tar -C /etc -xzf -\"",
		})
		if err := support.ExecCmd(args); err != nil {
			return err
		}
		sshSession, err := support.GetSSHInterSession(sshClient)
		if err != nil {
			return err
		}
		archName := filepath.Base(support.TSLArchSource)
		cmd := fmt.Sprintf("sudo tar -C /etc -xvzf %s", archName)
		cmd = cmd + " && sleep 1 && " + "docker pull redmoo/unavailable:latest"
		// cmd = cmd + " && sleep 1 && " + "docker rm $(docker stop $(docker -a -q --filter ancestor=redmoo/unavailable))"
		cmd = cmd + " && sleep 2 && " + "docker run -it -v /etc/letsencrypt/:/etc/letsencrypt/" +
			" -p 443:443 -p 80:80 -p 10254:10254 -d redmoo/unavailable"
		support.YellowPf("executing %s\n", cmd)
		if err = sshSession.Run(cmd); err != nil {
			return err
		}
		sshSession.Close()
		if err := ip.AssignIP(); err != nil {
			return err
		}
	}
	return nil
}
