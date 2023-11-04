package dns

import (
	"context"
	"github.com/cloudflare/cloudflare-go"
	"torch-client/utils"
)

type Cloudflare struct {
	Svc      *cloudflare.API
	ZoneID   string
	RecordID string
}

func (c *Cloudflare) UpdateDNS() error {
	ip, err := utils.GetIP()
	if err != nil {
		return err
	}

	return c.UpdateDNSForce(ip)
}

func (c *Cloudflare) UpdateDNSForce(ip string) error {
	retry := 0
	var final error
	for retry < 5 {
		_, err := c.Svc.UpdateDNSRecord(context.Background(), cloudflare.ZoneIdentifier(c.ZoneID), cloudflare.UpdateDNSRecordParams{ID: c.RecordID, Content: ip})
		if err != nil {
			retry++
			final = err
			continue
		}

		break
	}

	return final
}

func init() {
	var f NewDNS = NewCloudflare
	DNSMap.Store("Cloudflare", f)
}

func NewCloudflare(instance *DNS) error {
	api, err := cloudflare.NewWithAPIToken(utils.Conf.GetString("Cloudflare.Token"))
	if err != nil {
		return err
	}

	var i = Cloudflare{Svc: api, ZoneID: utils.Conf.GetString("Cloudflare.ZoneID"), RecordID: utils.Conf.GetString("Cloudflare.RecordID")}
	*instance = &i

	return nil
}
