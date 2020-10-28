package crontab

import (
	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

func RemoveStatic() error {
	c := cron.New()
	spec := "0 0 1 * * *"
	err := c.AddFunc(spec, func() {
		logrus.Info(1)
		err := RemoveContents("/opt/static/")
		if err != nil {

		}
	})
	if err != nil {
		return err
	}
	c.Start()
	return nil
}

func RemoveContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}