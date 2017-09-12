package admin

import (
	"flag"
	"fmt"
	"log"

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
		if err := parseArgsListAdmin(args[1:]); err != nil {
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

func parseArgsListAdmin(args []string) error {
	log.Println("called list admin")
	subCmd := flag.NewFlagSet("admin", flag.ExitOnError)
	// regPtr := subCmd.String("region", "fra1", "-region=fra1")
	subCmd.Parse(args)
	if subCmd.Parsed() {
		// support.ValidateRegions(regPtr)
	}
	// if len(args) < 1 {
	// 	support.Red.Println("Provide the args, please.")
	// 	subCmd.PrintDefaults()
	// 	os.Exit(1)
	// }
	// fmt.Printf("*regPtr = %+v\n", *regPtr)
	return nil
}

func parseArgsCreateAdmin(args []string) error {
	log.Println("called create admin")
	return nil
}
func parseArgsDeleteAdmin(args []string) error {
	log.Println("called delete admin")
	return nil
}
