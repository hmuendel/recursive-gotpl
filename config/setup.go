package config

import (
	"flag"
	"fmt"
	"github.com/golang/glog"
	"github.com/spf13/viper"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

//The prefix for which config environment variables are considered
const (
	EnvPrefix     = "RGTPL"
	NoConfigPanic = "panicking cowardly without a config to read"
)

// Setup should be the initial call in the application.
// It configures logging before the config is read and reads the config at the default paths
// or the path configured via environment variables. It logs version and commit hash
func Setup(version, commit string, defaults map[string]interface{}) {

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

	viper.SetEnvPrefix(EnvPrefix)
	err = viper.BindEnv("config", EnvPrefix+"_CONFIG")
	if err != nil {
		fmt.Print(fmt.Errorf("error binding viper conifg env: %s", err))
	}
	file := path.Base(viper.GetString("config"))
	name := strings.TrimSuffix(file, filepath.Ext(file))
	dir := path.Dir(viper.GetString("config"))
	if glog.V(5) {
		glog.Infof("reading config from %s/%s.<extension>", dir, name)
	}
	viper.SetConfigName(name)
	viper.AddConfigPath(dir)
	err = viper.ReadInConfig()
	if err != nil {
		glog.Errorf("Fatal error config file: %s \n", err)
		panic(NoConfigPanic)
	}
	viper.AutomaticEnv()
	viper.WatchConfig()
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
		viper.SetDefault(k, v)
	}
}
