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
	Project    *Project
}

func (p *Prompt) Prompt() error {
	forceSleep := time.Duration(p.Parameters.ForceSleep) * time.Second

	if p.Parameters.Force {
		fmt.Printf("no-prompt flag set, continuing without prompting user")
		fmt.Printf("waiting %v before continuing", forceSleep)
		time.Sleep(forceSleep)
	} else {
		fmt.Printf("Do you really want to nuke the Hetzner project with the ID %d and the name %q.\n", p.Project.ID, p.Project.Name)
		fmt.Printf("Do you want to continue? Enter project name %q to continue.\n", p.Project.Name)
		if err := utils.Prompt(p.Project.Name); err != nil {
			return err
		}
	}

	return nil
}
