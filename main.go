// Deploy k8s on digitalocean.
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/altnometer/godoapi/droplet"
	"github.com/altnometer/godoapi/k8s/setupk8s"
	"github.com/altnometer/godoapi/lib/support"
	"github.com/altnometer/godoapi/volume"
)

var argEntryFailMsg = fmt.Sprintf("Provide <%s|%s|%s> subcommand, please.",
	support.YellowSp("droplet"), support.YellowSp("volume"), support.YellowSp("setupk8s"))

func main() {
	if len(os.Args) < 2 {
		fmt.Println(argEntryFailMsg)
		os.Exit(1)
	}
	switch os.Args[1] {
	case "droplet":
		if err := droplet.ParseArgs(os.Args[2:]); err != nil {
			log.Fatal(support.RedSp(err))
		}
	case "volume":
		volume.ParseArgs(os.Args[2:])
	case "setupk8s":
		setupk8s.ParseArgs(os.Args[2:])
	default:
		fmt.Println(argEntryFailMsg)
		// fmt.Println("")
		// fmt.Println("Flags for droplet subcommand:")
		// dropCmd.PrintDefaults()
		// fmt.Println("Flags for volume subcommand:")
		// flag.PrintDefaults()
		os.Exit(1)
	}
	// if dropCmd.Parsed() {
	// 	if *dropNamePtr == "" {
	// 		dropCmd.PrintDefaults()
	// 	}
	// }
	// fmt.Printf("*dropNamePtr = %+v\n", *dropNamePtr)
	// fmt.Printf("*dropRegPtr = %+v\n", *dropRegPtr)
}
