package controllers

import (
	serverv1alpha1 "github.com/BuxOrg/bux-kube-controller/api/v1alpha1"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// ReconcileConsoleService is the service
func (r *BuxReconciler) ReconcileConsoleMongoService(log logr.Logger) (bool, error) {
	bux := serverv1alpha1.Bux{}
	if err := r.Get(r.Context, r.NamespacedName, &bux); err != nil {
		return false, err
	}
	svc := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "bux-console-mongodb",
			Namespace: r.NamespacedName.Namespace,
		},
	}
	_, err := controllerutil.CreateOrUpdate(r.Context, r.Client, &svc, func() error {
		return r.updateConsoleMongodbService(&svc, &bux)
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *BuxReconciler) updateConsoleMongodbService(svc *corev1.Service, bux *serverv1alpha1.Bux) error {
	err := controllerutil.SetControllerReference(bux, svc, r.Scheme)
	if err != nil {
		return err
	}
	svc.Spec = *defaultConsoleMongodbServiceSpec()
	return nil
}

func defaultConsoleMongodbServiceSpec() *corev1.ServiceSpec {
	labels := map[string]string{
		"app": "bux-console-mongodb",
	}
	return &corev1.ServiceSpec{
		Selector: labels,
		Type:     corev1.ServiceTypeClusterIP,
		Ports: []corev1.ServicePort{
			{
				Name:       "27017",
				Port:       int32(27017),
				TargetPort: intstr.FromInt(27017),
			},
		},
	}
}
