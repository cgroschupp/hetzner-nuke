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

const LoadBalancerResource = "LoadBalancer"

func init() {
	registry.Register(&registry.Registration{
		Name:     LoadBalancerResource,
		Scope:    nuke.Account,
		Resource: &LoadBalancer{},
		Lister:   &LoadBalancerLister{},
		DependsOn: []string{
			NetworkResource,
			CertificateResource,
		},
	})
}

type LoadBalancerLister struct{}

func (l *LoadBalancerLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	resources := make([]resource.Resource, 0)
	loadBalancers, _, err := opts.Client.LoadBalancer.List(ctx, hcloud.LoadBalancerListOpts{})
	if err != nil {
		return resources, err
	}
	for _, loadBalancer := range loadBalancers {
		resources = append(resources, &LoadBalancer{obj: loadBalancer, Name: &loadBalancer.Name, Labels: loadBalancer.Labels, ID: &loadBalancer.ID, client: opts.Client, Created: &loadBalancer.Created, Location: &loadBalancer.Location.Name})
	}
	return resources, nil
}

type LoadBalancer struct {
	client   *hetzner.Client
	obj      *hcloud.LoadBalancer
	Name     *string
	ID       *int64
	Labels   map[string]string `description:"The labels associated with the loadbalancer"`
	Created  *time.Time        `description:"The time the loadbalancer was created"`
	Location *string           `description:"The location of the loadbalancer"`
}

func (r *LoadBalancer) Remove(ctx context.Context) error {
	_, err := r.client.LoadBalancer.Delete(ctx, r.obj)

	return err
}

func (r *LoadBalancer) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *LoadBalancer) String() string {
	return *r.Name
}
