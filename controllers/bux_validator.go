package controllers

import (
	"errors"
	"fmt"

	serverv1alpha1 "github.com/BuxOrg/bux-kube-controller/api/v1alpha1"
	"github.com/go-logr/logr"
)

// Validate will run validations
func (r *BuxReconciler) Validate(_ logr.Logger) (bool, error) {
	bux := serverv1alpha1.Bux{}
	if err := r.Get(r.Context, r.NamespacedName, &bux); err != nil {
		return false, err
	}
	if bux.Spec.Configuration.Datastore == "" {
		return false, errors.New("missing datastore configuration")
	}
	if err := validateDatastore(bux.Spec.Configuration.Datastore); err != nil {
		return false, err
	}
	return true, nil
}

func validateDatastore(datastore string) error {
	switch datastore {
	case "postgresql":
		return nil
	case "mongodb":
		return nil
	default:
		return fmt.Errorf("unsupported datastore %s", datastore)
	}
}
