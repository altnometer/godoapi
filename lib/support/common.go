package support

import (
	"fmt"
	"os"
)

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
