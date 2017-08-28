package droplet

import (
	"flag"
	"fmt"
	"time"

	"github.com/altnometer/godoapi/lib/support"
	"github.com/briandowns/spinner"
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
	}
	support.ValidateRegions(regPtr)
	// fmt.Printf("*regPtr = %+v\n", *regPtr)
	// fmt.Printf("*listAllPtr = %+v\n", *listAllPtr)

}
func listAllDroplets() {
	fmt.Print("Sending request to DO api: ")
	s := spinner.New(spinner.CharSets[9], 150*time.Millisecond)
	s.Start()
	opt := &godo.ListOptions{
		Page:    1,
		PerPage: 200,
	}
	droplets, _, err := support.DOClient.Droplets.List(support.Ctx, opt)
	s.Stop()
	fmt.Println("")
	if err != nil {
		panic("support.DOclient.Droplets.List() failed.")
	}
	if len(droplets) > 0 {
		for _, d := range droplets {
			fmt.Printf("d.Name                %+v\n", d.Name)
			fmt.Printf("d.ID                  %+v\n", d.ID)
			// fmt.Printf("d.Size = %+v\n", d.Size)
			fmt.Printf("d.Size.Slug           %+v\n", d.Size.Slug)
			fmt.Printf("d.Size.Memory         %+v\n", d.Size.Memory)
			fmt.Printf("d.Size.Vcpus          %+v\n", d.Size.Vcpus)
			fmt.Printf("d.Size.Disk           %+v\n", d.Size.Disk)
			fmt.Printf("d.Size.PriceMonthly   %+v\n", d.Size.PriceMonthly)
			// fmt.Printf("d = %+v\n", d)
			fmt.Println("***************************")
		}
	} else {
		support.GreenLn("No droplets exist.")
	}
	// fmt.Printf("droplets = %+v\n", droplets)
}
