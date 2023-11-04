package main

import (
	"github.com/go-co-op/gocron"
	"log"
	"time"
	"torch-client/instance"
	"torch-client/utils"
)

var Instance instance.Instance

func main() {
	i, err := instance.NewInstanceByType()
	if err != nil {
		panic(err)
	}
	Instance = i

	s := gocron.NewScheduler(time.UTC)
	_, err = s.Every(utils.Conf.GetString("Frequency")).Do(Check)
	if err != nil {
		return
	}

	s.StartBlocking()
}

func Check() {
	ok, err := utils.Check()
	if err != nil {
		utils.ErrorTimes++
	}

	if ok {
		return
	}

	if !utils.Conf.GetBool("Supervisor") {
		log.Println(`blocked`)
		return
	}

	if err != nil && utils.ErrorTimes > 3 {
		log.Println(`backend down`)
		return
	}

	err = Instance.UpdateIP()
	if err != nil {
		log.Println(`error: `, err.Error())
	}
}
