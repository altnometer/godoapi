package droplet

import (
	"flag"
	"fmt"
	"os"

	"github.com/altnometer/godoapi/lib/support"
	"github.com/digitalocean/godo"
)

// var client = support.GetDOClient()

var createRequestData = &godo.DropletMultiCreateRequest{
	Names:  []string{"sub-01.example.com"},
	Region: "nyc3",
	Size:   "512mb",
	Image: godo.DropletCreateImage{
		Slug: "ubuntu-16-04-x64",
	},
	IPv6: true,
	Tags: []string{"web"},
}

// ParseArgsCreateDrop handles os.Args and calls relevant functions in the package.
func ParseArgsCreateDrop(args []string) {
	// client := support.GetDOClient()
	subCmd := flag.NewFlagSet("create", flag.ExitOnError)
	var multiName support.NameList
	subCmd.Var(&multiName, "name", "-name=<name1[,name2...]>")
	regPtr := subCmd.String("region", "fra1", "-region=fra1")
	sizePtr := subCmd.String("size", "512mb", "-size=<512mb|1gb|2gb...>")
	var multiTag support.NameList
	subCmd.Var(&multiTag, "tag", "-tag=<tag1[,tag2...]>")

	subCmd.Parse(args)
	if len(args) < 1 {
		fmt.Println("Provide the args, please.")
		subCmd.PrintDefaults()
		os.Exit(1)
	}
	var sizeIsValid bool
	for _, size := range support.DropletSizes {
		if *sizePtr == size {
			sizeIsValid = true
			break
		}
	}
	if sizeIsValid == false {
		support.Red.Printf("Valid -size values %+v\n", support.DropletSizes)
	}
	fmt.Printf("sizeIsValid = %+v\n", sizeIsValid)
	// if subCmd.Parsed() {
	// 	if (&multiName).String() == "[]" {
	// 		// subCmd.PrintDefaults()
	// 		os.Exit(1)
	// 	}
	// }
	support.ValidateRegions(regPtr)
	createRequestData.Size = *sizePtr
	createRequestData.Region = *regPtr
	createRequestData.Names = multiName
	createRequestData.Tags = multiTag
	// fmt.Printf("createRequestData = %+v\n", createRequestData)
	// fmt.Printf("(&multiName).String() = %+v\n", (&multiName).String())
	fmt.Printf("multiName[0] = %+v\n", multiName[0])
	fmt.Printf("(&multiTag).String() = %+v\n", (&multiTag).String())
	// fmt.Printf("*namePtr = %+v\n", *namePtr)
	fmt.Printf("*regPtr = %+v\n", *regPtr)
	fmt.Printf("*sizePtr = %+v\n", *sizePtr)
	createDroplet()

}
func createDroplet() {
	droplet, _, err := support.DOClient.Droplets.CreateMultiple(support.Ctx, createRequestData)
	if err != nil {
		panic(err)
	}

	fmt.Printf("droplet = %+v\n", droplet)

}
