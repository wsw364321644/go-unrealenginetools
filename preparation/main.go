package main

import (
	"github.com/wsw364321644/unrealenginetools/preparation/deploy"
	"github.com/wsw364321644/unrealenginetools/sharedcode/log"
)

func main() {
	log.Init(&log.LogSettings{"preparation",false,""})

	err:=deploy.DeployConfig()
	if err!=nil{
		log.Error(err)
	}

}

