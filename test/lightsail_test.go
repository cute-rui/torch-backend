package test

import (
	"log"
	"testing"
	"torch-client/instance"
	"torch-client/utils"
)

func TestLightsailUpdate(t *testing.T) {
	var i instance.Instance
	err := instance.NewLightsail(&i)
	if err != nil {
		log.Println(err)
		return
	}

	err = i.UpdateIP()
	if err != nil {
		log.Println(err)
	}
}

func TestGetIP(t *testing.T) {
	log.Println(utils.GetIP())
}

func TestCheck(t *testing.T) {
	log.Println(utils.Check())
}
