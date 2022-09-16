package controllers

import (
	serverv1alpha1 "github.com/BuxOrg/bux-kube-controller/api/v1alpha1"
	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// ReconcileDeployment is the deployment
func (r *BuxReconciler) ReconcileDeployment(_ logr.Logger) (bool, error) {
	bux := serverv1alpha1.Bux{}
	if err := r.Get(r.Context, r.NamespacedName, &bux); err != nil {
		return false, err
	}
	dep := appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "bux",
			Namespace: r.NamespacedName.Namespace,
			Labels:    r.getAppLabels(),
		},
	}
	_, err := controllerutil.CreateOrUpdate(r.Context, r.Client, &dep, func() error {
		return r.updateDeployment(&dep, &bux)
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *BuxReconciler) updateDeployment(dep *appsv1.Deployment, bux *serverv1alpha1.Bux) error {
	err := controllerutil.SetControllerReference(bux, dep, r.Scheme)
	if err != nil {
		return err
	}
	dep.Spec = *defaultDeploymentSpec()
	return nil
}

func defaultDeploymentSpec() *appsv1.DeploymentSpec {
	podLabels := map[string]string{
		"app":        "bux",
		"deployment": "bux",
	}
	var envFrom []corev1.EnvFromSource
	envVars := []corev1.EnvVar{
		{
			Name:  "BUX_ENVIRONMENT",
			Value: "development",
		},
	}
	image := "docker.io/buxorg/bux-server:latest"
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
						EnvFrom:                  envFrom,
						Env:                      envVars,
						Image:                    image,
						ImagePullPolicy:          corev1.PullAlways,
						Name:                     "bux",
						TerminationMessagePolicy: corev1.TerminationMessageFallbackToLogsOnError,
						Ports: []corev1.ContainerPort{
							{
								ContainerPort: 3003,
								Protocol:      corev1.ProtocolTCP,
							},
							{
								ContainerPort: 443,
								Protocol:      corev1.ProtocolTCP,
							},
							{
								ContainerPort: 80,
								Protocol:      corev1.ProtocolTCP,
							},
						},
						VolumeMounts: []corev1.VolumeMount{
							{
								MountPath: "config/envs",
								Name:      "config",
							},
						},
					},
				},
				Volumes: []corev1.Volume{
					{
						Name: "config",
						VolumeSource: corev1.VolumeSource{
							ConfigMap: &corev1.ConfigMapVolumeSource{
								LocalObjectReference: corev1.LocalObjectReference{
									Name: "bux-config",
								},
							},
						},
					},
				},
			},
		},
	}
}
