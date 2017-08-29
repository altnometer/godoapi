package support

import (
	"fmt"
	"os"
	"strings"

	"github.com/digitalocean/godo"
	"github.com/fatih/color"
)

// YellowSp colors str yellow.
var YellowSp = color.New(color.FgYellow).SprintFunc()

// RedLn colors a line red.
var RedLn = color.New(color.FgRed).PrintlnFunc()

// GreenLn colors a line red.
var GreenLn = color.New(color.FgGreen).PrintlnFunc()

// RedBold colors the output.
var RedBold = color.New(color.FgRed, color.Bold)

// Red colors the output.
var Red = color.New(color.FgRed)

/////////////// Common droplet properties//////////////////////////////////////

// SSHKeys for the droplet.
var SSHKeys = godo.DropletCreateSSHKey{Fingerprint: "2f:2a:c4:eb:ec:38:35:cd:2a:d9:65:cf:59:12:df:44"}

// DropletSizes holds valid values.
var DropletSizes = []string{"512mb", "1gb", "2gb"}

// NameList to hold custom flag value for multiple names.
type NameList []string

func (n *NameList) String() string {
	return fmt.Sprintf("%v", *n)
}

// Set parsed the flag value
func (n *NameList) Set(value string) error {
	*n = strings.Split(value, ",")
	return nil
}

// Value interface for cusmom flag argument.
type Value interface {
	String() string
	Set(string) error
}

// ValidateRegions prints out error msg and exits for invalid regions.
func ValidateRegions(regPtr *string) {
	regions := map[string]bool{"fra1": true}
	if _, validChoice := regions[*regPtr]; !validChoice {
		keys := make([]string, len(regions))
		i := 0
		for k := range regions {
			keys[i] = k
			i++
		}
		fmt.Printf("valid choices for region field are: %+v\n", keys)
		os.Exit(1)
	}
}
