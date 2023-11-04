package instance

import (
	"errors"
	"sync"
	"torch-client/utils"
)

var InstanceMap = sync.Map{}

type Instance interface {
	UpdateIP() error
	SetupDNS() error
}

type NewInstance func(instance *Instance) error

func NewInstanceByType() (Instance, error) {
	f, ok := InstanceMap.Load(utils.Conf.GetString("Instance.Type"))

	if !ok {
		return nil, errors.New(`instance type not found`)
	}

	nf, ok := f.(NewInstance)
	if !ok {
		return nil, errors.New(`error on assert`)
	}

	var instance Instance
	if err := nf(&instance); err != nil {
		return nil, err
	}

	return instance, nil
}
