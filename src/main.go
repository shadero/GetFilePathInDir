package main

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"io/ioutil"
	"net/http"
	"path/filepath"
)

type SearchDirConf struct {
	Path	string		"json:path"
	Port	string		"json:port"	
}

type FileConf struct {
	Path	string		"json:path"
	Date	int64		"json:date"
}
func loadConfig(dir, name string) (SearchDirConf, error) {
	var conf SearchDirConf
	confPath := dir + name

	j, err := ioutil.ReadFile(confPath)
	if err != nil {
		fmt.Println(name + "が確認できないため生成します。")
		dConfJ, err := json.Marshal(SearchDirConf{Path: dir, Port: "1323"})
		if err != nil {
			return conf, err
		}
		if err := ioutil.WriteFile(confPath, dConfJ, 0664); err != nil {
			return conf, err
		}
		j, err = ioutil.ReadFile(confPath)
		if err != nil {
			return conf, err
		}	
		fmt.Println(name + "を生成しました。")
	}
	if err := json.Unmarshal(j, &conf); err != nil {
		return conf, err
	}
	return conf, nil
}

func loadFilePath(dirPath string) ([]FileConf) {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		panic(err)
	}

	var fileConfs []FileConf	
	for _, file := range files {
		p := filepath.Join(dirPath, file.Name())

		if(file.IsDir()){
			fileConfs = append(
				fileConfs,
				loadFilePath(p)...
			)
			continue
		}

		fileConfs = append(
			fileConfs,
			FileConf {
				Path: p,
				Date: file.ModTime().Unix(),
			},
		)
	}
	return fileConfs
}

func main(){
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	conf, err := loadConfig("./", "config.json")
	if err != nil {
		panic(err)
	}
	
	e.GET ("/", func(c echo.Context) error {
		/*files, err := ioutil.ReadDir(conf.Path)
		if err != nil {
			panic(err)
		}*/
		/*var fileConfs []FileConf	
		for _, file := range files {
			if(file.IsDir()){
				
				continue
			}
			fileConfs = append(
				fileConfs,
				FileConf {
					Path: filepath.Join(conf.Path, file.Name()),
					Date: file.ModTime().Unix(),
				},
			)
		}*/
		return c.JSON(http.StatusOK, loadFilePath(conf.Path))
	})

	e.Start(":" + conf.Port)
}
