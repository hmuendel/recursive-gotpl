package config_test

import (
	"github.com/hmuendel/recursive-gotpl/config"
	"os"
	"reflect"
	"testing"
	"time"
)

var LogConfigTests = []struct {
	name   string
	config string
	output *config.TemplateConfig
}{
	{"only_source", `
template:
  sourcePath: "/var/tmp/"
`,
		&config.TemplateConfig{
			SourcePath: "/var/tmp/",
		}},
	{"wrong_source", `
template:
  sourcePath: 123
`,
		nil},
	{"wrong_target", `
template:
  sourcePath: "/var/tmp/test"
  targetPath: 123
`,
		nil},
}

func TestParseTemplateConfigs(t *testing.T) {
	prefix := envPrefix()
	err := cleanEnv(prefix)
	if err != nil {
		t.Fatalf("unable to clear environment: %v", err)
	}
	err = os.Setenv(prefix+"_LOG_LEVEL", "10")
	if err != nil {
		t.Fatalf("unable to environment variable: %v", err)
	}
	for _, testConfigCase := range LogConfigTests {
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
			lc, err := config.NewTemplateConfig()
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
