/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apimeta "k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	serverv1alpha1 "github.com/BuxOrg/bux-kube-controller/api/v1alpha1"
	"github.com/go-logr/logr"
)

// BuxReconciler reconciles a Bux object
type BuxReconciler struct {
	client.Client
	Log            logr.Logger
	Scheme         *runtime.Scheme
	Context        context.Context
	NamespacedName types.NamespacedName
}

// +kubebuilder:rbac:groups=server.getbux.io,resources=buxes,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=server.getbux.io,resources=buxes/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=server.getbux.io,resources=buxes/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Bux object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *BuxReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	r.Log = log.FromContext(ctx)
	logger := r.Log.WithValues("bux", req.NamespacedName)
	result := ctrl.Result{}
	r.Context = ctx
	r.NamespacedName = req.NamespacedName
	bux := serverv1alpha1.Bux{}

	if err := r.Get(ctx, req.NamespacedName, &bux); err != nil {
		logger.Error(err, "unable to fetch Bux CR")
		return result, nil
	}

	_, err := ReconcileBatch(r.Log,
		r.Validate,
		r.ReconcileConfig,
		r.ReconcileDatastore,
		r.ReconcileRedis,
		r.ReconcileService,
		r.ReconcileIngress,
		r.ReconcileDeployment,
	)

	if err != nil {
		apimeta.SetStatusCondition(&bux.Status.Conditions,
			metav1.Condition{
				Type:    serverv1alpha1.ConditionReconciled,
				Status:  metav1.ConditionFalse,
				Reason:  serverv1alpha1.ReconciledReasonError,
				Message: err.Error(),
			},
		)
	} else {
		apimeta.SetStatusCondition(&bux.Status.Conditions,
			metav1.Condition{
				Type:    serverv1alpha1.ConditionReconciled,
				Status:  metav1.ConditionTrue,
				Reason:  serverv1alpha1.ReconciledReasonComplete,
				Message: serverv1alpha1.ReconcileCompleteMessage,
			},
		)
	}

	statusErr := r.Client.Status().Update(ctx, &bux)
	if err == nil {
		err = statusErr
	}

	return ctrl.Result{}, err
}

// SetupWithManager sets up the controller with the Manager.
func (r *BuxReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&serverv1alpha1.Bux{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Owns(&corev1.ConfigMap{}).
		Complete(r)
}

// ReconcileFunc is a reconcile function type
type ReconcileFunc func(logr.Logger) (bool, error)

// ReconcileBatch will reconcile the batch of functions
func ReconcileBatch(l logr.Logger, reconcileFunctions ...ReconcileFunc) (bool, error) {
	for _, f := range reconcileFunctions {
		if cont, err := f(l); !cont || err != nil {
			return cont, err
		}
	}
	return true, nil
}
