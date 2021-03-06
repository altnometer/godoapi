package setupk8s

import (
	"errors"
	"flag"
	"fmt"

	"github.com/altnometer/godoapi/droplet"
	"github.com/altnometer/godoapi/lib/support"
)

var argK8SFailMsg = fmt.Sprintf("Provide <%s|%s|%s> subcommand, please.",
	support.YellowSp("create"),
	support.YellowSp("add-node"),
	support.YellowSp("delete"))

// ParseArgs handles os.Args and calls relevant functions in the package.
func ParseArgs(args []string) error {
	if len(args) < 1 {
		fmt.Println(argK8SFailMsg)
		return support.ErrBadArgs
	}
	switch args[0] {
	case "create":
		if err := parseArgsSetupK8S(args[1:]); err != nil {
			return err
		}
	case "delete":
		fmt.Println("Not implemented yet.")
		return nil
	case "add-node":
		if err := parseAddNodeArgs(args[1:]); err != nil {
			return err
		}
	default:
		fmt.Print("Incorrect arg: ")
		support.RedBold.Println(args[0])
		fmt.Println(argK8SFailMsg)
		return support.ErrBadArgs
	}
	return nil
}

func parseAddNodeArgs(args []string) error {
	subCmd := flag.NewFlagSet("add-node", flag.ExitOnError)
	envPtr := subCmd.String("env", "dev", "-env=<prod|test|stage|dev>")
	regPtr := subCmd.String("region", "fra1", "-region=fra1")
	sizePtr := subCmd.String("size", "1gb", "-size=<512mb|1gb|2gb...>")
	subCmd.Parse(args)
	if subCmd.Parsed() {
		if err := support.ValidateRegions(regPtr); err != nil {
			fmt.Println(err)
			subCmd.PrintDefaults()
			return err
		}
	}
	// if len(args) < 1 {
	// 	err := support.ErrBadArgs
	// 	support.Red.Println(err)
	// 	subCmd.PrintDefaults()
	// 	return err
	// }
	if err := addNode(*envPtr, *regPtr, *sizePtr); err != nil {
		return err
	}
	return nil
}
func parseArgsSetupK8S(args []string) error {
	subCmd := flag.NewFlagSet("setupk8s", flag.ExitOnError)
	envPtr := subCmd.String("env", "dev", "-env=<prod|test|stage|dev>")
	regPtr := subCmd.String("region", "fra1", "-region=fra1")
	sizePtr := subCmd.String("size", "1gb", "-size=<512mb|1gb|2gb...>")
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
	createDropData := droplet.GetDefaultDropCreateData()
	createDropData.Size = *sizePtr
	createDropData.Region = *regPtr
	createDropData.Names = []string{"master-1"}
	createDropData.Tags = []string{"k8s", "master", *envPtr}
	ip, token, err := SetUpMaster(
		createDropData,
		*userNamePtr,
		*passwordPtr,
		*sshKeyPathPtr,
	)
	if err != nil {
		return err
	}
	if err := SetUpNode(*envPtr, *regPtr, ip, token); err != nil {
		return err
	}
	return nil
}
