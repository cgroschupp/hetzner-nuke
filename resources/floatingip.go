package resources

import (
	"context"

	"github.com/cgroschupp/hetzner-nuke/pkg/nuke"
	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
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
		resources = append(resources, &FloatingIP{obj: floatingIP, Name: &floatingIP.Name, Labels: floatingIP.Labels, ID: &floatingIP.ID})
	}
	return resources, nil
}

type FloatingIP struct {
	client *hcloud.Client
	obj    *hcloud.FloatingIP
	Name   *string
	ID     *int64
	Labels map[string]string
}

func (r *FloatingIP) Remove(ctx context.Context) error {
	_, err := r.client.FloatingIP.Delete(ctx, r.obj)

	return err
}

func (r *FloatingIP) String() string {
	return *r.Name
}
