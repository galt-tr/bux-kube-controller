package controllers

import (
	serverv1alpha1 "github.com/BuxOrg/bux-kube-controller/api/v1alpha1"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// ReconcileMongoPVC is the postgres PVC
func (r *BuxReconciler) ReconcileConsoleMongoPVC(log logr.Logger) (bool, error) {
	bux := serverv1alpha1.Bux{}
	if err := r.Get(r.Context, r.NamespacedName, &bux); err != nil {
		return false, err
	}
	pvc := corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "bux-console-mongo",
			Namespace: r.NamespacedName.Namespace,
			Labels:    r.getAppLabels(),
		},
	}

	_, _ = controllerutil.CreateOrUpdate(r.Context, r.Client, &pvc, func() error {
		return r.updateMongoPVC(&pvc, &bux)
	})
	// for now ignore errors since there are immutable fields
	/*if err != nil && !k8serrors.IsForbidden(err) {
		return false, err
	}*/
	return true, nil
}

func (r *BuxReconciler) updateMongoPVC(pvc *corev1.PersistentVolumeClaim, bux *serverv1alpha1.Bux) error {
	err := controllerutil.SetControllerReference(bux, pvc, r.Scheme)
	if err != nil {
		return err
	}
	pvc.Spec = *defaultMongoPVCSpec()
	return nil
}

func defaultMongoPVCSpec() *corev1.PersistentVolumeClaimSpec {
	return &corev1.PersistentVolumeClaimSpec{
		AccessModes: []corev1.PersistentVolumeAccessMode{
			corev1.ReadWriteOnce,
		},
		Resources: corev1.ResourceRequirements{
			Requests: corev1.ResourceList{
				"storage": resource.MustParse("1Gi"),
			},
		},
	}
}
