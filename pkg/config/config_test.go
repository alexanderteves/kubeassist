package config

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestKubeconfig(t *testing.T) {
	kcfg, err := Load("../../test/kube.config")
	if err != nil {
		t.Fatal(err)
	}

	err = kcfg.SetNamespace("kube-system")
	if err != nil {
		t.Fatal(err)
	}

	err = kcfg.SetContext("missing_context")
	if err == nil {
		t.Fatal("Setting missing context did not fail")
	}

	err = kcfg.SetContext("test")
	if err != nil {
		t.Fatal(err)
	}

	connInfo, err := kcfg.GetConnectionInfo()
	if err != nil {
		t.Fatal(err)
	}

	switch {
	case connInfo.Token != "ahR5eich1Oulahpo7aht2ahci":
		t.Fatal("Token does not match test data")
	case connInfo.Server != "https://k8s.example.com:6443":
		t.Fatal("Server does not match test data")
	case connInfo.CA != "../../test/tls/ca.crt":
		t.Fatal("CA does not match test data")
	case connInfo.Namespace != "kube-system":
		t.Fatal("Namespace does not match test data")
	}

	tmpFile, err := ioutil.TempFile("", "kubeconfig_")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	err = kcfg.Dump(tmpFile.Name())
	if err != nil {
		t.Fatal(err)
	}
}
