package setupk8s

import (
	"flag"

	"github.com/altnometer/godoapi/lib/support"
)

// ParseArgs handles os.Args and calls relevant functions in the package.
func ParseArgs(args []string) {
	subCmd := flag.NewFlagSet("setupk8s", flag.ExitOnError)
	envPtr := subCmd.String("env", "dev", "-env=<prod|test|stage|dev>")
	regPtr := subCmd.String("region", "fra1", "-region=fra1")
	subCmd.Parse(args)
	if subCmd.Parsed() {
		// support.ValidateRegions(regPtr)
		switch *envPtr {
		case "dev":
			ip, token := SetUpMaster(*envPtr, *regPtr)
			SetUpNode(*envPtr, *regPtr, ip, token)
		default:
			support.RedLn("Provide valid -env value, please.")
		}
	}
	// if len(args) < 1 {
	// 	support.Red.Println("Provide the args, please.")
	// 	subCmd.PrintDefaults()
	// 	os.Exit(1)
	// }
	// fmt.Printf("*regPtr = %+v\n", *regPtr)
}
