package admin

import (
	"flag"
	"fmt"
	"log"

	"github.com/altnometer/godoapi/droplet"
	"github.com/altnometer/godoapi/lib/support"
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
	subCmd.Parse(args)
	if subCmd.Parsed() {
		if err := support.ValidateRegions(regPtr); err != nil {
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
	droplet.CreateDroplet(createDropData)
	return nil
}
func parseArgsDeleteAdmin(args []string) error {
	log.Println("called delete admin")
	return nil
}
