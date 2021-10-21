package sharedcode

import (
	"github.com/wsw364321644/go-botil/log"
	"gopkg.in/ini.v1"
)

func GetEngineConfigPath() string{
	return "config/baseengine.conf"
}
func GetEngineConfig(section string,key string)(value *ini.Key,err error){
	file,err:=ini.ShadowLoad(GetEngineConfigPath())
	if err!=nil{
		log.Error(err)
		return
	}
	s,err:=file.GetSection(section)
	value,err= s.GetKey(key)
	return
}




