package admin

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/altnometer/godoapi/droplet"
	"github.com/altnometer/godoapi/lib/support"
	"github.com/digitalocean/godo"
)

var argDropletFailMsg = fmt.Sprintf("Provide <%s|%s|%s> subcommand, please.",
	support.YellowSp("list"), support.YellowSp("create"), support.YellowSp("delete"))

// ParseArgs handles os.Args and calls relevant functions in the package.
func ParseArgs(args []string) error {
	if len(args) < 1 {
		fmt.Println(argDropletFailMsg)
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
		fmt.Println(argDropletFailMsg)
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
	createDropData.Tags = []string{"admin"}
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
		dr = droplet.CreateDroplet(crData)[0]
		fmt.Println("Wait for the droplet to boot up...")
		time.Sleep(3 * time.Second)
		dr = droplet.ReturnDropletByID(dr.ID)
	}
	publicIP, err := dr.PublicIPv4()
	if err != nil {
		return err
	}
	args := []string{
		"/home/sam/redmoo/devops/k8s/setupcluster/docean/admin-1.sh",
		"--TARGET_MACHINE_IP",
		publicIP,
		"--PATH_TO_SSH_PRIV_KEYS",
		sshKeyPath,
		"--USERNAME",
		userName,
		"--USER_PASSWORD",
		password,
	}
	if err := os.Setenv("DOSSHKeyPath", sshKeyPath); err != nil {
		return err
	}
	// argstr := strings.Join(arg, " ")
	// support.YellowLn("Executing ssh from golang with following args: ")
	// fmt.Printf("argstr = %+v\n", argstr)

	// cmdOut, err := exec.Command("ssh", arg...).Output()
	cmd := exec.Command("bash", args...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return err
	}
	// TODO: scp vimsetup and install vim
	// TODO: scp bashsetup and configure
	// support.ExecSSH("root", publicIP, args)
	// if err != nil {
	// 	return err
	// }
	return nil
}
