package app

import (
	"encoding/json"
	"os"
)

type App struct {
}

func New() *App {
	return &App{}
}

func Build() error {
	navi, err := genNavi()
	if err != nil {
		return err
	}

	naviBytes, err := json.Marshal(navi)
	if err != nil {
		return err
	}

	_ = os.MkdirAll("public/api/", os.ModePerm)
	err = os.WriteFile("public/api/navi.json", naviBytes, os.ModePerm)
	if err != nil {
		return err
	}
	return nil

}
