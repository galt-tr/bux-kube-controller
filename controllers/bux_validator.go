package controllers

import "github.com/go-logr/logr"

func (r *BuxReconciler) Validate(log logr.Logger) (bool, error) {
	return true, nil
}
