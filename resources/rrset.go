package resources

import (
	"context"

	"github.com/cgroschupp/hetzner-nuke/pkg/hetzner"
	"github.com/cgroschupp/hetzner-nuke/pkg/nuke"
	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

const RRSetResource = "RRSet"

func init() {
	registry.Register(&registry.Registration{
		Name:     RRSetResource,
		Scope:    nuke.Account,
		Resource: &RRSet{},
		Lister:   &RRSetLister{},
	})
}

type RRSetLister struct{}

func (l *RRSetLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	resources := make([]resource.Resource, 0)
	zones, _, err := opts.Client.Zone.List(ctx, hcloud.ZoneListOpts{})
	if err != nil {
		return resources, err
	}
	for _, zone := range zones {
		rrsets, _, err := opts.Client.Zone.ListRRSets(ctx, zone, hcloud.ZoneRRSetListOpts{})
		if err != nil {
			return resources, err
		}
		for _, rrset := range rrsets {
			if rrset.Type == hcloud.ZoneRRSetTypeSOA || rrset.Type == hcloud.ZoneRRSetTypeNS {
				continue
			}
			resources = append(resources, &RRSet{obj: rrset, client: opts.Client, ID: &rrset.ID, Name: &rrset.Name, Labels: rrset.Labels, Type: rrset.Type})
		}

	}
	return resources, nil
}

type RRSet struct {
	client *hetzner.Client
	obj    *hcloud.ZoneRRSet
	Name   *string
	ID     *string
	Type   hcloud.ZoneRRSetType
	Labels map[string]string `description:"The labels associated with the rrset"`
}

func (r *RRSet) Remove(ctx context.Context) error {
	_, _, err := r.client.Zone.DeleteRRSet(ctx, r.obj)

	return err
}

func (r *RRSet) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *RRSet) String() string {
	return *r.Name
}
