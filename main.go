// Deploy k8s on digitalocean.
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/altnometer/godoapi/admin"
	"github.com/altnometer/godoapi/droplet"
	"github.com/altnometer/godoapi/ip"
	"github.com/altnometer/godoapi/k8s/setupk8s"
	"github.com/altnometer/godoapi/lib/support"
	"github.com/altnometer/godoapi/volume"
)

var argEntryFailMsg = fmt.Sprintf("Provide <%s|%s|%s|%s|%s> subcommand, please.",
	support.YellowSp("admin"),
	support.YellowSp("droplet"),
	support.YellowSp("ip"),
	support.YellowSp("setupk8s"),
	support.YellowSp("volume"),
)

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
	case "ip":
		if err := ip.ParseArgs(os.Args[2:]); err != nil {
			log.Fatal(support.RedSp(err))
		}
	case "volume":
		volume.ParseArgs(os.Args[2:])
	case "setupk8s":
		if err := setupk8s.ParseArgs(os.Args[2:]); err != nil {
			log.Fatal(support.RedSp(err))
		}
	case "admin":
		err := admin.ParseArgs(os.Args[2:])
		if err != nil {
			log.Fatal(err)
		}
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
