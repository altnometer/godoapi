package volume

import (
	"fmt"
	"os"

	"github.com/altnometer/godoapi/lib/support"
)

var argVolumeFailMsg = fmt.Sprintf("Provide <%s|%s|%s|%s> subcommand, please.",
	support.YellowSp("list"), support.YellowSp("create"), support.YellowSp("detach"), support.YellowSp("delete"))

// ParseArgs handles os.Args and calls relevant functions in the package.
func ParseArgs(args []string) {
	if len(args) < 1 {
		fmt.Println(argVolumeFailMsg)
		os.Exit(1)
	}
	switch args[0] {
	case "list":
		ParseArgsListVol(args[1:])
	case "create":
		ParseArgsCreateVol(args[1:])
	case "detach":
		ParseArgsDetachVol(args[1:])
	case "delete":
		ParseArgsDeleteVol(args[1:])
	default:
		fmt.Print("Incorrect arg: ")
		support.RedBold.Println(args[0])
		fmt.Println(argVolumeFailMsg)
		os.Exit(1)
	}
}

// ParseArgsListVol handles 'volume list' subcommand.
func ParseArgsListVol(args []string) {
	// volCmd := flag.NewFlagSet("list", flag.ExitOnError)
	fmt.Println("'list' subcmd for volumes is not implemented yet")
	return
	// namePtr := volCmd.String("name", "", "-name=<volname>")
	// regPtr := volCmd.String("region", "fra1", "-region=fra1")
	// volCmd.Parse(args)
	// if len(args) < 1 {
	// 	fmt.Println("Provide the args, please.")
	// 	volCmd.PrintDefaults()
	// 	os.Exit(1)
	// }
	// if volCmd.Parsed() {
	// 	if *namePtr == "" {
	// 		volCmd.PrintDefaults()
	// 		os.Exit(1)
	// 	}
	// }
	// support.ValidateRegions(regPtr)
	// fmt.Printf("*namePtr = %+v\n", *namePtr)
	// fmt.Printf("*regPtr = %+v\n", *regPtr)

}

// ParseArgsCreateVol handles 'volume create' subcommand.
func ParseArgsCreateVol(args []string) {
	// volCmd := flag.NewFlagSet("create", flag.ExitOnError)
	fmt.Println("'create' subcmd for volumes is not implemented yet")
	return
}

// ParseArgsDetachVol handles 'volume detach' subcommand.
func ParseArgsDetachVol(args []string) {
	// volCmd := flag.NewFlagSet("detach", flag.ExitOnError)
	fmt.Println("'detach' subcmd for volumes is not implemented yet")
	return
}

// ParseArgsDeleteVol handles 'volume delete' subcommand.
func ParseArgsDeleteVol(args []string) {
	// volCmd := flag.NewFlagSet("delete", flag.ExitOnError)
	fmt.Println("'delete' subcmd for volumes is not implemented yet")
	return
}
