package volume

import (
	"bufio"
	"fmt"
	"log"
	"os"
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
			fmt.Printf("v = %+v\n", v)
			fmt.Println("************************************")
		}
	} else {
		support.GreenLn("No volumes exist.")
	}
	return
}

// GetAllVols gets all volumes.
func GetAllVols() *[]godo.Volume {
	fmt.Println("Fetching existing volumes...")
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
func Attach(
	vd *godo.VolumeCreateRequest,
	dropName string) (vol *godo.Volume, drop *godo.Droplet, err error) {
	vols := *GetAllVols()
	for _, v := range vols {
		if v.Name == vd.Name {
			vol = &v
			break
		}
	}
	if vol == nil {
		vol, err = Create(vd)
		if err != nil {
			return nil, nil, err
		}
	}
	var droplets []godo.Droplet
	if droplets, err = droplet.ReturnDroplets(); err != nil {
		return nil, nil, err
	}
	for _, d := range droplets {
		if d.Name == dropName {
			drop = &d
		}
	}
	if drop == nil {
		drCreateData := droplet.GetDefaultDropCreateData()
		drCreateData.Names = []string{dropName}
		drCreateData.Region = vd.Region
		drCreateData.Tags = []string{"volume", vd.Name}
		droplets := droplet.CreateDroplet(drCreateData)
		if len(droplets) == 0 {
			return nil, nil, fmt.Errorf("failed to create droplet to attach volume to")
		}
		// give time for droplet to boot up.
		fmt.Println("Waiting for created droplet to boot up...")
		time.Sleep(10 * time.Second)
		drop = &droplets[0]
	}
	// func ReturnDropletsData() *[]map[string]string {
	if len(vol.DropletIDs) != 0 {
		for _, di := range vol.DropletIDs {
			if di == drop.ID {
				return vol, drop, nil
			}
		}
		err := support.ErrVolAttached{
			Msg: fmt.Sprintf(
				"volume %s is attached to different droplet(s) with id(s) %v",
				vol.Name, vol.DropletIDs)}
		return nil, nil, err
	}

	actionPtr, resPtr, err := support.DOClient.StorageActions.Attach(
		support.Ctx, vol.ID, drop.ID)
	if err != nil {
		return nil, nil, err
	}
	if actionPtr != nil {
		fmt.Printf("*actionPtr = %+v\n", *actionPtr)
	}
	fmt.Printf("resPtr = %+v\n", resPtr)
	return vol, drop, nil
}

// Mount function mounts volume with vol ID to droplet with drop ID.
func Mount(vd *godo.VolumeCreateRequest, dropName string) error {
	vol, drop, err := Attach(vd, dropName)
	if err != nil {
		return err
	}

	fmt.Printf("vol = %+v\n", vol)
	fmt.Printf("drop.Networks.V4 = %+v\n", drop.Networks.V4)
	ip, err := drop.PublicIPv4()
	fmt.Printf("ip = %+v\n", ip)

	volSysName := support.VolByIDPrefix + vol.Name
	cmd := "if /sbin/sfdisk -d /dev/disk/by-id/" + volSysName + " 2>/dev/null ; then echo \"partitionen\"; else echo \"not partitioned\";fi"
	fmt.Printf("cmd = %+v\n", cmd)
	sshCmds := []string{cmd}
	sshKeyPath := support.GetSSHKeyPath()
	sshOutput := support.FetchSSHOutput("root", ip, sshKeyPath, sshCmds)
	fmt.Printf("sshOutput = %+v\n", sshOutput)
	// sudo parted /dev/disk/by-id/scsi-0DO_Volume_volume-nyc1-01 mklabel gpt

	return err
	// partitions: -part1, part2 etc
	// Partition the volume
	// Format the partitions
	// Create mount points
	// Mount the filesystems
	// Adjust the /etc/fstab
}
