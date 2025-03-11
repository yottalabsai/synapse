package utils

import (
	"crypto/tls"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"synapse/log"
	"time"
)

func GetProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", os.ErrNotExist
		}
		dir = parent
	}
}

func CheckIngressURLHealth(url string) bool {
	client := http.Client{
		Timeout: 5 * time.Second,
		// todo: remove this when we have a valid SSL certificate
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	resp, err := client.Get(url)
	if err != nil {
		fmt.Printf("Error checking URL: %v\n", err)
		return false
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Log.Errorw("Error closing response body", zap.Error(err))
		}
	}(resp.Body)

	return resp.StatusCode == http.StatusOK
}

func GenerateRequestId() string {
	return fmt.Sprintf("req-%d", time.Now().UnixNano())
}
