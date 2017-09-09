package droplet

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/altnometer/godoapi/lib/support"
	"github.com/briandowns/spinner"
	"github.com/digitalocean/godo"
)

// var client = support.GetDOClient()

// ParseArgsDeleteDrop deletes droplets.
func ParseArgsDeleteDrop(args []string) error {
	// client := support.GetDOClient()
	subCmd := flag.NewFlagSet("delete", flag.ExitOnError)
	// regPtr := subCmd.String("region", "fra1", "-region=fra1")
	var multiTag support.NameList
	// TODO: add functionality fo 'all' tag option
	subCmd.Var(&multiTag, "tag", "-tag=<all|tag1[,tag2...]>")

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
		return deleteAllDroplets()
	}
	return nil
}
func deleteAllDroplets() error {
	volsByDropIds := getAttVolsByDropIds()
	// fmt.Printf("volsByDropIds = %+v\n", volsByDropIds)
	droplets, err := ReturnDroplets()
	if err != nil {
		return err
	}
	for i, d := range droplets {
		for id, v := range volsByDropIds {
			if id == d.ID {
				if err := deleteDropWithAttVols(d, v); err != nil {
					return err
				}
				droplets = append(droplets[:i], droplets[i+1:]...)
			}
		}
	}
	for _, dData := range droplets {
		dID := dData.ID
		fmt.Println("Delete droplet?[y/N]")
		fmt.Printf("  Name  %v\n", support.YellowSp(dData.Name))
		fmt.Printf("  ID    %+v\n", dID)
		fmt.Printf("  tag   %+v\n", dData.Tags)
		reader := bufio.NewReader(os.Stdin)
		char, _, err := reader.ReadRune()
		if err != nil {
			panic(err)
		}
		if char == 10 && char != 'y' && char != 'Y' {
			continue
		}
		s := spinner.New(spinner.CharSets[9], 150*time.Millisecond)
		s.Start()
		res, err := support.DOClient.Droplets.Delete(support.Ctx, dID)
		s.Stop()
		if err != nil {
			panic(err)
		}
		fmt.Printf("res = %+v\n\n", res)
	}
	return nil
}

func getAttVolsByDropIds() (volsByDropIds map[int]godo.Volume) {
	volsByDropIds = make(map[int]godo.Volume)
	fmt.Println("Fetching attached volumes...")
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
	for _, v := range volumes {
		if len(v.DropletIDs) != 0 {
			for _, id := range v.DropletIDs {
				volsByDropIds[id] = v
			}
		}
	}
	return volsByDropIds
}

func deleteDropWithAttVols(d godo.Droplet, v godo.Volume) error {
	// fmt.Printf("d = %+v\n", d)
	var (
		ip  string
		err error
	)
	if ip, err = d.PublicIPv4(); err != nil {
		return err
	}
	diskByIDName := support.VolByIDPrefix + v.Name
	partName := diskByIDName + "-part1"
	pathDiskDir := "/dev/disk/by-id/"
	pathDiskByID := pathDiskDir + diskByIDName
	// pathPartition := pathDiskDir + partName
	mntPoint := "/mnt/" + v.Name
	// lsof | grep mntPoint
	// would not detect NFS exported
	// showmount -e remote_nfs_server
	cmdUmount := fmt.Sprintf("if [ $(mount | grep %s --count) -ge 1 ]; then ", mntPoint) +
		fmt.Sprintf(
			"for i in $(lsof +f -- %s | awk 'NR>1 { print $2 }'); do ", mntPoint) +
		"echo \"kill process $i\" && kill $i; done && " +
		fmt.Sprintf("umount %[1]s && echo -e '\\e[32m unmounted %[1]s \\e[00m';", mntPoint) +
		fmt.Sprintf("else echo -e '\\e[33m %s is not mounted. \\e[00m'; fi", pathDiskByID)
	sshCmds := []string{cmdUmount}
	support.ExecSSH("root", ip, sshCmds)
	promtMsg := fmt.Sprintf("Confirm %s is not mounted", pathDiskByID)
	confirmed, err := support.UserConfirmDefaultN(promtMsg)
	if err != nil {
		return err
	}
	if !confirmed {
		return fmt.Errorf("User did not confirm unmouning of %s", pathDiskByID)
	}
	action, _, err := support.DOClient.StorageActions.DetachByDropletID(support.Ctx, v.ID, d.ID)
	if err != nil {
		return err
	}
	fmt.Printf("action = %+v\n", action)
	fmt.Println("Sleep until the volume is detached.")
	time.Sleep(3 * time.Second)
	promtMsg = fmt.Sprintf("Delete droplet %s?\n  ID %v\n  tag %v\n",
		support.YellowSp(d.Name), d.ID, d.Tags)
	confirmed, err = support.UserConfirmDefaultN(promtMsg)
	if !confirmed {
		return fmt.Errorf("User did not confirm unmouning of %s", pathDiskByID)
	}
	res, err := support.DOClient.Droplets.Delete(support.Ctx, d.ID)
	if err != nil {
		return err
	}
	fmt.Printf("res = %+v\n", res)
	return nil
}
