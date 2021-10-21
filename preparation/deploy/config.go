package deploy

import (
	"github.com/gobuffalo/packd"
	"github.com/gobuffalo/packr"
	"io"
	"os"
	"path/filepath"
)

func DeployConfig()(err error){
	fileinfo,err:=os.Stat("config")
	if err!=nil||!fileinfo.IsDir() {
		err=nil
		configbox := packr.NewBox("../../config")
		os.Mkdir("config",os.ModeDir)
		configbox.Walk(func(s string, file packd.File) error {

			path:=filepath.Dir(s)
			di,err:=os.Stat(filepath.Join("config",path))
			if(err!=nil||!di.IsDir()){
				os.Mkdir(filepath.Join("config",path),os.ModeDir)
			}


			var out *os.File
			out, err = os.Create(filepath.Join("config",s))
			if err != nil {
				return err
			}
			defer func() {
				if e := out.Close(); e != nil {
					err = e
				}
			}()

			_, err = io.Copy(out, file)
			if err != nil {
				return err
			}

			err = out.Sync()
			if err != nil {
				return err
			}

			err = os.Chmod(filepath.Join("config",s), 0222)
			if err != nil {
				return err
			}

			return err
		})

	}
	return
}








