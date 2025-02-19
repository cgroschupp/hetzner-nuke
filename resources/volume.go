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

const VolumeResource = "Volume"

func init() {
	registry.Register(&registry.Registration{
		Name:      VolumeResource,
		Scope:     nuke.Account,
		Resource:  &Volume{},
		Lister:    &VolumeLister{},
		DependsOn: []string{ServerResource},
	})
}

type VolumeLister struct{}

func (l *VolumeLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	resources := make([]resource.Resource, 0)
	volumes, _, err := opts.Client.Volume.List(ctx, hcloud.VolumeListOpts{})
	if err != nil {
		return resources, err
	}
	for _, volume := range volumes {
		resources = append(resources, &Volume{obj: volume, Name: &volume.Name, Labels: volume.Labels, ID: &volume.ID, client: opts.Client, Created: &volume.Created, Location: &volume.Location.Name})
	}
	return resources, nil
}

type Volume struct {
	client   *hetzner.Client
	obj      *hcloud.Volume
	Name     *string
	ID       *int64
	Labels   map[string]string `description:"The labels associated with the volume"`
	Created  *time.Time        `description:"The time the volume was created"`
	Location *string           `description:"The location of the volume"`
}

func (r *Volume) Remove(ctx context.Context) error {
	_, err := r.client.Volume.Delete(ctx, r.obj)

	return err
}

func (r *Volume) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *Volume) String() string {
	return *r.Name
}
