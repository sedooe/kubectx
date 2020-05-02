package kubeconfig

import (
	"github.com/ahmetb/kubectx/internal/cmdutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

var (
	DefaultLoader Loader = new(StandardKubeconfigLoader)
)

type StandardKubeconfigLoader struct{}

type kubeconfigFile struct{ *os.File }

func (*StandardKubeconfigLoader) Load(cfgPath string) (ReadWriteResetCloser, error) {
	f, err := os.OpenFile(cfgPath, os.O_RDWR, 0)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errors.Wrap(err, "kubeconfig file not found")
		}
		return nil, errors.Wrap(err, "failed to open file")
	}
	return &kubeconfigFile{f}, nil
}

func (kf *kubeconfigFile) Reset() error {
	if err := kf.Truncate(0); err != nil {
		return errors.Wrap(err, "failed to truncate file")
	}
	_, err := kf.Seek(0, 0)
	return errors.Wrap(err, "failed to seek in file")
}

func kubeconfigPaths() ([]string, error) {
	// KUBECONFIG env var
	if v := os.Getenv("KUBECONFIG"); v != "" {
		return filepath.SplitList(v), nil
	}

	// default path
	home := cmdutil.HomeDir()
	if home == "" {
		return nil, errors.New("HOME or USERPROFILE environment variable not set")
	}
	return []string{filepath.Join(home, ".kube", "config")}, nil
}

// IsNotFoundErr determines if the underlying error is os.IsNotExist. Right now
// errors from github.com/pkg/errors doesn't work with os.IsNotExist.
func IsNotFoundErr(err error) bool {
	for e := err; e != nil; e = errors.Unwrap(e) {
		if os.IsNotExist(e) {
			return true
		}
	}
	return false
}
