package main

import (
	"flag"
	"github.com/wsw364321644/unrealenginetools/sharedcode"
	"github.com/wsw364321644/go-botil/log"
	"golang.org/x/sys/windows/registry"
	"path/filepath"
)

func main() {
	var check bool
	flag.BoolVar(&check,"check",true,"check privilege")
	flag.Parse()
	log.Init(&log.LogSettings{"regengine",false,""})

	key,err:=sharedcode.GetEngineConfig(sharedcode.GeneralSectionName,sharedcode.GeneralBuildpath)
	if err!=nil{
		log.Error(err)
		return
	}

	newk,_,err :=registry.CreateKey(registry.CURRENT_USER,"SOFTWARE\\Epic Games\\Unreal Engine\\Builds",registry.WRITE)
	if err!=nil{
		log.Error(err)
		return
	}
	newk.SetStringValue("MyCustom",filepath.Join(key.String(),"Engine","Windows"))



	//registry.DeleteKey(registry.CURRENT_USER,"\\SOFTWARE\\Epic Games\\Unreal Engine\\Builds\\MyCustom")

	//exec.Command("reg add \"HKEY_CURRENT_USER\\SOFTWARE\\Epic Games\\Unreal Engine\\Builds" /v \"MyCustom\" /t REG_SZ /d \""+key.String()+"\" /f")
}


