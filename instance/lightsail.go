package instance

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/lightsail"
	"torch-client/dns"
	"torch-client/utils"
)

type Lightsail struct {
	Svc *lightsail.Client
	DNS dns.DNS

	IPName       string
	InstanceName string
	Region       string
	Token        string
}

func (l *Lightsail) SetupDNS() error {
	d, err := dns.NewDNSByType()
	if err != nil {
		return err
	}

	l.DNS = d
	return nil
}

func WithRegion(r string) func(*lightsail.Options) {
	return func(o *lightsail.Options) {
		o.Region = r
	}
}

func (l *Lightsail) UpdateIP() error {
	retry := 0
	var final error
	for retry < 5 {
		_, err := l.Svc.ReleaseStaticIp(context.Background(), &lightsail.ReleaseStaticIpInput{StaticIpName: &l.IPName}, WithRegion(utils.Conf.GetString("Lightsail.Region")))
		if err != nil {
			retry++
			final = err
			_ = l.RefreshLogin()
			continue
		}
		break
	}

	if retry == 5 {
		return final
	}
	retry = 0
	for retry < 5 {
		_, err := l.Svc.AllocateStaticIp(context.Background(), &lightsail.AllocateStaticIpInput{StaticIpName: &l.IPName}, WithRegion(utils.Conf.GetString("Lightsail.Region")))
		if err != nil {
			retry++
			final = err
			_ = l.RefreshLogin()
			continue
		}
		break
	}

	if retry == 5 {
		return final
	}
	retry = 0
	for retry < 5 {
		_, err := l.Svc.AttachStaticIp(context.Background(), &lightsail.AttachStaticIpInput{StaticIpName: &l.IPName, InstanceName: &l.InstanceName}, WithRegion(utils.Conf.GetString("Lightsail.Region")))
		if err != nil {
			retry++
			final = err
			_ = l.RefreshLogin()
			continue
		}
		break
	}
	if retry == 5 {
		return final
	}

	return l.DNS.UpdateDNS()
}

func (l *Lightsail) RefreshLogin() error {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(utils.Conf.GetString("Lightsail.AKID"),
			utils.Conf.GetString("Lightsail.SECRET_KEY"),
			utils.Conf.GetString("Lightsail.TOKEN"))),
		config.WithRegion(utils.Conf.GetString("Lightsail.Region")),
	)
	if err != nil {
		return err
	}

	// Create a Lightsail client from just a session.
	l.Svc = lightsail.NewFromConfig(cfg)

	return nil
}

func init() {
	var f NewInstance = NewLightsail
	InstanceMap.Store("Lightsail", f)
}

func NewLightsail(instance *Instance) error {
	var i = Lightsail{IPName: utils.Conf.GetString("Lightsail.IPName"), InstanceName: utils.Conf.GetString("Lightsail.InstanceName"), Region: utils.Conf.GetString("Lightsail.Region")}
	err := i.RefreshLogin()
	if err != nil {
		return err
	}

	err = i.SetupDNS()
	if err != nil {
		return err
	}
	*instance = &i

	return nil
}
