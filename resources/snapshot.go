package resources

import (
	"context"

	"github.com/cgroschupp/hetzner-nuke/pkg/nuke"
	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

const SnapshotResource = "Snapshot"

func init() {
	registry.Register(&registry.Registration{
		Name:     SnapshotResource,
		Scope:    nuke.Account,
		Resource: &Snapshot{},
		Lister:   &SnapshotLister{},
	})
}

type SnapshotLister struct{}

func (l *SnapshotLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	resources := make([]resource.Resource, 0)
	snapshots, _, err := opts.Client.Image.List(ctx, hcloud.ImageListOpts{Type: []hcloud.ImageType{hcloud.ImageTypeSnapshot}})
	if err != nil {
		return resources, err
	}
	for _, snapshot := range snapshots {
		resources = append(resources, &Snapshot{obj: snapshot, Name: &snapshot.Name, Labels: snapshot.Labels, ID: &snapshot.ID})
	}
	return resources, nil
}

type Snapshot struct {
	client *hcloud.Client
	obj    *hcloud.Image
	Name   *string
	ID     *int64
	Labels map[string]string
}

func (r *Snapshot) Remove(ctx context.Context) error {
	_, err := r.client.Image.Delete(ctx, r.obj)

	return err
}

func (r *Snapshot) String() string {
	return *r.Name
}
