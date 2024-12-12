package hetzner

import (
	"fmt"
	"time"

	libnuke "github.com/ekristen/libnuke/pkg/nuke"
	"github.com/ekristen/libnuke/pkg/utils"
)

// Prompt is a struct that contains the parameters and tenant details use to craft a unique prompt
// for the user to confirm the nuke operation.
type Prompt struct {
	Parameters *libnuke.Parameters
}

func (p *Prompt) Prompt() error {
	forceSleep := time.Duration(p.Parameters.ForceSleep) * time.Second

	fmt.Println("Do you really want to nuke the Hetzner account.")
	if p.Parameters.Force {
		fmt.Printf("Waiting %v before continuing.\n", forceSleep)
		time.Sleep(forceSleep)
	} else {
		fmt.Printf("Do you want to continue? Enter %q to continue.\n", "yes")
		if err := utils.Prompt("yes"); err != nil {
			return err
		}
	}

	return nil
}
