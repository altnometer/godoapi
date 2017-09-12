package droplet

import (
	"fmt"
	"os"

	"github.com/altnometer/godoapi/lib/support"
)

var argDropletFailMsg = fmt.Sprintf("Provide <%s|%s|%s> subcommand, please.",
	support.YellowSp("list"), support.YellowSp("create"), support.YellowSp("delete"))

// ParseArgs handles os.Args and calls relevant functions in the package.
func ParseArgs(args []string) error {
	if len(args) < 1 {
		fmt.Println(argDropletFailMsg)
		os.Exit(1)
	}
	switch args[0] {
	case "list":
		ParseArgsListDrop(args[1:])
	case "create":
		if err := ParseArgsCreateDrop(args[1:]); err != nil {
			return err
		}
	case "delete":
		if err := ParseArgsDeleteDrop(args[1:]); err != nil {
			return err
		}

	default:
		fmt.Print("Incorrect arg: ")
		support.RedBold.Println(args[0])
		fmt.Println(argDropletFailMsg)
		os.Exit(1)
	}
	return nil
}
