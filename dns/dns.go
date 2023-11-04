package dns

import (
	"errors"
	"sync"
	"torch-client/utils"
)

var DNSMap = sync.Map{}

type DNS interface {
	UpdateDNS() error
	UpdateDNSForce(string) error
}

type NewDNS func(instance *DNS) error

func NewDNSByType() (DNS, error) {
	f, ok := DNSMap.Load(utils.Conf.GetString("DNS.Type"))

	if !ok {
		return nil, errors.New(`dns type not found`)
	}

	nf, ok := f.(NewDNS)
	if !ok {
		return nil, errors.New(`error on assert`)
	}

	var instance DNS
	if err := nf(&instance); err != nil {
		return nil, err
	}

	return instance, nil
}
