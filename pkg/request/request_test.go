package request

import (
	"crypto/tls"
	"fmt"
	"github.com/alexanderteves/kubeassist/pkg/config"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequest(t *testing.T) {
	kcfg, err := config.Load("../../test/kube.config")
	if err != nil {
		t.Fatal(err)
	}

	testBody := "Theevn0i"

	// Setup TLS mock server
	ts := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, testBody)
	}))
	cert, err := tls.LoadX509KeyPair("../../test/tls/localhost.crt", "../../test/tls/localhost.key")
	if err != nil {
		t.Fatal(err)
	}
	ts.TLS = &tls.Config{Certificates: []tls.Certificate{cert}}
	ts.StartTLS()
	defer ts.Close()

	// Replace server in test kube.config with mock server
	for i, cluster := range kcfg.Clusters {
		if cluster.Name == "test" {
			kcfg.Clusters[i].Cluster.Server = ts.URL
		}
	}

	data, err := GetApiData(kcfg, "/")
	if err != nil {
		t.Fatal(err)
	}

	// Compare if func returns what mock server served
	if string(data) != testBody {
		t.Fatal("Response does not match")
	}
}
