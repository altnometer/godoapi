package volume

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/altnometer/godoapi/droplet"
	"github.com/altnometer/godoapi/lib/support"
	"github.com/briandowns/spinner"
	"github.com/digitalocean/godo"
)

func getDefaultVolCreateData() *godo.VolumeCreateRequest {
	return &godo.VolumeCreateRequest{
		Region: "fra1",
		// Name:          string `json:"name"`
		// Description:   string `json:"description"`
		SizeGigaBytes: 10,
		// SnapshotID:    string `json:"snapshot_id"`
	}
}

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
func Create(volCreateData *godo.VolumeCreateRequest) (*godo.Volume, error) {
	reader := bufio.NewReader(os.Stdin)
	support.YellowPf("Creating %v volume?[Y/n] ", volCreateData.Name)
	char, _, err := reader.ReadRune()
	if err != nil {
		return nil, fmt.Errorf("failed reading userinput: %v", err)
	}
	if char != 10 && char != 'y' && char != 'Y' {
		return nil, support.ErrUserSaysQuit
	}
	s := spinner.New(spinner.CharSets[9], 150*time.Millisecond)
	s.Start()
	volume, _, err := support.DOClient.Storage.CreateVolume(support.Ctx, volCreateData)
	if err != nil {
		log.Fatal(err)
	}
	v := Vol{*volume}
	// fmt.Printf("volume = %+v\n", volume)
	fmt.Println(v)
	return volume, nil
}

// ListAll lists all volumes.
func ListAll() {
	s := spinner.New(spinner.CharSets[9], 150*time.Millisecond)
	s.Start()
	volumes := *GetAllVols()
	s.Stop()
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
func Attach(vd *godo.VolumeCreateRequest, dropName string) error {
	s := spinner.New(spinner.CharSets[9], 150*time.Millisecond)
	// func Attach(volID string, dropID int) {
	vols := *GetAllVols()
	var (
		volID  string
		dropID int
		err    error
	)
	for _, v := range vols {
		if v.Name == vd.Name {
			volID = v.ID
		}
	}
	if volID == "" {
		volume, err := Create(vd)
		if err != nil {
			return err
		}

		volID = volume.ID
	}
	droplets := *droplet.ReturnDropletsData()
	for _, d := range droplets {
		if d["Name"] == dropName {
			dropID, err = strconv.Atoi(d["ID"])
			if err != nil {
				return err
			}
		}
	}
	if dropID == 0 {
		drCreateData := droplet.GetDefaultDropCreateData()
		drCreateData.Names = []string{dropName}
		drCreateData.Region = vd.Region
		drCreateData.Tags = []string{"volume", vd.Name}
		s.Start()
		droplets := droplet.CreateDroplet(drCreateData)
		if len(droplets) == 0 {
			err := fmt.Errorf("failed to create droplet to attach volume to")
			panic(err)
		}
		dropID = droplets[0].ID
		// give time for droplet to boot up.
		time.Sleep(10 * time.Second)
		s.Stop()
	}
	// func ReturnDropletsData() *[]map[string]string {
	actionPtr, resPtr, err := support.DOClient.StorageActions.Attach(support.Ctx, volID, dropID)
	if actionPtr != nil {
		fmt.Printf("*actionPtr = %+v\n", *actionPtr)
	}
	fmt.Printf("resPtr = %+v\n", resPtr)
	if err != nil {
		return err
	}
	return nil
}
