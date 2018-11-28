package config

import (
	"flag"
	"fmt"
	"github.com/golang/glog"
	"github.com/spf13/viper"
	"gopkg.in/go-playground/validator.v9"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

//The prefix for which config environment variables are considered
const (
	NoConfigPanic = "panicking cowardly without a config to read"
)

var validate *validator.Validate

var configStore *viper.Viper

// Setup should be the initial call in the application.
// It configures logging before the config is read and reads the config at the default paths
// or the path configured via environment variables. It logs version and commit hash
func Setup(version, commit string, EnvPrefix string, defaults map[string]interface{}) {

	//preparing glog flags because we don't configure it via commandline flags and we want to log something
	//before going into the config reading logic which might fail
	err := flag.Set("logtostderr", "true")
	if err != nil {
		fmt.Print(fmt.Errorf("error configuring glog flags: %s", err))
	}
	err = flag.Set("v", "0")
	if err != nil {
		fmt.Print(fmt.Errorf("error configuring glog flags: %s", err))
	}
	flag.Parse()
	glog.Infof("starting recursive-gotpl in version: %s, commit: %s", version, commit)

	//making it possible to get verbose logging before config is read, by specifying env variable
	v := os.Getenv(EnvPrefix + "_LOG_LEVEL")
	if _, err := strconv.Atoi(v); err != nil {
		v = "0"
	}
	glog.Infof("setting verbosity level to %s for pre config logging. This can be changed via %s_LOG_LEVEL",
		v, EnvPrefix)
	err = flag.Set("v", v)
	if err != nil {
		fmt.Print(fmt.Errorf("error configuring glog flags: %s", err))
	}
	flag.Parse()

	if glog.V(7) {
		glog.Info("instantiating viper")
	}
	configStore = viper.New()
	if glog.V(7) {
		glog.Info("instantiating validator")
	}
	validate = validator.New()
	//preparing the config
	if glog.V(7) {
		glog.Info("setting defaults")
		if glog.V(8) {
			glog.Info(defaults)
		}
	}
	setDefaults(defaults)
	if glog.V(9) {
		glog.Info("successfully set defaults")
		if glog.V(10) {
			glog.Infof("setting viper env prefix to %s and binding config env", EnvPrefix)
		}
	}

	configStore.SetEnvPrefix(EnvPrefix)
	err = configStore.BindEnv("config", EnvPrefix+"_CONFIG")
	if err != nil {
		fmt.Print(fmt.Errorf("error binding viper conifg env: %s", err))
	}
	file := path.Base(configStore.GetString("config"))
	name := strings.TrimSuffix(file, filepath.Ext(file))
	dir := path.Dir(configStore.GetString("config"))
	if glog.V(5) {
		glog.Infof("reading config from %s/%s.<extension>", dir, name)
	}
	configStore.SetConfigName(name)
	configStore.AddConfigPath(dir)
	err = configStore.ReadInConfig()
	if err != nil {
		glog.Errorf("Fatal error config file: %s \n", err)
		panic(NoConfigPanic)
	}
	configStore.AutomaticEnv()
	configStore.WatchConfig()
	if glog.V(9) {
		glog.Info("successfully read config")
	}
}

func setDefaults(defaults map[string]interface{}) {
	for k, v := range defaults {
		if glog.V(9) {
			glog.Infof("setting default %s", k)
			if glog.V(10) {
				glog.Info(v)
			}
		}
		configStore.SetDefault(k, v)
	}
}
