package test

import (
	"log"
	"testing"
	"torch-client/dns"
)

func TestCloudflare(t *testing.T) {
	var i dns.DNS
	err := dns.NewCloudflare(&i)
	if err != nil {
		log.Println(err)
		return
	}

	err = i.UpdateDNS()
	if err != nil {
		log.Println(err)
		return
	}
}
