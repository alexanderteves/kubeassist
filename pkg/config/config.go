package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
)

// ConnectionInfo bundles information from KubeConfig necessary to send requests
// to a Kubernetes API.
type ConnectionInfo struct {
	Token     string
	Server    string
	CA        string
	Namespace string
}

// KubeConfig represents a parsed kubeconfig file from the file system.
type Kubeconfig struct {
	ApiVersion string `yaml:"apiVersion"`
	Clusters   []struct {
		Cluster struct {
			CA     string `yaml:"certificate-authority"`
			Server string `yaml:"server"`
		} `yaml:"cluster"`
		Name string `yaml:"name"`
	} `yaml:"clusters"`
	Contexts []struct {
		Context struct {
			Cluster   string `yaml:"cluster"`
			Namespace string `yaml:"namespace"`
			User      string `yaml:"user"`
		} `yaml:"context"`
		Name string `yaml:"name"`
	} `yaml:"contexts"`
	CurrentContext string            `yaml:"current-context"`
	Kind           string            `yaml:"kind"`
	Preferences    map[string]string `yaml:"preferences"`
	Users          []struct {
		Name string `yaml:"name"`
		User struct {
			Token      string `yaml:"token"`
			ClientCert string `yaml:"client-certificate"`
			ClientKey  string `yaml:"client-key"`
		} `yaml:"user"`
	} `yaml:"users"`
}

// Load reads and unmarshals a new KubeConfig from path.
func Load(path string) (kcfg Kubeconfig, err error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return kcfg, fmt.Errorf("Reading file %v: %v", path, err)
	}
	err = yaml.Unmarshal(data, &kcfg)
	if err != nil {
		return kcfg, fmt.Errorf("Decoding %v: %v", path, err)
	}
	return kcfg, nil
}

// SetContext sets the KubeConfig's current context to newContext. Returns an
// error if the desired context does not exists.
func (k *Kubeconfig) SetContext(newContext string) error {
	var contextExists = false
	for _, context := range k.Contexts {
		if context.Name == newContext {
			contextExists = true
			break
		}
	}
	if !contextExists {
		return fmt.Errorf("Context %v does not exist", newContext)
	} else {
		k.CurrentContext = newContext
		return nil
	}
}

// SetNamespace sets the KubeConfig's namespace in its current context to
// newNamespace.
func (k *Kubeconfig) SetNamespace(newNamespace string) error {
	for index, context := range k.Contexts {
		if context.Name == k.CurrentContext {
			k.Contexts[index].Context.Namespace = newNamespace
			return nil
		}
	}
	return fmt.Errorf("Cannot set namespace for non-existent context %v", k.CurrentContext)
}

// GetConnectinfo is a convenient method to extract HTTP connection related
// information from KubeConfig.
func (k *Kubeconfig) GetConnectionInfo() (connInfo ConnectionInfo, err error) {
	if k.CurrentContext == "" {
		return connInfo, fmt.Errorf("No current context set")
	}

	var currentCluster, currentNamespace, currentUser, currentServer, currentToken, currentCa string
	for _, context := range k.Contexts {
		if context.Name == k.CurrentContext {
			currentCluster = context.Context.Cluster
			currentNamespace = context.Context.Namespace
			currentUser = context.Context.User
			break
		}
	}

	if currentNamespace == "" {
		currentNamespace = "default"
	}

	if currentCluster == "" || currentUser == "" {
		return connInfo, fmt.Errorf("Incomplete context section")
	}

	// TODO: Use cert and key if present
	for _, user := range k.Users {
		if user.Name == currentUser {
			currentToken = user.User.Token
			break
		}
	}

	if currentToken == "" {
		return connInfo, fmt.Errorf("No token found")
	}

	for _, cluster := range k.Clusters {
		if cluster.Name == currentCluster {
			currentServer = cluster.Cluster.Server
			currentCa = cluster.Cluster.CA
			break
		}
	}
	if currentServer == "" {
		return connInfo, fmt.Errorf("No server found")
	}

	connInfo.Token = currentToken
	connInfo.Server = currentServer
	connInfo.CA = currentCa
	connInfo.Namespace = currentNamespace
	return connInfo, nil
}

// Dump marshals and writes the KubeConfig to path.
func (k *Kubeconfig) Dump(path string) error {
	data, err := yaml.Marshal(k)
	if err != nil {
		return fmt.Errorf("Encoding: %v", err)
	}
	err = os.WriteFile(path, data, 0600)
	if err != nil {
		return fmt.Errorf("Writing file: %v", err)
	}
	return nil
}
