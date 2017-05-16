package wapiq

import (
	"os"
)

type Include struct {
	File string
}

func (i *Include) Path() (string, error) {
	var err error
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	pth := wd + "/" + i.File + ".wapiq"
	_, err = os.Stat(pth)
	if err != nil {
		return "", err
	}

	return pth, nil
}
