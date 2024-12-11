package resources

import (
	"context"

	"github.com/cgroschupp/hetzner-nuke/pkg/nuke"
	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

const FirewallResource = "Firewall"

func init() {
	registry.Register(&registry.Registration{
		Name:     FirewallResource,
		Scope:    nuke.Account,
		Resource: &Firewall{},
		Lister:   &FirewallLister{},
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
		resources = append(resources, &Firewall{obj: firewall, Name: &firewall.Name, Labels: firewall.Labels, ID: &firewall.ID, client: opts.Client})
	}
	return resources, nil
}

type Firewall struct {
	client *hcloud.Client
	obj    *hcloud.Firewall
	Name   *string
	ID     *int64
	Labels map[string]string
}

func (r *Firewall) Remove(ctx context.Context) error {
	_, err := r.client.Firewall.Delete(ctx, r.obj)

	return err
}

func (r *Firewall) String() string {
	return *r.Name
}
