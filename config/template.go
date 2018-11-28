package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/golang/glog"
)

type TemplateConfig struct {
	MissingKey string `validate:""`
	SourcePath string `validate:"required,uri"`
	TargetPath string `validate:"omitempty,uri"`
}

//New Login returns an initialized and validated TemplateConfig
func NewTemplateConfig() (*TemplateConfig, error) {
	tc := TemplateConfig{}
	err := tc.Init()
	if err != nil {
		return nil, err
	}
	err = tc.Validate()
	if err != nil {
		return nil, err
	}
	return &tc, nil
}

func (tc *TemplateConfig) Init() error {
	configStore.OnConfigChange(func(e fsnotify.Event) {
		if glog.V(5) {
			glog.Info("LoggingConfig file changed:", e.Name)
			if glog.V(6) {
				glog.Info(e)
			}
		}
		err := tc.configure()
		if err != nil {
			if glog.V(1) {
				glog.Warningf("changed config produce err : %v falling back to previous", err)
			}
		}
	})
	err := tc.configure()
	if err != nil {
		return fmt.Errorf("error initializing log config: %v", err)
	}
	return nil
}

func (tc *TemplateConfig) configure() error {
	var newConfig = TemplateConfig{}
	err := configStore.UnmarshalKey("template", &newConfig)
	if err != nil {
		return err
	}
	err = newConfig.Validate()
	if err != nil {
		return fmt.Errorf("errror setting up template config: %v", err)
	}

	tc.MissingKey = newConfig.MissingKey
	tc.SourcePath = newConfig.SourcePath
	tc.TargetPath = newConfig.TargetPath
	return nil
}

func (tc *TemplateConfig) Validate() error {
	err := validate.Struct(tc)
	if err != nil {
		return fmt.Errorf("error validating template config: %v", err)
	}
	return nil
}
