package tools

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func SetLastLash(text string) string {
	last := text[len(text)-1:]

	if last != "/" {
		return text + "/"

	}
	return text

}

func SetFirstLash(text string) string {
	if text == "" {
		return "/"
	}
	first := text[0:1]

	if first != "/" {
		return "/" + text

	}
	return text

}

func CopyDir(source, destination string) error {
	var err error = filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		var relPath string = strings.Replace(path, source, "", 1)
		if relPath == "" {
			return nil
		}
		if info.IsDir() {
			return os.Mkdir(filepath.Join(destination, relPath), 0755)
		} else {
			var data, err1 = ioutil.ReadFile(filepath.Join(source, relPath))
			if err1 != nil {
				return err1
			}
			return ioutil.WriteFile(filepath.Join(destination, relPath), data, 0777)
		}
	})
	return err
}
