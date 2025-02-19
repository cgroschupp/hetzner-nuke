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

const FloatingIPResource = "FloatingIP"

func init() {
	registry.Register(&registry.Registration{
		Name:     FloatingIPResource,
		Scope:    nuke.Account,
		Resource: &FloatingIP{},
		Lister:   &FloatingIPLister{},
	})
}

type FloatingIPLister struct{}

func (l *FloatingIPLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	resources := make([]resource.Resource, 0)
	floatingIPS, _, err := opts.Client.FloatingIP.List(ctx, hcloud.FloatingIPListOpts{})
	if err != nil {
		return resources, err
	}
	for _, floatingIP := range floatingIPS {
		resources = append(resources, &FloatingIP{obj: floatingIP, Name: &floatingIP.Name, Labels: floatingIP.Labels, ID: &floatingIP.ID, client: opts.Client, Created: &floatingIP.Created, Location: &floatingIP.HomeLocation.Name})
	}
	return resources, nil
}

type FloatingIP struct {
	client   *hetzner.Client
	obj      *hcloud.FloatingIP
	Name     *string
	ID       *int64
	Labels   map[string]string `description:"The labels associated with the floatingip"`
	Created  *time.Time        `description:"The time the floatingip was created"`
	Location *string           `description:"The location of the floatingip"`
}

func (r *FloatingIP) Remove(ctx context.Context) error {
	_, err := r.client.FloatingIP.Delete(ctx, r.obj)

	return err
}

func (r *FloatingIP) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *FloatingIP) String() string {
	return *r.Name
}
