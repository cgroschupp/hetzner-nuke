package nuke

import (
	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/sirupsen/logrus"
)

const Account registry.Scope = "account"

type ListerOpts struct {
	Client *hcloud.Client
	Logger *logrus.Entry
}
