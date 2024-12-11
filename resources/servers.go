package resources

import (
	"context"

	"github.com/cgroschupp/hetzner-nuke/pkg/nuke"
	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

const ServerResource = "Server"

func init() {
	registry.Register(&registry.Registration{
		Name:     ServerResource,
		Scope:    nuke.Account,
		Resource: &Server{},
		Lister:   &ServerLister{},
	})
}

type ServerLister struct{}

func (l *ServerLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	resources := make([]resource.Resource, 0)
	servers, _, err := opts.Client.Server.List(ctx, hcloud.ServerListOpts{})
	if err != nil {
		return resources, err
	}
	for _, server := range servers {
		resources = append(resources, &Server{obj: server, Name: &server.Name, Labels: server.Labels, ID: &server.ID, client: opts.Client})
	}
	return resources, nil
}

type Server struct {
	client *hcloud.Client
	obj    *hcloud.Server
	Name   *string
	ID     *int64
	Labels map[string]string
}

func (r *Server) Remove(ctx context.Context) error {
	_, _, err := r.client.Server.DeleteWithResult(ctx, r.obj)

	return err
}

func (r *Server) String() string {
	return *r.Name
}
