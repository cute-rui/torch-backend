package test

import (
	"log"
	"testing"
	"torch-client/instance"
)

func TestDigitalOceanUpdate(t *testing.T) {
	var i instance.Instance
	err := instance.NewDigitalOcean(&i)
	if err != nil {
		log.Println(err)
		return
	}

	err = i.UpdateIP()
	if err != nil {
		log.Println(err)
	}
}
