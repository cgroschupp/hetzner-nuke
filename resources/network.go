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

const NetworkResource = "Network"

func init() {
	registry.Register(&registry.Registration{
		Name:     NetworkResource,
		Scope:    nuke.Account,
		Resource: &Network{},
		Lister:   &NetworkLister{},
	})
}

type NetworkLister struct{}

func (l *NetworkLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	resources := make([]resource.Resource, 0)
	networks, _, err := opts.Client.Network.List(ctx, hcloud.NetworkListOpts{})
	if err != nil {
		return resources, err
	}
	for _, network := range networks {
		resources = append(resources, &Network{obj: network, Name: &network.Name, Labels: network.Labels, ID: &network.ID, client: opts.Client, Created: &network.Created})
	}
	return resources, nil
}

type Network struct {
	client  *hetzner.Client
	obj     *hcloud.Network
	Name    *string
	ID      *int64
	Labels  map[string]string `description:"The labels associated with the network"`
	Created *time.Time        `description:"The time the network was created"`
}

func (r *Network) Remove(ctx context.Context) error {
	_, err := r.client.Network.Delete(ctx, r.obj)

	return err
}

func (r *Network) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *Network) String() string {
	return *r.Name
}
