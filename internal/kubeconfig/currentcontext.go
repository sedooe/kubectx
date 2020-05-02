package kubeconfig

// GetCurrentContext returns "current-context" value in given
// kubeconfig object Node, or returns "" if not found.
func (k *Kubeconfig) GetCurrentContext() string {
	kc := k.Current
	v := valueOf(kc.rootNode, "current-context")
	if v == nil {
		return ""
	}
	return v.Value
}

func (k *Kubeconfig) UnsetCurrentContext() error {
	kc := k.Current
	curCtxValNode := valueOf(kc.rootNode, "current-context")
	curCtxValNode.Value = ""
	return nil
}
