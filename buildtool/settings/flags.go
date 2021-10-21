package settings

import (
	"errors"
	"flag"
	"log"
)


type flags struct {

	SILENCE bool
	PLATFORM string
	CONFIG string

	ONLYCOOK bool
	PACK bool
}
var Flags flags
func init() {

	Flags=flags{}

	flag.BoolVar(&Flags.SILENCE, "silence",  false, "without console input")

	plat:=""
	for _,v:=range Platforms{
		plat=v
		break
	}

	flag.StringVar(&Flags.PLATFORM, "plat", plat, GetAllPlatformStr())
	flag.StringVar(&Flags.CONFIG, "conf", Configurations[0], GetAllConf())

	flag.BoolVar(&Flags.ONLYCOOK,"onlycook" , false, "only cook resource")
	flag.BoolVar(&Flags.PACK,"pack" , false, "pack resource")


	flag.Parse()


	str:=ParseConfigurations(Flags.CONFIG)
	if(str==""){
		flag.PrintDefaults()
		log.Panic("wrong conf")
	}
	_,err:=ParsePlatformStr(Flags.PLATFORM)
	if(err!=nil){
		flag.PrintDefaults()
		log.Panic("wrong plat")
	}


}


var Configurations = []string{
	"Debug",
	"Development",
	"Shipping",
}
func GetDefaultConfiguration() string {
	return Configurations[0]
}
func ParseConfigurations(conf string) string{
	for _,c:=range Configurations{
		if c==conf{
			return conf
		}
	}
	return ""
}

func GetAllConf()string{
	root:=""
	for _,c:=range Configurations{
		root+=c
		root+="\n"
	}
	return root
}


type PlatformType int
const (
	Plat_Begin PlatformType = iota
	Plat_WindowsClient
	Plat_WindowsNoEditor
	Plat_WindowsServer
	Plat_Windows
	Plat_LinuxClient
	Plat_LinuxNoEditor
	Plat_LinuxServer
	Plat_Linux
	Plat_End
)
var Platforms = map[PlatformType]string{
	Plat_WindowsNoEditor:"WindowsNoEditor",
	Plat_WindowsServer:"WindowsServer",
	Plat_LinuxServer:"LinuxServer",
}
func GetPlatformStr(p PlatformType) string{
	t,ok:=Platforms[p]
	if(ok){
		return t
	}else{
		return ""
	}
}
func GetAllPlatformStr()string{
	root:=""
	for _,v:=range Platforms{
		root+=v
		root+="\n"
	}
	return root
}
func GetOSStr(p PlatformType) string{
	if(p>=Plat_WindowsClient&&p<=Plat_Windows){
		return "Win64"
	}else if(p>=Plat_LinuxClient&&p<=Plat_Linux){
		return "Linux"
	}
	return ""
}
func ParsePlatformStr(str string)(PlatformType,error){
	for k,v:=range Platforms{
		if v==str{
			return k,nil
		}
	}
	return Plat_Begin,errors.New("plat error")
}

