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

const ZoneResource = "Zone"

func init() {
	registry.Register(&registry.Registration{
		Name:      ZoneResource,
		Scope:     nuke.Account,
		Resource:  &Zone{},
		Lister:    &ZoneLister{},
		DependsOn: []string{RRSetResource},
	})
}

type ZoneLister struct{}

func (l *ZoneLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	resources := make([]resource.Resource, 0)
	zones, _, err := opts.Client.Zone.List(ctx, hcloud.ZoneListOpts{})
	if err != nil {
		return resources, err
	}
	for _, zone := range zones {
		resources = append(resources, &Zone{obj: zone, client: opts.Client, ID: &zone.ID, Name: &zone.Name, Labels: zone.Labels, Created: &zone.Created})
	}
	return resources, nil
}

type Zone struct {
	client  *hetzner.Client
	obj     *hcloud.Zone
	Name    *string
	ID      *int64
	Labels  map[string]string `description:"The labels associated with the zone"`
	Created *time.Time        `description:"The time the zone was created"`
}

func (r *Zone) Remove(ctx context.Context) error {
	_, _, err := r.client.Zone.Delete(ctx, r.obj)

	return err
}

func (r *Zone) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *Zone) String() string {
	return *r.Name
}
