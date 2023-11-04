package instance

import (
	"context"
	"github.com/digitalocean/godo"
	"time"
	"torch-client/dns"
	"torch-client/utils"
)

type DigitalOcean struct {
	DNS        dns.DNS
	client     *godo.Client
	DropletID  int
	ReservedIP string
}

func (do *DigitalOcean) UpdateIP() error {
	retry := 0
	var final error
	for retry < 5 {
		reservedIPs, _, err := do.client.ReservedIPs.List(context.Background(), nil)
		if err != nil {
			retry++
			final = err
			continue
		}

		for i := range reservedIPs {
			if reservedIPs[i].Droplet.ID == do.DropletID {
				_, err = do.client.ReservedIPs.Delete(context.Background(), reservedIPs[i].IP)
				if err != nil {
					break
				}
			}
		}

		break
	}

	if retry == 5 {
		return final
	}
	retry = 0
	time.Sleep(5 * time.Second)
	for retry < 5 {
		ip, _, err := do.client.ReservedIPs.Create(context.Background(), &godo.ReservedIPCreateRequest{DropletID: do.DropletID})
		if err != nil {
			retry++
			final = err
			continue
		}
		do.ReservedIP = ip.IP
		break
	}

	return do.DNS.UpdateDNSForce(do.ReservedIP)
}

func (do *DigitalOcean) SetupDNS() error {
	d, err := dns.NewDNSByType()
	if err != nil {
		return err
	}

	do.DNS = d
	return nil
}

func init() {
	var f NewInstance = NewDigitalOcean
	InstanceMap.Store("DigitalOcean", f)
}

func NewDigitalOcean(instance *Instance) error {
	client := godo.NewFromToken(utils.Conf.GetString("DigitalOcean.Token"))

	i := DigitalOcean{client: client, DropletID: utils.Conf.GetInt("DigitalOcean.ID")}
	err := i.SetupDNS()
	if err != nil {
		return err
	}
	*instance = &i
	return nil
}
