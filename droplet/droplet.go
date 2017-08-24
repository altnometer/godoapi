package droplet

import (
	"flag"
	"fmt"
	"os"

	"github.com/altnometer/godoapi/lib/support"
)

// ParseArgs handles os.Args and calls relevant functions in the package.
func ParseArgs(args []string) {
	dropCmd := flag.NewFlagSet("droplet", flag.ExitOnError)
	// droplet  subcommand flag pointers
	namePtr := dropCmd.String("name", "", "-name=<volname1[,volname2...]>")
	regPtr := dropCmd.String("region", "fra1", "-region=fra1")
	dropCmd.Parse(args)
	if len(args) < 1 {
		fmt.Println("Provide the args, please.")
		dropCmd.PrintDefaults()
		os.Exit(1)
	}
	if dropCmd.Parsed() {
		if *namePtr == "" {
			dropCmd.PrintDefaults()
			os.Exit(1)
		}
	}
	support.ValidateRegions(regPtr)
	fmt.Printf("*namePtr = %+v\n", *namePtr)
	fmt.Printf("*regPtr = %+v\n", *regPtr)

}
