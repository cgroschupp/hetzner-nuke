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

const PlacementGroupResource = "PlacementGroup"

func init() {
	registry.Register(&registry.Registration{
		Name:      PlacementGroupResource,
		Scope:     nuke.Account,
		Resource:  &PlacementGroup{},
		Lister:    &PlacementGroupLister{},
		DependsOn: []string{ServerResource},
	})
}

type PlacementGroupLister struct{}

func (l *PlacementGroupLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	resources := make([]resource.Resource, 0)
	placementGroups, _, err := opts.Client.PlacementGroup.List(ctx, hcloud.PlacementGroupListOpts{})
	if err != nil {
		return resources, err
	}
	for _, placementGroup := range placementGroups {
		resources = append(resources, &PlacementGroup{obj: placementGroup, Name: &placementGroup.Name, Labels: placementGroup.Labels, ID: &placementGroup.ID, client: opts.Client, Created: &placementGroup.Created})
	}
	return resources, nil
}

type PlacementGroup struct {
	client  *hetzner.Client
	obj     *hcloud.PlacementGroup
	Name    *string
	ID      *int64
	Labels  map[string]string `description:"The labels associated with the placement group"`
	Created *time.Time        `description:"The time the placement group was created"`
}

func (r *PlacementGroup) Remove(ctx context.Context) error {
	_, err := r.client.PlacementGroup.Delete(ctx, r.obj)

	return err
}

func (r *PlacementGroup) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *PlacementGroup) String() string {
	return *r.Name
}
