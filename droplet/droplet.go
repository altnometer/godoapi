package droplet

import (
	"flag"
	"fmt"
	"os"
)

// ParseArgs handles os.Args and calls relevant functions in the package.
func ParseArgs(args []string) {
	dropCmd := flag.NewFlagSet("droplet", flag.ExitOnError)
	// droplet  subcommand flag pointers
	NamePtr := dropCmd.String("name", "", "-name=<volname1[,volname2...]>")
	RegPtr := dropCmd.String("region", "fra1", "-region=fra1")
	dropCmd.Parse(args)
	if len(args) < 1 {
		fmt.Println("Provide the args, please.")
		dropCmd.PrintDefaults()
		os.Exit(1)
	}
	if dropCmd.Parsed() {
		if *NamePtr == "" {
			dropCmd.PrintDefaults()
			os.Exit(1)
		}
	}
	fmt.Printf("*NamePtr = %+v\n", *NamePtr)
	fmt.Printf("*RegPtr = %+v\n", *RegPtr)

}
