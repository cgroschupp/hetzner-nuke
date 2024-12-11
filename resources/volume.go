package resources

import (
	"context"

	"github.com/cgroschupp/hetzner-nuke/pkg/nuke"
	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

const VolumeResource = "Volume"

func init() {
	registry.Register(&registry.Registration{
		Name:     VolumeResource,
		Scope:    nuke.Account,
		Resource: &Volume{},
		Lister:   &VolumeLister{},
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
		resources = append(resources, &Volume{obj: volume, Name: &volume.Name, Labels: volume.Labels, ID: &volume.ID, client: opts.Client})
	}
	return resources, nil
}

type Volume struct {
	client *hcloud.Client
	obj    *hcloud.Volume
	Name   *string
	ID     *int64
	Labels map[string]string
}

func (r *Volume) Remove(ctx context.Context) error {
	_, err := r.client.Volume.Delete(ctx, r.obj)

	return err
}

func (r *Volume) String() string {
	return *r.Name
}
