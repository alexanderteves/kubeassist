package request

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/alexanderteves/kubeassist/pkg/config"
	"io/ioutil"
	"net/http"
	"os"
)

// GetApiData is a simple HTTPS client that uses data from kcfg to send
// requests to an Kubernetes API.
func GetApiData(kcfg config.Kubeconfig, path string) (apiData []byte, err error) {
	connInfo, err := kcfg.GetConnectionInfo()
	if err != nil {
		return apiData, err
	}

	certPool := x509.NewCertPool()
	caCert, err := os.ReadFile(connInfo.CA)
	if err != nil {
		return apiData, err
	}

	certPool.AppendCertsFromPEM(caCert)
	config := &tls.Config{RootCAs: certPool}
	tr := &http.Transport{TLSClientConfig: config}

	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("GET", connInfo.Server+path, nil)
	if err != nil {
		return apiData, err
	}

	// TODO: Use cert and key if present
	req.Header.Add("Authorization", "Bearer "+connInfo.Token)
	resp, err := client.Do(req)
	if err != nil {
		return apiData, err
	}
	defer resp.Body.Close()

	apiData, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return apiData, err
	}

	return apiData, nil
}
