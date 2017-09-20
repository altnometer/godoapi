// Package ip handle digitalocean floating ips.
package ip

import (
	"fmt"
	"strconv"

	"github.com/altnometer/godoapi/droplet"
	"github.com/altnometer/godoapi/lib/support"
	"github.com/digitalocean/godo"
)

var argIPFailMsg = fmt.Sprintf("Provide <%s|%s|%s> subcommand, please.",
	support.YellowSp("list"), support.YellowSp("create"), support.YellowSp("delete"))

// ParseArgs handles os.Args and calls relevant functions in the package.
func ParseArgs(args []string) error {
	if len(args) < 1 {
		fmt.Println(argIPFailMsg)
		return support.ErrBadArgs
	}
	switch args[0] {
	case "assign":
		if err := assignIP(); err != nil {
			return err
		}
	case "list":
		if err := listIPs(); err != nil {
			return err
		}
	case "create":
		fmt.Println("'ip create' is not implemented yet.")
		return nil
	case "delete":
		fmt.Println("'ip delete' is not implemented yet.")
		return nil
	default:
		fmt.Print("Incorrect arg: ")
		support.RedBold.Println(args[0])
		fmt.Println(argIPFailMsg)
		return support.ErrBadArgs
	}
	return nil
}
func listIPs() error {
	fIPs, err := getFloatIps()
	if err != nil {
		return err
	}
	for _, ip := range fIPs {
		// fmt.Printf("ip = %+v\n", ip)
		d := ip.Droplet
		if d == nil {
			support.YellowLn("Available")
		} else {
			fmt.Printf("ip.Droplet.ID   %+v\n", ip.Droplet.ID)
		}
		fmt.Printf("ip.IP           %+v\n", ip.IP)
		fmt.Printf("ip.Region.Slug  %+v\n", ip.Region.Slug)
		fmt.Println("**************************")
	}
	return nil

}

func getFloatIps() ([]godo.FloatingIP, error) {
	opt := &godo.ListOptions{
		Page:    1,
		PerPage: 200,
	}
	fIPs, _, err := support.DOClient.FloatingIPs.List(support.Ctx, opt)
	if err != nil {
		return nil, err
	}
	return fIPs, nil
}

func assignIP() error {
	fmt.Println("called assignIP()")
	fIPs, err := getFloatIps()
	if err != nil {
		return err
	}
	for i, ip := range fIPs {
		d := ip.Droplet
		if d == nil {
			fmt.Printf("%d - %s, region - %s\n", i, ip.IP, ip.Region.Slug)
		}
	}
	ipNum, err := strconv.ParseInt(
		support.GetUserInput("Type in a number of the desired ip: "), 10, 0)
	if err != nil {
		return err
	}
	if ipNum < 0 || ipNum >= int64(len(fIPs)) {
		return fmt.Errorf("choice must be withing 0 to %d range", len(fIPs)-1)
	}
	ip := fIPs[ipNum]
	fmt.Printf("You selected ip = %+v\n", support.YellowSp(ip.IP))
	drops, err := droplet.ReturnDroplets()
	if err != nil {
		return err
	}
	for i, d := range drops {
		fmt.Printf("%d - %s, region - %s, tags - %v\n", i, d.Name, d.Region.Slug, d.Tags)
	}
	dNum, err := strconv.ParseInt(
		support.GetUserInput("Type in a number of the desired host: "), 10, 0)
	if err != nil {
		return err
	}
	if dNum < 0 || dNum >= int64(len(drops)) {
		return fmt.Errorf("choice must be withing 0 to %d range", len(drops)-1)
	}
	d := drops[dNum]
	fmt.Printf("You selected host %+v\n", support.YellowSp(d.Name))
	support.YellowPf("Assigning %s to %s\n", ip.IP, d.Name)
	action, _, err := support.DOClient.FloatingIPActions.Assign(
		support.Ctx, ip.IP, d.ID)
	if err != nil {
		return err
	}
	fmt.Printf("action = %+v\n", action)
	return nil
}
