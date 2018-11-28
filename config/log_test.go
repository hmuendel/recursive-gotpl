package config_test

import (
	"github.com/hmuendel/recursive-gotpl/config"
	"os"
	"reflect"
	"testing"
	"time"
)

var configTests = []struct {
	name   string
	config string
	output *config.LogConfig
}{
	{"only_level", `
log:
  level: "42"
`,
		&config.LogConfig{
			Level: "42",
		}},
	{"wrong_level", `
log:
  level: "a"
`,
		nil},
}

func TestParseConfigs(t *testing.T) {
	prefix := envPrefix()
	err := cleanEnv(prefix)
	if err != nil {
		t.Fatalf("unable to clear environment: %v", err)
	}
	for _, testConfigCase := range configTests {
		t.Run(testConfigCase.name, func(t *testing.T) {
			tempFile, err := createConfigInTmp(testConfigCase.config)
			if err != nil {
				t.Fatalf("error creating test tempdir: %s", err)
			}
			cfg := make(map[string]interface{})
			cfg["config"] = tempFile
			_, err = CaptureStderrLines(1*time.Millisecond, func() {
				config.Setup("0.1.0", "deadbeaf", prefix, cfg)
			})
			if err != nil {
				t.Fatal(err)
			}
			lc, err := config.NewLogConfig()
			if testConfigCase.output == nil {
				if err == nil {
					t.Errorf("did not error for config %v", testConfigCase.config)
				}
				//todo check for exact validation error types
			} else {
				if err != nil {
					t.Errorf("error: %v for config %v", err, testConfigCase.config)
				}
				if !reflect.DeepEqual(lc, testConfigCase.output) {
					t.Errorf("actual output : %#v does not match expected: %#v", lc, testConfigCase.output)
				}
			}

		})
	}
}

func TestParseConfigsWithHighLogLevel(t *testing.T) {
	prefix := envPrefix()
	err := cleanEnv(prefix)
	if err != nil {
		t.Fatalf("unable to clear environment: %v", err)
	}
	err = os.Setenv(prefix+"_LOG_LEVEL", "10")
	for _, testConfigCase := range configTests {
		t.Run(testConfigCase.name, func(t *testing.T) {
			tempFile, err := createConfigInTmp(testConfigCase.config)
			if err != nil {
				t.Fatalf("error creating test tempdir: %s", err)
			}
			cfg := make(map[string]interface{})
			if err != nil {
				t.Fatalf("error setting config path variable: %s", err)
			}
			cfg["config"] = tempFile
			_, err = CaptureStderrLines(1*time.Millisecond, func() {
				config.Setup("0.1.0", "deadbeaf", prefix, cfg)
			})
			if err != nil {
				t.Fatal(err)
			}
			lc, err := config.NewLogConfig()
			if testConfigCase.output == nil {
				if err == nil {
					t.Errorf("did not error for config %v", testConfigCase.config)
				}
				//todo check for exact validation error types
			} else {
				if err != nil {
					t.Errorf("error: %v for config %v", err, testConfigCase.config)
				}
				if !reflect.DeepEqual(lc, testConfigCase.output) {
					t.Errorf("actual output : %#v does not match expected: %#v", lc, testConfigCase.output)
				}
			}
		})
	}
}
