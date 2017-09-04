package volume

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/altnometer/godoapi/lib/support"
	"github.com/briandowns/spinner"
	"github.com/digitalocean/godo"
)

// Vol struct is a local version of godo.Volume
type Vol struct {
	godo.Volume
}

func (v Vol) String() string {
	// return fmt.Sprintln(v.Name, v.ID)
	return fmt.Sprintf("Name: %s\nID: %s\nReg: %s",
		v.Name, v.ID, v.Region.Slug)
}

// Delete godo.Volume method deletes the volume
func (v Vol) Delete() (*godo.Response, error) {
	if len(v.DropletIDs) != 0 {
		return nil, fmt.Errorf("Cannot delete attached volume %s", v.Name)
	}
	return support.DOClient.Storage.DeleteVolume(support.Ctx, v.ID)
}

// Create creates a volume with provided specs.
func Create(volCreateData *godo.VolumeCreateRequest) {
	volume, _, err := support.DOClient.Storage.CreateVolume(support.Ctx, volCreateData)
	if err != nil {
		log.Fatal(err)
	}
	v := Vol{*volume}
	// fmt.Printf("volume = %+v\n", volume)
	fmt.Println(v)
	return
}

// ListAll lists all volumes.
func ListAll() {
	volumes := *GetAllVols()
	if len(volumes) > 0 {
		for _, v := range volumes {
			fmt.Printf("v.Name          = %+v\n", v.Name)
			fmt.Printf("v.ID            = %+v\n", v.ID)
			fmt.Printf("v.Region.Slug   = %+v\n", v.Region.Slug)
			fmt.Printf("v.SizeGigaBytes = %+v\n", v.SizeGigaBytes)
			fmt.Printf("v.Description   = %+v\n", v.Description)
			fmt.Printf("v.DropletIDs    = %+v\n", v.DropletIDs)
			fmt.Println("************************************")
		}
	} else {
		support.GreenLn("No volumes exist.")
	}
	return
}

// GetAllVols gets all volumes.
func GetAllVols() *[]godo.Volume {
	// ListVolumeParams stores the options you can set for a ListVolumeCall
	// type ListVolumeParams struct {
	// 	Region      string       `json:"region"`
	// 	Name        string       `json:"name"`
	// 	ListOptions *ListOptions `json:"list_options,omitempty"`
	// }
	opt := &godo.ListOptions{
		Page:    1,
		PerPage: 200,
	}
	listOpt := &godo.ListVolumeParams{
		ListOptions: opt,
	}
	volumes, _, err := support.DOClient.Storage.ListVolumes(support.Ctx, listOpt)
	if err != nil {
		log.Fatal(err)
	}
	return &volumes
}

// DeleteAll deletes all volumes if confirmed by user.
func DeleteAll() {
	volumes := *GetAllVols()
	if len(volumes) != 0 {
		for _, vol := range volumes {
			v := Vol{vol}
			reader := bufio.NewReader(os.Stdin)
			fmt.Println("Delete droplet?[y/N]")
			fmt.Printf("  Name  %v\n", support.YellowSp(v.Name))
			fmt.Printf("  ID    %+v\n", v.ID)
			fmt.Printf("  DropletIDs   %+v\n", v.DropletIDs)
			char, _, err := reader.ReadRune()
			if err != nil {
				panic(err)
			}
			if char == 10 && char != 'y' && char != 'Y' {
				continue
			}
			s := spinner.New(spinner.CharSets[9], 150*time.Millisecond)
			s.Start()
			res, err := v.Delete()
			s.Stop()
			if err != nil {
				panic(err)
			}
			fmt.Printf("res = %+v\n\n", res)
		}
	} else {
		support.GreenLn("No volumes exist.")
	}

}

// Attach function attaches volume with volID to droplet with dropID.
func Attach(volID string, dropID int) {

}
