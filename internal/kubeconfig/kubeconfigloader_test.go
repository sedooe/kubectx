package kubeconfig

import (
	"github.com/ahmetb/kubectx/internal/testutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func Test_kubeconfigPath(t *testing.T) {
	defer testutil.WithEnvVar("HOME", "/x/y/z")()

	expected := []string{filepath.FromSlash("/x/y/z/.kube/config")}
	got, err := kubeconfigPaths()
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("got=%q expected=%q", got, expected)
	}
}

func Test_kubeconfigPath_noEnvVars(t *testing.T) {
	defer testutil.WithEnvVar("XDG_CACHE_HOME", "")()
	defer testutil.WithEnvVar("HOME", "")()
	defer testutil.WithEnvVar("USERPROFILE", "")()

	_, err := kubeconfigPaths()
	if err == nil {
		t.Fatalf("expected error")
	}
}

func Test_kubeconfigPath_envOvveride(t *testing.T) {
	defer testutil.WithEnvVar("KUBECONFIG", "foo")()

	v, err := kubeconfigPaths()
	if err != nil {
		t.Fatal(err)
	}
	if expected := []string{"foo"}; !reflect.DeepEqual(v, expected) {
		t.Fatalf("expected=%q, got=%q", expected, v)
	}
}

func Test_kubeconfigPath_envOvverideDoesNotSupportPathSeparator(t *testing.T) {
	expected := []string{"file1", "file2"}
	path := strings.Join(expected, string(os.PathListSeparator))
	defer testutil.WithEnvVar("KUBECONFIG", path)()

	got, err := kubeconfigPaths()
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("got=%q expected=%q", got, expected)
	}
}

func TestStandardKubeconfigLoader_returnsNotFoundErr(t *testing.T) {
	defer testutil.WithEnvVar("KUBECONFIG", "foo")()
	kc := new(Kubeconfig).WithLoader(DefaultLoader)
	err := kc.Parse()
	if err == nil {
		t.Fatal("expected err")
	}
	if !IsNotFoundErr(err) {
		t.Fatalf("expected ENOENT error; got=%v", err)
	}
}
