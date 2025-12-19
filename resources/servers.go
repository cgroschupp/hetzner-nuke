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

const ServerResource = "Server"

func init() {
	registry.Register(&registry.Registration{
		Name:     ServerResource,
		Scope:    nuke.Account,
		Resource: &Server{},
		Lister:   &ServerLister{},
		DependsOn: []string{
			NetworkResource,
		},
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
		resources = append(resources, &Server{obj: server, Name: &server.Name, Labels: server.Labels, ID: &server.ID, client: opts.Client, Created: &server.Created, Location: &server.Location.Name})
	}
	return resources, nil
}

type Server struct {
	client   *hetzner.Client
	obj      *hcloud.Server
	Name     *string
	ID       *int64
	Labels   map[string]string `description:"The labels associated with the server"`
	Created  *time.Time        `description:"The time the server was created"`
	Location *string           `description:"The location of the server"`
}

func (r *Server) Remove(ctx context.Context) error {
	_, _, err := r.client.Server.DeleteWithResult(ctx, r.obj)

	return err
}

func (r *Server) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *Server) String() string {
	return *r.Name
}
