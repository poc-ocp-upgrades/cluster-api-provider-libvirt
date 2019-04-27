package controller

import (
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

var AddToManagerFuncs []func(manager.Manager) error

func AddToManager(m manager.Manager) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	for _, f := range AddToManagerFuncs {
		if err := f(m); err != nil {
			return err
		}
	}
	return nil
}
