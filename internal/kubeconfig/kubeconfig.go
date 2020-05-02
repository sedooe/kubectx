package kubeconfig

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"io"
	"os"
)

type ReadWriteResetCloser interface {
	io.ReadWriteCloser

	// Reset truncates the file and seeks to the beginning of the file.
	Reset() error
}

type Loader interface {
	Load(string) (ReadWriteResetCloser, error)
}

type Kubeconfig struct {
	loader Loader
	Current *SingleKubeconfig
	ConfigsMap map[string]*SingleKubeconfig
	kcPaths	[]string
}

type SingleKubeconfig struct {
	path	 string
	f        ReadWriteResetCloser
	rootNode *yaml.Node
}

func (k *Kubeconfig) WithLoader(l Loader) *Kubeconfig {
	k.loader = l
	return k
}

func (k *Kubeconfig) Close() error {
	if k.Current == nil {
		return nil
	}
	return k.Current.f.Close()
}

func (k *Kubeconfig) Parse() error {
	paths, err := kubeconfigPaths()
	if err != nil {
		return errors.Wrap(err, "cannot determine kubeconfig path")
	}

	k.kcPaths = paths
	configsMap := map[string]*SingleKubeconfig{}

	for _, path := range paths {
		f, err := k.loader.Load(path)
		if err != nil {
			return errors.Wrap(err, "failed to load")
		}

		kconf := &SingleKubeconfig{}
		kconf.path = path
		kconf.f = f

		var v yaml.Node
		if err := yaml.NewDecoder(f).Decode(&v); err != nil {
			return errors.Wrap(err, "failed to decode")
		}
		kconf.rootNode = v.Content[0]
		if kconf.rootNode.Kind != yaml.MappingNode {
			return errors.New("kubeconfig file is not a map document")
		}

		configsMap[path] = kconf
	}

	k.Current = configsMap[k.getDefaultFilename()]
	k.ConfigsMap = configsMap

	return nil
}

func (k *Kubeconfig) Bytes() ([]byte, error) {
	return yaml.Marshal(k.Current.rootNode)
}

func (k *Kubeconfig) Save() error {
	if err := k.Current.f.Reset(); err != nil {
		return errors.Wrap(err, "failed to reset file")
	}
	enc := yaml.NewEncoder(k.Current.f)
	enc.SetIndent(0)
	return enc.Encode(k.Current.rootNode)
}

func (k *Kubeconfig) getDefaultFilename() string {
	envVarFiles := k.kcPaths

	if len(envVarFiles) == 1 {
		return envVarFiles[0]
	}

	// if any of the envvar files already exists, return it
	for _, envVarFile := range envVarFiles {
		if _, err := os.Stat(envVarFile); err == nil {
			return envVarFile
		}
	}

	// otherwise, return the last one in the list
	return envVarFiles[len(envVarFiles)-1]
}
