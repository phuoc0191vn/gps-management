package utilities

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

func LoadEnvFromFile(config interface{}, configPrefix, envPath string) (err error) {
	godotenv.Load(envPath)
	err = envconfig.Process(configPrefix, config)
	return
}

func LoadEnv(config interface{}, prefix string, source string) error {
	if err := LoadEnvFromDir(config, prefix, source); err != nil {
		return LoadEnvFromFile(config, prefix, source)
	}

	return nil
}

func LoadEnvFromDir(config interface{}, configPrefix, dir string) error {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	filePaths := make([]string, 0)
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		filePaths = append(filePaths, filepath.Join(dir, f.Name()))
	}

	if err := godotenv.Load(filePaths...); err != nil {
		return err
	}
	return envconfig.Process(configPrefix, config)
}

func TimeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Printf("%s took %s\n", name, elapsed)
}

func StringInArray(str string, arr []string) bool {
	if len(arr) == 0 {
		return false
	}

	for _, val := range arr {
		if strings.TrimSpace(str) == strings.TrimSpace(val) {
			return true
		}
	}
	return false
}

func GetQuery(req *http.Request, key string) (string, bool) {
	if values, ok := req.URL.Query()[key]; ok && len(values) > 0 {
		return values[0], true
	}
	return "", false
}

func IntInArray(i int, arr []int) bool {
	if len(arr) == 0 {
		return false
	}

	for _, val := range arr {
		if val == i {
			return true
		}
	}
	return false
}

func PostService(endpoint, authToken string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, endpoint, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authToken))
	req.Header.Set("Content-Type", "application/json")

	return http.DefaultClient.Do(req)
}
