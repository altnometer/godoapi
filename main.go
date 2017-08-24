package main

import (
	"fmt"
	"os"

	"github.com/altnometer/godoapi/droplet"
	"github.com/altnometer/godoapi/volume"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("droplet or volume subcommand is required")
		os.Exit(1)
	}
	switch os.Args[1] {
	case "droplet":
		droplet.ParseArgs(os.Args[2:])
	case "volume":
		volume.ParseArgs(os.Args[2:])
	default:
		fmt.Println("droplet or volume subcommand is required")
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
