package droplet

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/altnometer/godoapi/lib/support"
	"github.com/briandowns/spinner"
)

// var client = support.GetDOClient()

// ParseArgsDeleteDrop deletes droplets.
func ParseArgsDeleteDrop(args []string) {
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
		deleteAllDroplets()
	}

}
func deleteAllDroplets() {
	droplets := *ReturnDropletData()
	for _, dData := range droplets {
		if dData != nil {
			dID, err := strconv.Atoi(dData["ID"])
			if err != nil {
				panic(err)
			}

			// fmt.Printf("dData[ID] %+v\n", dData["ID"])
			reader := bufio.NewReader(os.Stdin)
			fmt.Println("Delete droplet?[y/N]")
			fmt.Printf("  Name  %v\n", support.YellowSp(dData["Name"]))
			fmt.Printf("  ID    %+v\n", dID)
			fmt.Printf("  tag   %+v\n", dData["Tags"])
			char, _, err := reader.ReadRune()
			if err != nil {
				panic(err)
			}
			if char == 10 && char != 'y' && char != 'Y' {
				os.Exit(0)
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
	}
}
