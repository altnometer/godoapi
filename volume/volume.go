package volume

import (
	"flag"
	"fmt"
	"os"

	"github.com/altnometer/godoapi/lib/support"
)

// ParseArgs handles os.Args and calls relevant functions in the package.
func ParseArgs(args []string) {
	volCmd := flag.NewFlagSet("volume", flag.ExitOnError)
	// Volume  subcommand flag pointers
	NamePtr := volCmd.String("name", "", "-name=<volname>")
	regPtr := volCmd.String("region", "fra1", "-region=fra1")
	volCmd.Parse(args)
	if len(args) < 1 {
		fmt.Println("Provide the args, please.")
		volCmd.PrintDefaults()
		os.Exit(1)
	}
	if volCmd.Parsed() {
		if *NamePtr == "" {
			volCmd.PrintDefaults()
			os.Exit(1)
		}
	}
	support.ValidateRegions(regPtr)
	fmt.Printf("*NamePtr = %+v\n", *NamePtr)
	fmt.Printf("*regPtr = %+v\n", *regPtr)

}
