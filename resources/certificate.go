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

const CertificateResource = "Certificate"

func init() {
	registry.Register(&registry.Registration{
		Name:     CertificateResource,
		Scope:    nuke.Account,
		Resource: &Certificate{},
		Lister:   &CertificateLister{},
	})
}

type CertificateLister struct{}

func (l *CertificateLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	resources := make([]resource.Resource, 0)
	certificates, _, err := opts.Client.Certificate.List(ctx, hcloud.CertificateListOpts{})
	if err != nil {
		return resources, err
	}
	for _, certificate := range certificates {
		resources = append(resources, &Certificate{obj: certificate, Name: &certificate.Name, Labels: certificate.Labels, ID: &certificate.ID, client: opts.Client, Created: &certificate.Created})
	}
	return resources, nil
}

type Certificate struct {
	client  *hetzner.Client
	obj     *hcloud.Certificate
	Name    *string
	ID      *int64
	Labels  map[string]string `description:"The labels associated with the certificate"`
	Created *time.Time        `description:"The time the certificate was created"`
}

func (r *Certificate) Remove(ctx context.Context) error {
	_, err := r.client.Certificate.Delete(ctx, r.obj)

	return err
}

func (r *Certificate) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *Certificate) String() string {
	return *r.Name
}
