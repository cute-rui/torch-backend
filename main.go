package main

import (
	"flag"
	"github.com/go-co-op/gocron"
	"log"
	"time"
	"torch-client/instance"
	"torch-client/utils"
)

var Instance instance.Instance
var force bool

func main() {
	flag.BoolVar(&force, "f", false, "force run once")
	flag.Parse()

	i, err := instance.NewInstanceByType()
	if err != nil {
		panic(err)
	}
	Instance = i

	if force {
		log.Println(utils.Check())
		log.Println(Instance.UpdateIP())
		return
	}

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
		log.Println(`check success`)
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
