package config_test

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"
	"time"
)

var dirCleaner = make([]func(), 0)

//redirects stderr and run function f for a duration. returns all lines captured during this time
func CaptureStderrLines(duration time.Duration, f func()) ([]string, error) {
	lines := make([]string, 0)
	oldStdErr := os.Stderr
	readFile, writeFile, err := os.Pipe()
	if err != nil {
		return nil, err
	}
	os.Stderr = writeFile
	defer func() {
		os.Stderr = oldStdErr
	}()
	c := make(chan string, 1)
	go func() {
		scanner := bufio.NewScanner(readFile)
		for scanner.Scan() {
			var line = scanner.Text()
			c <- line
		}
	}()
	f()
	for {
		select {
		case line := <-c:
			lines = append(lines, line)
		case <-time.After(duration):
			err = writeFile.Close()
			if err != nil {
				fmt.Println("error closing stderr redirect", err)
			}
			return lines, nil
		}
	}
}

func createConfigInTmp(content string) (string, error) {
	tempDir, err := ioutil.TempDir("", "recursive-gotpl-test")
	f := func() {
		fmt.Println("deleting dir", tempDir)
		err := os.RemoveAll(tempDir)
		if err != nil {
			fmt.Printf("Error removing tempdir: %s", tempDir)
		}
	}
	dirCleaner = append(dirCleaner, f)

	if err != nil {
		return "", err
	}
	config := path.Join(tempDir, "config.yaml")
	err = ioutil.WriteFile(config, []byte(content), 0644)
	if err != nil {
		return "", err
	}
	return config, nil
}

func cleanup() {
	fmt.Println("startiing cleanup")
	for _, f := range dirCleaner {
		f()
	}
}

func TestMain(m *testing.M) {
	m.Run()
	cleanup()
}
