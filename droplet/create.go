package droplet

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/altnometer/godoapi/lib/support"
	"github.com/briandowns/spinner"
	"github.com/digitalocean/godo"
)

// GetDefaultDropCreateData return &godo.DropletMultiCreateRequest with
// some default values
func GetDefaultDropCreateData() *godo.DropletMultiCreateRequest {
	return &godo.DropletMultiCreateRequest{
		Names:             []string{"sub-01.example.com"},
		SSHKeys:           []godo.DropletCreateSSHKey{support.SSHKeys},
		PrivateNetworking: true,
		Region:            "fra1",
		Size:              "512mb",
		Image: godo.DropletCreateImage{
			Slug: "ubuntu-16-04-x64",
		},
		IPv6: true,
		Tags: []string{"web"},
	}
}

// CreatedDrpSpecs holds specs of newly created droplet.
type CreatedDrpSpecs struct {
	Name      string
	ID        int
	PublicIP  string
	PrivateIP string
}

// ParseArgsCreateDrop handles os.Args and calls relevant functions in the package.
func ParseArgsCreateDrop(args []string) error {
	// client := support.GetDOClient()
	subCmd := flag.NewFlagSet("create", flag.ExitOnError)
	var multiName support.NameList
	subCmd.Var(&multiName, "name", "-name=<name1[,name2...]>")
	regPtr := subCmd.String("region", "fra1", "-region=fra1")
	envPtr := subCmd.String("env", "dev", "-env=<prod|test|stage|dev>")
	sizePtr := subCmd.String("size", "512mb", "-size=<512mb|1gb|2gb...>")
	var multiTag support.NameList
	subCmd.Var(&multiTag, "tag", "-tag=<tag1[,tag2...]>")

	subCmd.Parse(args)
	if len(args) < 1 {
		fmt.Println("Provide the args, please.")
		subCmd.PrintDefaults()
		os.Exit(1)
	}
	if subCmd.Parsed() {
		if len(multiName) == 0 {
			support.RedLn("No name(s) provided.")
			subCmd.PrintDefaults()
			os.Exit(1)
		}
	}
	sizeIsValid, err := support.ValidateMemSize(*sizePtr)
	if !sizeIsValid {
		return err
	}
	// if subCmd.Parsed() {
	// 	if (&multiName).String() == "[]" {
	// 		// subCmd.PrintDefaults()
	// 		os.Exit(1)
	// 	}
	// }
	if err := support.ValidateRegions(regPtr); err != nil {
		return err

	}
	createDropData := GetDefaultDropCreateData()
	createDropData.Size = *sizePtr
	createDropData.Region = *regPtr
	createDropData.Names = multiName
	createDropData.Tags = append(multiTag, *envPtr)
	// fmt.Printf("createDropData = %+v\n", createDropData)
	// fmt.Printf("(&multiName).String() = %+v\n", (&multiName).String())
	// fmt.Printf("multiName[0] = %+v\n", multiName[0])
	// fmt.Printf("(&multiTag).String() = %+v\n", (&multiTag).String())
	// // fmt.Printf("*namePtr = %+v\n", *namePtr)
	// fmt.Printf("*regPtr = %+v\n", *regPtr)
	// fmt.Printf("*sizePtr = %+v\n", *sizePtr)
	_, err = CreateDroplet(createDropData)
	if err != nil {
		return err
	}
	return nil
}

// CreateDroplet creates a droplet with provided specs.
func CreateDroplet(
	reqDataPtr *godo.DropletMultiCreateRequest) ([]godo.Droplet, error) {
	confirmed, err := support.UserConfirmDefaultY(
		fmt.Sprintf("Creating %v droplet(s)?",
			support.YellowSp(reqDataPtr.Names)))
	if err != nil {
		return nil, err
	}
	if !confirmed {
		return nil, nil
	}
	s := spinner.New(spinner.CharSets[9], 150*time.Millisecond)
	s.Start()
	droplets, _, err := support.DOClient.Droplets.CreateMultiple(support.Ctx, reqDataPtr)
	s.Stop()
	if err != nil {
		return nil, err
	}

	// fmt.Printf("droplets = %+v\n", droplets)
	// dspecs := make([]CreatedDrpSpecs, support.MaxDroplets)
	for _, d := range droplets {
		// ds := CreatedDrpSpecs{Name: d.Name, ID: d.ID}
		// Networks are not available at this stage.
		// for _, n := range d.Networks.V4 {
		// 	if n.Type == "public" {
		// 		ds.PublicIP = n.IPAddress
		// 	}
		// 	if n.Type == "private" {
		// 		ds.PrivateIP = n.IPAddress
		// 	}
		// }
		// dspecs[i] = ds
		fmt.Println("Created droplet with:")
		fmt.Printf("  d.Name = %+v\n", d.Name)
		fmt.Printf("  d.ID = %+v\n", d.ID)
		// fmt.Printf("d.Size = %+v\n", d.Size)
		// fmt.Printf("d = %+v\n", d)
		fmt.Println("***************************")
	}
	return droplets, nil

}
