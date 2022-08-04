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

// ReconcileConsole will reconcile console
func (r *BuxReconciler) ReconcileConsole(log logr.Logger) (bool, error) {
	bux := serverv1alpha1.Bux{}
	if err := r.Get(r.Context, r.NamespacedName, &bux); err != nil {
		return false, err
	}
	if !bux.Spec.Console {
		return false, nil
	}
	return ReconcileBatch(log,
		r.ReconcileConsoleDeployment,
		r.ReconcileConsoleService,
		r.ReconcileConsoleMongoDeployment,
		r.ReconcileConsoleMongoService,
		r.ReconcileConsoleMongoPVC,
		r.ReconcileConsoleIngress,
	)

}

// ReconcileConsoleDeployment is the deployment
func (r *BuxReconciler) ReconcileConsoleDeployment(_ logr.Logger) (bool, error) {
	bux := serverv1alpha1.Bux{}
	if err := r.Get(r.Context, r.NamespacedName, &bux); err != nil {
		return false, err
	}
	dep := appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "bux-console",
			Namespace: r.NamespacedName.Namespace,
			Labels:    r.getAppLabels(),
		},
	}
	_, err := controllerutil.CreateOrUpdate(r.Context, r.Client, &dep, func() error {
		return r.updateConsoleDeployment(&dep, &bux)
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *BuxReconciler) updateConsoleDeployment(dep *appsv1.Deployment, bux *serverv1alpha1.Bux) error {
	err := controllerutil.SetControllerReference(bux, dep, r.Scheme)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("https://%s-console.%s", bux.Namespace, bux.Spec.Domain)
	dep.Spec = *defaultConsoleDeploymentSpec(url)
	return nil
}

func defaultConsoleDeploymentSpec(url string) *appsv1.DeploymentSpec {
	podLabels := map[string]string{
		"app":        "bux-console",
		"deployment": "bux-console",
	}
	var envFrom []corev1.EnvFromSource
	envVars := []corev1.EnvVar{
		{
			Name:  "ROOT_URL",
			Value: url,
		},
		{
			Name:  "PORT",
			Value: "3000",
		},
		{
			Name:  "MONGO_URL",
			Value: "mondogb://bux-console-mongodb:27017/meteor",
		},
	}
	image := "docker.io/galtbv/bux-console:latest"
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
						Name:                     "bux-console",
						TerminationMessagePolicy: corev1.TerminationMessageFallbackToLogsOnError,
						Ports: []corev1.ContainerPort{
							{
								ContainerPort: 3000,
								Protocol:      corev1.ProtocolTCP,
							},
							{
								ContainerPort: 80,
								Protocol:      corev1.ProtocolTCP,
							},
							{
								ContainerPort: 443,
								Protocol:      corev1.ProtocolTCP,
							},
						},
					},
				},
			},
		},
	}
}
