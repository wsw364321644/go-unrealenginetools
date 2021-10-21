package builder

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/wsw364321644/go-botil"
	"github.com/wsw364321644/go-botil/log"
	"github.com/wsw364321644/unrealenginetools/buildtool/settings"
	"github.com/wsw364321644/unrealenginetools/sharedcode"
	"io/ioutil"

	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

var ToolchainUrl="http://cdn.unrealengine.com/CrossToolchain_Linux/v10_clang-5.0.0-centos7.zip"

type toolchainstep int
const(
	checkurl  toolchainstep = iota
	checkrootpath
	checktoolchainfolder
	checktoolchainzip
	downloadtoolchain
	setupzip
	setuptoolchain
)
func BuildProject()error{
	PackFlag:=settings.Flags.PACK
	buildConfiguration:=settings.Flags.CONFIG
	UATFlag := true
	OnlyCookFlag:=settings.Flags.ONLYCOOK

	plat,_:=settings.ParsePlatformStr(settings.Flags.PLATFORM)

	osstr:=settings.GetOSStr(plat)
	if(!settings.Flags.SILENCE){
		for plat := settings.Plat_Begin + 1; plat < settings.Plat_End; plat++ {
			platstr := settings.GetPlatformStr(plat)
			if (platstr != "") {
				fmt.Printf("%d-%s\n", plat, platstr)
			}
		}
		indexstr := botil.CheckedScanfln("choose plat index:", func(input string) bool {
			i, err := strconv.ParseInt(input, 10, 64)
			if (err == nil && settings.GetPlatformStr(settings.PlatformType(i)) != "") {
				return true
			}
			return false
		})
		index, _ := strconv.ParseInt(indexstr, 10, 64)
		plat = settings.PlatformType(index)
		osstr = settings.GetOSStr(plat)
		if (osstr == "Linux" && botil.GetScanBoolFlag("Setup  Toolchain(n/y):", false)) {
			checker := NewToolchainChecker(ToolchainUrl)
			err := checker.checkToolchain()
			if (err != nil) {
				return err
			}
			fmt.Println("Toolchain check success")
		}

		for i, conf := range settings.Configurations {
			fmt.Printf("%d-%s\n", i, conf)
		}
		indexstr = botil.CheckedScanfln("choose conf index:", func(input string) bool {
			i, err := strconv.ParseInt(input, 10, 64)
			if (err == nil && int64(len(settings.Configurations)-1) >= i) {
				return true
			}
			return false
		})
		index, _ = strconv.ParseInt(indexstr, 10, 64)
		buildConfiguration = settings.Configurations[index]

		UATFlag=botil.GetScanBoolFlag("compile UAT AutomationToolLauncher(y/n):",true)

		OnlyCookFlag = botil.GetScanBoolFlag("onlycook (n/y):", false)
		PackFlag = false
		if (!OnlyCookFlag) {
			PackFlag = botil.GetScanBoolFlag("Pack Resource (y/n):", true)
		}
	}

	v,err:=sharedcode.GetEngineConfig(sharedcode.GeneralSectionName,sharedcode.GeneralEngineprojectPath)
	if (err != nil ){
		return err
	}
	exepath:=filepath.Join(v.String(),"Engine","Build","BatchFiles","RunUAT.bat")

	v,err=sharedcode.GetEngineConfig(sharedcode.GeneralSectionName,sharedcode.GeneralGameprojectpath)
	if (err != nil ){
		return err
	}
	projectfilename:=""
	rd, err := ioutil.ReadDir(v.String())
	for _, fi := range rd {
		if !fi.IsDir() {
			if(filepath.Ext(fi.Name())==".uproject"){
				projectfilename=fi.Name()
				break;
			}
		}
	}
	projectpath:=filepath.Join(v.String(),projectfilename)

	v,err=sharedcode.GetEngineConfig(sharedcode.GeneralSectionName,sharedcode.GeneralBuildpath)
	if (err != nil ){
		return err
	}
	archivedirectory:=filepath.Join(v.String(),buildConfiguration)
	archivefolder:=filepath.Join(archivedirectory, settings.GetPlatformStr(plat))
	os.RemoveAll(archivefolder)

	cmd := exec.Command("cmd.exe")
	outpipe,_:=cmd.StdoutPipe()
	errpipe,_:=cmd.StderrPipe()
	cmd.Env=os.Environ()
	input, _ := cmd.StdinPipe()

	c := make(chan os.Signal, 1)
	cwd,err:=os.Getwd()
	if(err!=nil) {return err}
	logfile, err := os.OpenFile(filepath.Join(cwd,"UATbuild.log"), os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0755)
	//logerrfile, err := os.OpenFile(filepath.Join(cwd,"llbuilderr.log"), os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0755)
	if err == nil {
		var logmux sync.Mutex
		go func() {
			for{
				select {
				case <-c:
					return
				case <-time.After(time.Second):
					{
						size := 32 * 1024
						buf := make([]byte, size)
						nr, er := outpipe.Read(buf)
						if nr > 0 {
							nw, ew := os.Stdout.Write(buf[0:nr])
							logmux.Lock()
							nw, ew = logfile.Write(buf[0:nr])
							logmux.Unlock()
							if ew != nil {
								return
							}
							if nr != nw {
								return
							}
						}
						if er != nil {
							if er != io.EOF {
								return
							}
						}
					}
				}
			}
		}()
		go func() {
			for{
				select {
				case <-c:
					return
				case <-time.After(time.Second):
					{
						size := 32 * 1024
						buf := make([]byte, size)
						nr, er := errpipe.Read(buf)
						if nr > 0 {
							nw, ew := os.Stderr.Write(buf[0:nr])
							logmux.Lock()
							nw, ew = logfile.Write(buf[0:nr])
							logmux.Unlock()
							if ew != nil {
								return
							}
							if nr != nw {
								return
							}
						}
						if er != nil {
							if er != io.EOF {
								return
							}
						}
					}
				}
			}
		}()
	}

	cmd.Start()
	fmt.Fprint(input, exepath+" BuildCookRun ")
	fmt.Fprintf(input, " -project=\"%s\" ",projectpath)
	fmt.Fprint(input, " -noP4 ")
	//fmt.Fprint(input, " -platform="+osstr+" -serverplatform="+osstr)
	fmt.Fprint(input, " -clientconfig="+buildConfiguration+" -serverconfig="+buildConfiguration)
	if(UATFlag){
		fmt.Fprint(input, " -compile ")
	}else{
		fmt.Fprint(input, " -nocompile ")
	}
	//fmt.Fprint(input, " -utf8output -build -cook -allmaps -SkipCookingEditorContent -package -compressed -cmdline=\" -Messaging\"")
	fmt.Fprint(input, " -utf8output -build -cook -allmaps -package -compressed -cmdline=\" -Messaging\"")
	if(PackFlag){
		fmt.Fprint(input, " -pak ")
	}
	switch plat{
	case settings.Plat_LinuxServer:
		fallthrough
	case settings.Plat_WindowsServer:
		fmt.Fprint(input, " -server  -noclient -serverplatform="+osstr)
	case settings.Plat_WindowsClient:
		fmt.Fprintf(input, " -platform=%s -targetplatform=%s"+osstr)
	}
	if(OnlyCookFlag){
		fmt.Fprint(input," -skipstage ")
	}else{
		fmt.Fprint(input," -stage ")
		fmt.Fprintf(input," -archive -archivedirectory=\"%s\" ",archivedirectory)
	}
	fmt.Fprintln(input,"")
	fmt.Fprintln(input,"exit")
	//go cmd.Wait()
	//fmt.Fprintln(input,botil.Scanfln())
	cmd.Wait()
	c <- syscall.Signal(10)

	if(!OnlyCookFlag) {
		myarchivefolder:=archivefolder
		if(!PackFlag){
			myarchivefolder+="_notpaked"
		}
		if(myarchivefolder!=archivefolder) {
			err = os.RemoveAll(myarchivefolder)
			if (err != nil) {
				return err
			}
			retrytime := 0
			for {
				err = os.Rename(archivefolder, myarchivefolder)
				if (err != nil) {
					if (retrytime > 5) {
						return err
					} else {
						time.Sleep(time.Second * 5)
						retrytime++
					}
				} else {
					break;
				}
			}
		}
	}


	return nil
}
type ToolchainChecker struct{
	ToolchainUrl string
	ToolchainName string
	ToolchainRootPath string
	ToolchainPath string
	toolchainstatus toolchainstep
}
func NewToolchainChecker(ToolchainUrl string)(checker *ToolchainChecker){
	checker=new(ToolchainChecker)
	checker.ToolchainUrl=ToolchainUrl
	checker.toolchainstatus=checkurl
	return checker
}
func(self *ToolchainChecker) checkToolchain()error{
	v,err:=sharedcode.GetEngineConfig(sharedcode.GeneralSectionName,sharedcode.GeneralEngineprojectPath)
	if (err != nil ){
		return err
	}
	enginepath:=v.String()
	switch self.toolchainstatus {
	case checkurl:
		ToolchainZipName := path.Base(ToolchainUrl)
		fileSuffix := path.Ext(ToolchainZipName)
		if(fileSuffix!=".zip"){return errors.New("error url")}
		self.ToolchainName=strings.TrimSuffix(ToolchainZipName,fileSuffix)
		self.toolchainstatus=checkrootpath
		return self.checkToolchain()
	case checkrootpath:

		self.ToolchainRootPath=filepath.Join(filepath.Dir(enginepath),"Linux_CrossCompileToolChain")
		_, err = os.Stat(self.ToolchainRootPath)
		if(err==nil){
			self.toolchainstatus=checktoolchainfolder
			return self.checkToolchain()
		}else if os.IsNotExist(err) {
			err := os.Mkdir(self.ToolchainRootPath, os.ModePerm)
			if (err != nil ){return err}
			self.toolchainstatus=downloadtoolchain
			return self.checkToolchain()
		}else if (err != nil ){
			return err
		}
	case checktoolchainfolder:
		self.ToolchainPath=filepath.Join(self.ToolchainRootPath,self.ToolchainName)
		_, err := os.Stat(self.ToolchainPath)
		if(err==nil){
			self.toolchainstatus=setuptoolchain
			return self.checkToolchain()
		}else if os.IsNotExist(err) {
			self.toolchainstatus=checktoolchainzip
			return self.checkToolchain()
		}else if (err != nil ){
			return err
		}
	case checktoolchainzip:
		_, err := os.Stat(filepath.Join(self.ToolchainRootPath,self.ToolchainName+".zip"))
		if(err==nil){
			self.toolchainstatus=setupzip
			return self.checkToolchain()
		}else if os.IsNotExist(err) {
			self.toolchainstatus=downloadtoolchain
			return self.checkToolchain()
		}else if (err != nil ){
			return err
		}
	case downloadtoolchain:
		err:=botil.DownloadFile(ToolchainUrl,self.ToolchainRootPath)
		if(err!=nil){return err}
		self.toolchainstatus=setupzip
		return self.checkToolchain()
	case setupzip:
		_,err:=botil.Unzip(filepath.Join(self.ToolchainRootPath,self.ToolchainName+".zip"),self.ToolchainPath)
		if(err!=nil){return err}
		self.toolchainstatus=setuptoolchain
		return self.checkToolchain()
	case setuptoolchain:
		path:=filepath.Join(self.ToolchainPath,"setup.bat")
		cmd := exec.Command("cmd.exe","/c",path)
		err:=cmd.Run()
		if(err!=nil){return err}

		envfile,err:=os.OpenFile(filepath.Join(self.ToolchainRootPath,"OutputEnvVars.txt"),os.O_RDONLY, os.ModePerm)
		defer envfile.Close()
		if(err!=nil){return err}
		envreader:=bufio.NewReader(envfile)

		cmd = exec.Command("cmd.exe")
		cmd.Stdout = os.Stdout
		cmd.Stderr=os.Stderr
		input, _ := cmd.StdinPipe()
		cmd.Start()
		for {
			line, err := envreader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					break
				}
				return err
			}
			line = strings.TrimRight(line,"\n\r")
			pair:=strings.Split(line,"=")
			if(len(pair)!=2){return errors.New("err env")}
			fmt.Fprintf(input, " setx %s %s\n",pair[0],pair[1])
		}
		fmt.Fprintln(input,"exit")
		err=cmd.Wait()
		if(err!=nil){return err}

		path=filepath.Join(enginepath,"Setup.bat")
		cmd= exec.Command("cmd.exe","/c",path)
		err=cmd.Run()
		if(err!=nil){return err}

		path=filepath.Join(enginepath,"GenerateProjectFiles.bat")
		cmd= exec.Command("cmd.exe","/c",path)
		err=cmd.Run()
		return err
	}
	return errors.New("toolchain check error")
}




func BuildInstalledEngine()error {
	v,err:=sharedcode.GetEngineConfig(sharedcode.GeneralSectionName,sharedcode.GeneralEngineprojectPath)
	if err!=nil{
		return err
	}
	enginepath:=v.String()
	v,err=sharedcode.GetEngineConfig(sharedcode.GeneralSectionName,sharedcode.GeneralBuildpath)
	if err!=nil{
		return err
	}
	//buildpath:=v.String()

	errch := make(chan error, 1)
	c:=exec.Command(filepath.Join(enginepath,"Engine","Build","BatchFiles","RunUAT.bat"),"BuildGraph",`-target=Make Installed Build Win64`, `-script=`+filepath.Join(enginepath,"Engine","Build","InstalledEngineBuild.xml"),`-set:HostPlatformOnly=true`,`-set:WithDDC=false`,`-set:VS2019=true`)
	stdout,err:=c.StdoutPipe()

	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			log.Infoln(line)
		}
		//select {
		//case <-time.After(time.Second * 1):
		//	fmt.Println("enter")
		//
		//case  <-errch:
		//	return
		//}
	}()

	err=c.Start()
	if err!=nil{
		return err
	}
	err= c.Wait()
	if err!=nil{
		return err
	}
	errch <-err
	return err
}