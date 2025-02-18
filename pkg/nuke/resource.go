package nuke

import (
	"github.com/cgroschupp/hetzner-nuke/pkg/hetzner"
	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/sirupsen/logrus"
)

const Account registry.Scope = "account"

type ListerOpts struct {
	Client *hetzner.Client
	Logger *logrus.Entry
}
