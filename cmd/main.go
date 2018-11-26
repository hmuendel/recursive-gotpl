package main

import (
	"github.com/golang/glog"
	"github.com/hmuendel/recursive-gotpl/config"
)

var (
	VERSION string = "none"
	COMMIT  string
)

func main() {
	//setting defaults for config values
	defaults := make(map[string]interface{})
	//first thing is to setup logging and read the config
	config.Setup(VERSION, COMMIT, defaults)

	//all finished
	glog.Info("recursive-gotpl finished successfully")
}
