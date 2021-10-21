package main

import (
	"github.com/wsw364321644/go-botil/log"
	"github.com/wsw364321644/go-unrealenginetools/preparation/deploy"
)

func main() {
	log.Init(&log.LogSettings{"preparation", false, ""})

	err := deploy.DeployConfig()
	if err != nil {
		log.Error(err)
	}

}
