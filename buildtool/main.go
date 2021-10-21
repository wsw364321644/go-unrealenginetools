package main

import (
	"fmt"
	"github.com/wsw364321644/go-botil"
	"github.com/wsw364321644/unrealenginetools/buildtool/builder"
	"github.com/wsw364321644/go-botil/log"

	"strconv"
)

type BuildType int
const (
	Build_Game BuildType = iota
	Build_Engine
	Build_End
)
var ActionStrMap = map[BuildType]string{
	Build_Game:"Build_Game",
	Build_Engine:"Build_Engine",
}
func GetBuildTypeStr(a BuildType) string {
	t, ok := ActionStrMap[a]
	if (ok) {
		return t
	} else {
		return ""
	}
}

func main() {
	log.Init(&log.LogSettings{"builder",false,""})

	buildtype:=Build_Game

	for buildtype := Build_Game; buildtype <Build_End; buildtype++ {
		actionstr :=GetBuildTypeStr(buildtype)
		if (actionstr != "") {
			fmt.Printf("%d-%s\n", buildtype, actionstr)
		}
	}
	indexstr := botil.CheckedScanfln("choose build type index:", func(input string) bool {
		if (input == "") {
			return true
		}
		i, err := strconv.ParseInt(input, 10, 64)
		if (err == nil && GetBuildTypeStr(BuildType(i)) != "") {
			return true
		}
		return false
	})
	index, err := strconv.ParseInt(indexstr, 10, 64)
	if (err == nil) {
		buildtype = BuildType(index)
	}

	switch buildtype {
	case Build_Game:
		err=builder.BuildProject()
		if err!=nil{
			log.Error(err)
		}
	case Build_Engine:
		err=builder.BuildInstalledEngine()
		if err!=nil{
			log.Error(err)
		}
	}
}


