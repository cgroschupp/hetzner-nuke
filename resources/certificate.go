package resources

import (
	"context"

	"github.com/cgroschupp/hetzner-nuke/pkg/nuke"
	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
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
		resources = append(resources, &Certificate{obj: certificate, Name: &certificate.Name, Labels: certificate.Labels, ID: &certificate.ID, client: opts.Client})
	}
	return resources, nil
}

type Certificate struct {
	client *hcloud.Client
	obj    *hcloud.Certificate
	Name   *string
	ID     *int64
	Labels map[string]string
}

func (r *Certificate) Remove(ctx context.Context) error {
	_, err := r.client.Certificate.Delete(ctx, r.obj)

	return err
}

func (r *Certificate) String() string {
	return *r.Name
}
