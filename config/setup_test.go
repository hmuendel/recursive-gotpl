package config_test

import (
	"github.com/hmuendel/recursive-gotpl/config"
	"os"
	"path"
	"strings"
	"testing"
	"testing/quick"
	"time"
)

func TestFailOnNoConfig(t *testing.T) {
	prefix := envPrefix()
	err := cleanEnv(prefix)
	if err != nil {
		t.Fatalf("unable to clear environment: %v", err)
	}
	cfg := make(map[string]interface{})
	tempFile, err := createConfigInTmp("")
	if err != nil {
		t.Fatalf("error creating test tempdir: %s", err)
	}
	cfg["config"] = path.Join(tempFile, "non-existent")
	// preparing recovery from expected panic
	defer func() {
		if r := recover(); r == nil {
			t.Error("did not panic without a config file")
		} else {
			panicString, ok := r.(string)
			if !ok {
				t.Error("did not receive string panic")
			}
			if panicString != config.NoConfigPanic {
				t.Errorf("unexpected panic string: '%s' instead of '%s'", panicString, config.NoConfigPanic)
			}
		}
	}()

	//running test function
	_, err = CaptureStderrLines(1*time.Millisecond, func() {
		config.Setup("0.1.0", "deadbeaf", prefix, cfg)
	})
	if err != nil {
		t.Fatal(err)
	}
	if err != nil {
		t.Fatalf("cannot crate tempdir: %v", err)
	}
}

func TestDontFail(t *testing.T) {
	prefix := envPrefix()
	err := cleanEnv(prefix)
	if err != nil {
		t.Fatalf("unable to clear environment: %v", err)
	}
	cfg := make(map[string]interface{})
	tempFile, err := createConfigInTmp("foo: bar")
	if err != nil {
		t.Fatalf("error creating test tempdir: %s", err)
	}
	cfg["config"] = tempFile
	//running test function
	_, err = CaptureStderrLines(1*time.Millisecond, func() {
		config.Setup("0.1.0", "deadbeaf", prefix, cfg)
	})

	if err != nil {
		t.Fatal(err)
	}
	if err != nil {
		t.Fatalf("cannot crate tempdir: %v", err)
	}
}

func TestEnvOverDefault(t *testing.T) {
	prefix := envPrefix()
	err := cleanEnv(prefix)
	if err != nil {
		t.Fatalf("unable to clear environment: %v", err)
	}
	cfg := make(map[string]interface{})
	tempFile, err := createConfigInTmp("foo: bar")
	if err != nil {
		t.Fatalf("error creating test tempdir: %s", err)
	}
	cfg["config"] = path.Join(tempFile, "non-existent")
	err = os.Setenv(prefix+"_CONFIG", tempFile)
	if err != nil {
		t.Fatalf("error setting config path variable: %s", err)

	}
	_, err = CaptureStderrLines(1*time.Millisecond, func() {
		config.Setup("0.1.0", "deadbeaf", prefix, cfg)
	})

	if err != nil {
		t.Fatal(err)
	}
	if err != nil {
		t.Fatalf("cannot crate tempdir: %v", err)
	}
}

func TestInitLogShowUp(t *testing.T) {
	prefix := envPrefix()
	err := cleanEnv(prefix)
	if err != nil {
		t.Fatalf("unable to clear environment: %v", err)
	}
	expectedLog := "starting recursive-gotpl"
	cfg := make(map[string]interface{})
	tempFile, err := createConfigInTmp("foo: bar")
	if err != nil {
		t.Fatalf("error creating test tempdir: %s", err)
	}
	cfg["config"] = tempFile
	output, err := CaptureStderrLines(1*time.Millisecond, func() {
		config.Setup("0.1.0", "deadbeaf", prefix, cfg)
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(output) != 2 {
		t.Error("Expected only two output lines got")
		for _, line := range output {
			t.Error(line)
		}
	}
	if !strings.Contains(output[0], expectedLog) {
		t.Error("Did not find expected string: ", expectedLog, " in ", output[0])
	}
}

func TestPreConfigVerbosityLevel(t *testing.T) {
	prefix := envPrefix()
	err := cleanEnv(prefix)
	if err != nil {
		t.Fatalf("unable to clear environment: %v", err)
	}
	var levelTests = []struct {
		level       string
		outputLines int
	}{
		{"1000", 12},
		{"10", 12},
		{"9", 10},
		{"8", 7},
		{"7", 6},
		{"6", 3},
		{"5", 3},
		{"4", 2},
		{"3", 2},
		{"2", 2},
		{"1", 2},
		{"0", 2},
		{"-1000", 2},
		{"foo", 2},
	}

	for _, testCase := range levelTests {
		t.Run("v"+testCase.level, func(t *testing.T) {
			cfg := make(map[string]interface{})
			tempFile, err := createConfigInTmp("foo: bar")
			if err != nil {
				t.Fatalf("error creating test tempdir: %s", err)
			}
			cfg["config"] = tempFile
			err = os.Setenv(prefix+"_LOG_LEVEL", testCase.level)
			if err != nil {
				t.Fatalf("error setting config path variable: %s", err)

			}
			output, err := CaptureStderrLines(1*time.Millisecond, func() {
				config.Setup("0.1.0", "deadbeaf", prefix, cfg)
			})
			if err != nil {
				t.Fatal(err)
			}
			if len(output) != testCase.outputLines {
				t.Errorf("Expected %d output lines for v = %s got %d",
					testCase.outputLines, testCase.level, len(output))
				for _, line := range output {
					t.Error(line)
				}
			}
		})

	}
}

func TestVariousVersionsAndCommitsToShow(t *testing.T) {
	prefix := envPrefix()
	err := cleanEnv(prefix)
	if err != nil {
		t.Fatalf("unable to clear environment: %v", err)
	}
	f := func(version, commit string) bool {
		cfg := make(map[string]interface{})
		tempFile, err := createConfigInTmp("foo: bar")
		if err != nil {
			t.Fatalf("error creating test tempdir: %s", err)
		}
		cfg["config"] = tempFile
		output, err := CaptureStderrLines(1*time.Millisecond, func() {
			config.Setup(version, commit, prefix, cfg)
		})
		if err != nil {
			t.Fatal(err)
		}
		if len(output) != 2 {
			t.Error("Expected only two output lines got")
			for _, line := range output {
				t.Error(line)
			}
			return false

		}
		if !(strings.Contains(output[0], version) && strings.Contains(output[0], commit)) {
			t.Error("Did not find expected version: ", version, " or commits: ", commit, " in ", output[0])
			return false
		}
		return true

	}

	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}
