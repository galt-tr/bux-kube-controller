package controllers

import "github.com/go-logr/logr"

// Validate will run validations
func (r *BuxReconciler) Validate(log logr.Logger) (bool, error) {
	return true, nil
}
