package resources

import (
	"context"

	"github.com/cgroschupp/hetzner-nuke/pkg/nuke"
	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

const SSHKeyResource = "SSHKey"

func init() {
	registry.Register(&registry.Registration{
		Name:     SSHKeyResource,
		Scope:    nuke.Account,
		Resource: &SSHKey{},
		Lister:   &SSHKeyLister{},
	})
}

type SSHKeyLister struct{}

func (l *SSHKeyLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	resources := make([]resource.Resource, 0)
	sshkeys, _, err := opts.Client.SSHKey.List(ctx, hcloud.SSHKeyListOpts{})
	if err != nil {
		return resources, err
	}
	for _, sshkey := range sshkeys {
		resources = append(resources, &SSHKey{obj: sshkey, Name: &sshkey.Name, Labels: sshkey.Labels, ID: &sshkey.ID, client: opts.Client})
	}
	return resources, nil
}

type SSHKey struct {
	client *hcloud.Client
	obj    *hcloud.SSHKey
	Name   *string
	ID     *int64
	Labels map[string]string
}

func (r *SSHKey) Remove(ctx context.Context) error {
	_, err := r.client.SSHKey.Delete(ctx, r.obj)

	return err
}

func (r *SSHKey) String() string {
	return *r.Name
}
