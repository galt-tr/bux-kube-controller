package controllers

import (
	"fmt"

	serverv1alpha1 "github.com/BuxOrg/bux-kube-controller/api/v1alpha1"
	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// ReconcileConsoleMongoDeployment is the deployment
func (r *BuxReconciler) ReconcileConsoleMongoDeployment(_ logr.Logger) (bool, error) {
	bux := serverv1alpha1.Bux{}
	if err := r.Get(r.Context, r.NamespacedName, &bux); err != nil {
		return false, err
	}
	if !bux.Spec.Console {
		return false, nil
	}
	dep := appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "bux-console-mongo",
			Namespace: r.NamespacedName.Namespace,
			Labels:    r.getAppLabels(),
		},
	}
	_, err := controllerutil.CreateOrUpdate(r.Context, r.Client, &dep, func() error {
		return r.updateConsoleMongoDeployment(&dep, &bux)
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *BuxReconciler) updateConsoleMongoDeployment(dep *appsv1.Deployment, bux *serverv1alpha1.Bux) error {
	err := controllerutil.SetControllerReference(bux, dep, r.Scheme)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("https://%s-console.%s", bux.Namespace, bux.Spec.Domain)
	dep.Spec = *defaultConsoleMongoDeploymentSpec(url)
	return nil
}

func defaultConsoleMongoDeploymentSpec(_ string) *appsv1.DeploymentSpec {
	podLabels := map[string]string{
		"app":        "bux-console-mongo",
		"deployment": "bux-console-mongo",
	}
	image := "docker.io/mongo:latest"
	return &appsv1.DeploymentSpec{
		Replicas: pointer.Int32Ptr(1),
		Selector: metav1.SetAsLabelSelector(podLabels),
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				CreationTimestamp: metav1.Time{},
				Labels:            podLabels,
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Args: []string{
							"--storageEngine=wiredTiger",
						},
						Image:                    image,
						ImagePullPolicy:          corev1.PullAlways,
						Name:                     "bux-console-mongo",
						TerminationMessagePolicy: corev1.TerminationMessageFallbackToLogsOnError,
						VolumeMounts: []corev1.VolumeMount{
							{
								Name:      "data",
								MountPath: "/data/db",
							},
						},
						Ports: []corev1.ContainerPort{
							{
								ContainerPort: 27017,
								Protocol:      corev1.ProtocolTCP,
							},
						},
					},
				},
				Volumes: []corev1.Volume{
					{
						Name: "data",
						VolumeSource: corev1.VolumeSource{
							PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
								ClaimName: "bux-console-mongo",
							},
						},
					},
				},
			},
		},
	}
}
