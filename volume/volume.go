package volume

import (
	"flag"
	"fmt"
	"os"

	"github.com/altnometer/godoapi/lib/support"
	"github.com/digitalocean/godo"
)

var argVolumeFailMsg = fmt.Sprintf("Provide <%s|%s|%s|%s|%s|%s> subcommand, please.",
	support.YellowSp("list"), support.YellowSp("create"),
	support.YellowSp("attach"), support.YellowSp("detach"),
	support.YellowSp("mount"),
	support.YellowSp("delete"))

// ParseArgs handles os.Args and calls relevant functions in the package.
func ParseArgs(args []string) error {
	if len(args) < 1 {
		fmt.Println(argVolumeFailMsg)
		os.Exit(1)
	}
	switch args[0] {
	case "list":
		ParseArgsListVol(args[1:])
	case "create":
		ParseArgsCreateVol(args[1:])
	case "attach":
		return ParseArgsAttachVol(args[1:])
	case "detach":
		ParseArgsDetachVol(args[1:])
	case "mount":
		return ParseArgsMountVol(args[1:])
	case "delete":
		ParseArgsDeleteVol(args[1:])
	default:
		fmt.Print("Incorrect arg: ")
		support.RedBold.Println(args[0])
		fmt.Println(argVolumeFailMsg)
		os.Exit(1)
	}
	return nil
}

// ParseArgsListVol handles 'volume list' subcommand.
func ParseArgsListVol(args []string) {
	ListAll()
	return
}

// ParseArgsCreateVol handles 'volume create' subcommand.
func ParseArgsCreateVol(args []string) error {
	volCmd := flag.NewFlagSet("create", flag.ExitOnError)
	namePtr := volCmd.String("name", "", "-name=<volname>")
	regPtr := volCmd.String("region", "fra1", "-region=fra1")
	descPtr := volCmd.String("description", "", "-description=<\"your volume description\">")
	sizePtr := volCmd.Int("size", 0, "-size=<5|10|...>")
	volCmd.Parse(args)
	if len(args) < 1 {
		fmt.Println("Provide the args, please.")
		volCmd.PrintDefaults()
		os.Exit(1)
	}
	if volCmd.Parsed() {
		if *namePtr == "" {
			volCmd.PrintDefaults()
			os.Exit(1)
		}
		if *descPtr == "" {
			volCmd.PrintDefaults()
			os.Exit(1)
		}
		if *sizePtr == 0 {
			volCmd.PrintDefaults()
			os.Exit(1)
		}
	}
	support.ValidateRegions(regPtr)
	// fmt.Printf("*namePtr = %+v\n", *namePtr)
	// fmt.Printf("*regPtr = %+v\n", *regPtr)
	// fmt.Printf("*descPtr = %+v\n", *descPtr)

	createVolData := &godo.VolumeCreateRequest{
		Region:        *regPtr,
		Name:          *namePtr,
		Description:   *descPtr,
		SizeGigaBytes: int64(*sizePtr),
	}
	_, err := Create(createVolData)
	if err != nil {
		return err
	}
	return nil
}

// ParseArgsAttachVol handles 'volume create' subcommand.
func ParseArgsAttachVol(args []string) error {
	volCmd := flag.NewFlagSet("attach", flag.ExitOnError)
	volNamePtr := volCmd.String("vol-name", "", "--vol-name=<volume-name>")
	dropNamePtr := volCmd.String("drop-name", "", "--drop-name=<droplet-name|volume-name>")
	regPtr := volCmd.String("region", "fra1", "-region=fra1")
	descPtr := volCmd.String("description", "", "-description=<\"your volume description\">")
	sizePtr := volCmd.Int("size", 10, "-size=<5|10|...>")
	volCmd.Parse(args)
	if len(args) < 1 {
		fmt.Println("Provide the args, please.")
		volCmd.PrintDefaults()
		return support.ErrBadArgs
	}
	if volCmd.Parsed() {
		if *volNamePtr == "" {
			volCmd.PrintDefaults()
			return support.ErrBadArgs
		}
		if *dropNamePtr == "" {
			*dropNamePtr = *volNamePtr
		}
	}
	if err := support.ValidateRegions(regPtr); err != nil {
		return err
	}
	// fmt.Printf("*namePtr = %+v\n", *namePtr)
	// fmt.Printf("*regPtr = %+v\n", *regPtr)
	// fmt.Printf("*descPtr = %+v\n", *descPtr)

	createVolData := &godo.VolumeCreateRequest{
		Region:        *regPtr,
		Name:          *volNamePtr,
		Description:   *descPtr,
		SizeGigaBytes: int64(*sizePtr),
	}
	_, _, err := Attach(createVolData, *dropNamePtr)
	return err
}

// ParseArgsMountVol mounts specified by args volume to droplet with given name
func ParseArgsMountVol(args []string) error {
	volCmd := flag.NewFlagSet("mount", flag.ExitOnError)
	volNamePtr := volCmd.String("vol-name", "", "--vol-name=<volume-name>")
	dropNamePtr := volCmd.String("drop-name", "", "--drop-name=<droplet-name|volume-name>")
	regPtr := volCmd.String("region", "fra1", "-region=fra1")
	sizePtr := volCmd.Int("size", 10, "-size=<5|10|...>")
	volCmd.Parse(args)
	if len(args) < 1 {
		fmt.Println("Provide the args, please.")
		volCmd.PrintDefaults()
		return support.ErrBadArgs
	}
	if volCmd.Parsed() {
		if *volNamePtr == "" {
			volCmd.PrintDefaults()
			return support.ErrBadArgs
		}
		if *dropNamePtr == "" {
			*dropNamePtr = *volNamePtr
		}
	}
	if err := support.ValidateRegions(regPtr); err != nil {
		return err
	}
	createVolData := &godo.VolumeCreateRequest{
		Region:        *regPtr,
		Name:          *volNamePtr,
		SizeGigaBytes: int64(*sizePtr),
	}
	return Mount(createVolData, *dropNamePtr)
}

// ParseArgsDetachVol handles 'volume detach' subcommand.
func ParseArgsDetachVol(args []string) {
	// volCmd := flag.NewFlagSet("detach", flag.ExitOnError)
	fmt.Println("'detach' subcmd for volumes is not implemented yet")
	return
}

// ParseArgsDeleteVol handles 'volume delete' subcommand.
func ParseArgsDeleteVol(args []string) {
	subCmd := flag.NewFlagSet("delete", flag.ExitOnError)
	// regPtr := subCmd.String("region", "fra1", "-region=fra1")
	var multiTag support.NameList
	// TODO: add functionality fo 'all' tag option
	// subCmd.Var(&multiTag, "tag", "-tag=<all|tag1[,tag2...]>")
	subCmd.Var(&multiTag, "tag", "-tag=<all>")

	subCmd.Parse(args)
	if len(args) < 1 {
		fmt.Println("Provide the args, please.")
		subCmd.PrintDefaults()
		os.Exit(1)
	}
	// if subCmd.Parsed() {
	// 	if (&multiName).String() == "[]" {
	// 		// subCmd.PrintDefaults()
	// 		os.Exit(1)
	// 	}
	// }
	if multiTag[0] == "all" {
		DeleteAll()
	}
	return
}
