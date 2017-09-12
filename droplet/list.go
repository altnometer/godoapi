package droplet

import (
	"errors"
	"flag"
	"fmt"
	"strconv"
	"strings"
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
		// ReturnDropletsData()
	}
	support.ValidateRegions(regPtr)
	// fmt.Printf("*regPtr = %+v\n", *regPtr)
	// fmt.Printf("*listAllPtr = %+v\n", *listAllPtr)

}
func listAllDroplets() {
	fmt.Println("Sending request to DO api: ")
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
	support.PrintDropData(droplets)
}

// ReturnDropletsData returns a list of data for each listed droplet.
func ReturnDropletsData() *[]map[string]string {
	fmt.Println("Collecting listed droplets data: ")
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
	dropletsData := make([]map[string]string, support.MaxDroplets)
	if len(droplets) > 0 {
		for i, d := range droplets {
			dData := make(map[string]string)
			dData["ID"] = strconv.Itoa(d.ID)
			dData["Name"] = d.Name
			dData["Tags"] = strings.Join(d.Tags[:], ",")
			dData["Region"] = d.Region.Slug
			dropletsData[i] = dData
			// fmt.Printf("d.Name                %+v\n", d.Name)
			// fmt.Printf("d.ID                  %+v\n", d.ID)
			// fmt.Printf("d.Region.Slug         %+v\n", d.Region.Slug)
			// fmt.Printf("d.Tags                %+v\n", d.Tags)
			// // fmt.Printf("d.Size = %+v\n", d.Size)
			// fmt.Printf("d.Size.Slug           %+v\n", d.Size.Slug)
			// fmt.Printf("d.Size.Memory         %+v\n", d.Size.Memory)
			// fmt.Printf("d.Size.Vcpus          %+v\n", d.Size.Vcpus)
			// fmt.Printf("d.Size.Disk           %+v\n", d.Size.Disk)
			// fmt.Printf("d.Size.PriceMonthly   %+v\n", d.Size.PriceMonthly)
			// fmt.Printf("d = %+v\n", d)
			// fmt.Println("***************************")
		}
	} else {
		support.GreenLn("No droplets exist.")
	}
	return &dropletsData
	// for _, dData := range dropletsData {
	// 	if dData != nil {
	// 		fmt.Printf("dData = %+v\n", dData)
	// 	}
	// }
	// fmt.Printf("dropletsData = %+v\n", dropletsData)
	// fmt.Printf("droplets = %+v\n", droplets)
}

// ReturnDroplets returns a list of data for each listed droplet.
func ReturnDroplets() (volumes []godo.Droplet, err error) {
	fmt.Println("Collecting listed droplets... ")
	s := spinner.New(spinner.CharSets[9], 150*time.Millisecond)
	s.Start()
	opt := &godo.ListOptions{
		Page:    1,
		PerPage: 200,
	}
	droplets, _, err := support.DOClient.Droplets.List(support.Ctx, opt)
	s.Stop()
	if err != nil {
		return nil, err
	}
	return droplets, nil
}

// ReturnDropletByID return droplet data identified by the provided id.
func ReturnDropletByID(id int) map[string]string {
	d, _, err := support.DOClient.Droplets.Get(support.Ctx, id)
	if err != nil {
		support.RedLn("No droplet is returned")
		panic(err)
	}
	dData := make(map[string]string)
	dData["ID"] = strconv.Itoa(d.ID)
	dData["Name"] = d.Name
	dData["Tags"] = strings.Join(d.Tags[:], ",")
	for _, n := range d.Networks.V4 {
		if n.Type == "public" {
			dData["PublicIP"] = n.IPAddress
		}
		if n.Type == "private" {
			dData["PrivateIP"] = n.IPAddress
		}
	}
	dData["Region"] = d.Region.Slug
	return dData
}

// ReturnDropletsByTag return droplets identified by the provided tag.
func ReturnDropletsByTag(tag string) (*[]map[string]string, error) {
	fmt.Println("Collecting listed droplets data fetched by tag name: ")
	s := spinner.New(spinner.CharSets[9], 150*time.Millisecond)
	s.Start()
	opt := &godo.ListOptions{
		Page:    1,
		PerPage: 200,
	}
	droplets, _, err := support.DOClient.Droplets.ListByTag(support.Ctx, tag, opt)
	s.Stop()
	fmt.Println("")
	if err != nil {
		return nil, errors.New("support.DOclient.Droplets.ListByTag() failed")
	}
	dropletsData := make([]map[string]string, support.MaxDroplets)
	if len(droplets) > 0 {
		for i, d := range droplets {
			dData := make(map[string]string)
			dData["ID"] = strconv.Itoa(d.ID)
			dData["Name"] = d.Name
			dData["Tags"] = strings.Join(d.Tags[:], ",")
			dData["Region"] = d.Region.Slug
			for _, n := range d.Networks.V4 {
				if n.Type == "public" {
					dData["PublicIP"] = n.IPAddress
				}
				if n.Type == "private" {
					dData["PrivateIP"] = n.IPAddress
				}
			}
			dropletsData[i] = dData
		}
	} else {
		support.GreenLn("No droplets exist.")
	}
	return &dropletsData, nil
}
