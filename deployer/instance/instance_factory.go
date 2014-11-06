package instance

import (
	boshlog "github.com/cloudfoundry/bosh-agent/logger"
	boshsys "github.com/cloudfoundry/bosh-agent/system"
	bmcloud "github.com/cloudfoundry/bosh-micro-cli/cloud"
	bmagentclient "github.com/cloudfoundry/bosh-micro-cli/deployer/agentclient"
	bmas "github.com/cloudfoundry/bosh-micro-cli/deployer/applyspec"
)

type Factory interface {
	Create(string, string, bmcloud.Cloud) Instance
}

type instanceFactory struct {
	agentClientFactory     bmagentclient.Factory
	templatesSpecGenerator TemplatesSpecGenerator
	applySpecFactory       bmas.Factory
	fs                     boshsys.FileSystem
	logger                 boshlog.Logger
}

func NewInstanceFactory(
	agentClientFactory bmagentclient.Factory,
	templatesSpecGenerator TemplatesSpecGenerator,
	applySpecFactory bmas.Factory,
	fs boshsys.FileSystem,
	logger boshlog.Logger,
) Factory {
	return &instanceFactory{
		agentClientFactory:     agentClientFactory,
		templatesSpecGenerator: templatesSpecGenerator,
		applySpecFactory:       applySpecFactory,
		fs:                     fs,
		logger:                 logger,
	}
}

func (f *instanceFactory) Create(vmCID string, mbusURL string, cloud bmcloud.Cloud) Instance {
	agentClient := f.agentClientFactory.Create(mbusURL)

	return NewInstance(
		vmCID,
		agentClient,
		cloud,
		f.templatesSpecGenerator,
		f.applySpecFactory,
		mbusURL,
		f.fs,
		f.logger,
	)
}
