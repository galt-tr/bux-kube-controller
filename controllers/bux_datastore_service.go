package controllers

import (
	serverv1alpha1 "github.com/BuxOrg/bux-kube-controller/api/v1alpha1"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// ReconcileDatastoreService is the datastore service
func (r *BuxReconciler) ReconcileDatastoreService(_ logr.Logger) (bool, error) {
	bux := serverv1alpha1.Bux{}
	if err := r.Get(r.Context, r.NamespacedName, &bux); err != nil {
		return false, err
	}
	svc := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "bux-datastore",
			Namespace: r.NamespacedName.Namespace,
			Labels:    r.getAppLabels(),
		},
	}
	_, err := controllerutil.CreateOrUpdate(r.Context, r.Client, &svc, func() error {
		return r.updateDatastoreService(&svc, &bux)
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *BuxReconciler) updateDatastoreService(svc *corev1.Service, bux *serverv1alpha1.Bux) error {
	err := controllerutil.SetControllerReference(bux, svc, r.Scheme)
	if err != nil {
		return err
	}
	svc.Spec = *defaultDatastoreServiceSpec()
	return nil
}

func defaultDatastoreServiceSpec() *corev1.ServiceSpec {
	labels := map[string]string{
		"app":        "bux",
		"deployment": "bux-postgresql",
	}
	return &corev1.ServiceSpec{
		Selector: labels,
		Type:     corev1.ServiceTypeClusterIP,
		Ports: []corev1.ServicePort{
			{
				Name:       "5432",
				Port:       int32(5432),
				TargetPort: intstr.FromInt(5432),
			},
		},
	}
}
