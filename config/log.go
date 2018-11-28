package config

import (
	"flag"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/golang/glog"
)

type LogConfig struct {
	Level   string `validate:"required,numeric"`
	LogDir  string `validate:"file"`
	Vmodule string `validate:""`
}

//New Login returns an initialized and validated LogConfig
func NewLogConfig() (*LogConfig, error) {
	lc := LogConfig{}
	err := lc.Init()
	if err != nil {
		return nil, err
	}
	err = lc.Validate()
	if err != nil {
		return nil, err
	}
	return &lc, nil
}

func (lc *LogConfig) Init() error {
	configStore.OnConfigChange(func(e fsnotify.Event) {
		if glog.V(5) {
			glog.Info("LoggingConfig file changed:", e.Name)
			if glog.V(6) {
				glog.Info(e)
			}
		}
		err := lc.configure()
		if err != nil {
			if glog.V(1) {
				glog.Warningf("changed config produce err : %v falling back to previous", err)
			}
		}
	})
	err := lc.configure()
	if err != nil {
		return fmt.Errorf("error initializing log config: %v", err)
	}
	return nil
}

func (lc *LogConfig) configure() error {
	var newConfig = LogConfig{}
	err := configStore.UnmarshalKey("log", &newConfig)
	if err != nil {
		return err
	}
	err = newConfig.Validate()
	if err != nil {
		return fmt.Errorf("errror setting up log config: %v", err)
	}
	err = newConfig.configureGlog()
	if err != nil {
		_ = lc.configureGlog()
		return err
	}
	lc.Level = newConfig.Level
	lc.LogDir = newConfig.LogDir
	lc.Vmodule = newConfig.Vmodule
	return nil
}

func (lc *LogConfig) Validate() error {
	err := validate.Struct(lc)
	if err != nil {
		return fmt.Errorf("error validating log config: %v", err)
	}
	return nil
}

func (lc *LogConfig) configureGlog() error {

	err := flag.Set("logtostderr", "true")
	if err != nil {
		return err
	}
	err = flag.Set("v", "0")
	if err != nil {
		return err
	}
	err = flag.Set("log_dir", "")
	if err != nil {
		return err
	}

	flag.Parse()
	return nil
}
