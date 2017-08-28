package droplet

import (
	"flag"
	"fmt"

	"github.com/altnometer/godoapi/lib/support"
	"github.com/digitalocean/godo"
)

// ParseArgsListDrop handles os.Args and calls relevant functions in the package.
func ParseArgsListDrop(args []string) {
	subCmd := flag.NewFlagSet("list", flag.ExitOnError)
	listAllPtr := subCmd.Bool("all", true, "-all=true")
	regPtr := subCmd.String("region", "fra1", "-region=fra1")
	subCmd.Parse(args)
	// if len(args) < 1 {
	// 	support.Red.Println("Provide the args, please.")
	// 	subCmd.PrintDefaults()
	// 	os.Exit(1)
	// }
	// if subCmd.Parsed() {
	// 	if (&multiName).String() == "[]" {
	// 		// subCmd.PrintDefaults()
	// 		os.Exit(1)
	// 	}
	// }
	if *listAllPtr {
		listAllDroplets()
		fmt.Println("Printing droplet")
	}
	support.ValidateRegions(regPtr)
	fmt.Printf("*regPtr = %+v\n", *regPtr)
	fmt.Printf("*listAllPtr = %+v\n", *listAllPtr)

}
func listAllDroplets() {
	opt := &godo.ListOptions{
		Page:    1,
		PerPage: 200,
	}
	droplets, _, err := support.DOClient.Droplets.List(support.Ctx, opt)
	if err != nil {
		panic("support.DOclient.Droplets.List() failed.")
	}
	for _, d := range droplets {
		fmt.Printf("d.Name = %+v\n", d.Name)
		fmt.Printf("d.ID = %+v\n", d.ID)
		fmt.Printf("d.Size = %+v\n", d.Size)
		// fmt.Printf("d = %+v\n", d)
		support.RedLn("***************************\n")
	}
	// fmt.Printf("droplets = %+v\n", droplets)
}
