package resources

import (
	"context"
	"time"

	"github.com/cgroschupp/hetzner-nuke/pkg/hetzner"
	"github.com/cgroschupp/hetzner-nuke/pkg/nuke"
	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

const FirewallResource = "Firewall"

func init() {
	registry.Register(&registry.Registration{
		Name:      FirewallResource,
		Scope:     nuke.Account,
		Resource:  &Firewall{},
		Lister:    &FirewallLister{},
		DependsOn: []string{ServerResource},
	})
}

type FirewallLister struct{}

func (l *FirewallLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	resources := make([]resource.Resource, 0)
	firewalls, _, err := opts.Client.Firewall.List(ctx, hcloud.FirewallListOpts{})
	if err != nil {
		return resources, err
	}
	for _, firewall := range firewalls {
		resources = append(resources, &Firewall{obj: firewall, Name: &firewall.Name, Labels: firewall.Labels, ID: &firewall.ID, client: opts.Client, Created: &firewall.Created})
	}
	return resources, nil
}

type Firewall struct {
	client  *hetzner.Client
	obj     *hcloud.Firewall
	Name    *string
	ID      *int64
	Labels  map[string]string `description:"The labels associated with the firewall"`
	Created *time.Time        `description:"The time the firewall was created"`
}

func (r *Firewall) Remove(ctx context.Context) error {
	_, err := r.client.Firewall.Delete(ctx, r.obj)

	return err
}

func (r *Firewall) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *Firewall) String() string {
	return *r.Name
}
