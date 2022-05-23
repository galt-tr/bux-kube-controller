package controllers

import (
	serverv1alpha1 "github.com/BuxOrg/bux-kube-controller/api/v1alpha1"
	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// ReconcileDatastore is the datastore
func (r *BuxReconciler) ReconcileDatastore(log logr.Logger) (bool, error) {
	return ReconcileBatch(log,
		r.ReconcilePostgresqlDeployment,
		r.ReconcilePostgresqlPVC,
		r.ReconcileDatastoreService,
	)
}

func (r *BuxReconciler) ReconcileDatastoreDeployment(log logr.Logger) (bool, error) {
	return r.ReconcilePostgresqlDeployment(log)
}

// ReconcilePostgresqlDeployment is the postgres deployment
func (r *BuxReconciler) ReconcilePostgresqlDeployment(log logr.Logger) (bool, error) {
	bux := serverv1alpha1.Bux{}
	if err := r.Get(r.Context, r.NamespacedName, &bux); err != nil {
		return false, err
	}
	dep := appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "bux-postgresql",
			Namespace: r.NamespacedName.Namespace,
			Labels:    r.getAppLabels(),
		},
	}
	_, err := controllerutil.CreateOrUpdate(r.Context, r.Client, &dep, func() error {
		return r.updatePostgresqlDeployment(&dep, &bux)
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

// ReconcilePostgresqlPVC is the postgres PVC
func (r *BuxReconciler) ReconcilePostgresqlPVC(log logr.Logger) (bool, error) {
	bux := serverv1alpha1.Bux{}
	if err := r.Get(r.Context, r.NamespacedName, &bux); err != nil {
		return false, err
	}
	pvc := corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "bux-postgresql",
			Namespace: r.NamespacedName.Namespace,
			Labels:    r.getAppLabels(),
		},
	}

	_, _ = controllerutil.CreateOrUpdate(r.Context, r.Client, &pvc, func() error {
		return r.updatePVC(&pvc, &bux)
	})
	// for now ignore errors since there are immutable fields
	/*if err != nil && !k8serrors.IsForbidden(err) {
		return false, err
	}*/
	return true, nil
}

func (r *BuxReconciler) updatePostgresqlDeployment(dep *appsv1.Deployment, bux *serverv1alpha1.Bux) error {
	err := controllerutil.SetControllerReference(bux, dep, r.Scheme)
	if err != nil {
		return err
	}
	dep.Spec = *defaultPostgresqlDeploymentSpec()
	return nil
}

func (r *BuxReconciler) updatePVC(pvc *corev1.PersistentVolumeClaim, bux *serverv1alpha1.Bux) error {
	err := controllerutil.SetControllerReference(bux, pvc, r.Scheme)
	if err != nil {
		return err
	}
	pvc.Spec = *defaultPVCSpec()
	return nil
}

func defaultPVCSpec() *corev1.PersistentVolumeClaimSpec {
	return &corev1.PersistentVolumeClaimSpec{
		AccessModes: []corev1.PersistentVolumeAccessMode{
			corev1.ReadWriteOnce,
		},
		Resources: corev1.ResourceRequirements{
			Requests: corev1.ResourceList{
				"storage": resource.MustParse("2Gi"),
			},
		},
	}
}

func defaultPostgresqlDeploymentSpec() *appsv1.DeploymentSpec {
	podLabels := map[string]string{
		"app":        "bux",
		"deployment": "bux-postgresql",
	}
	var envFrom []corev1.EnvFromSource
	envVars := []corev1.EnvVar{
		{
			Name:  "POSTGRESQL_USER",
			Value: "bux",
		},
		{
			Name:  "POSTGRESQL_PASSWORD",
			Value: "postgres",
		},
		{
			Name:  "POSTGRESQL_DATABASE",
			Value: "bux",
		},
	}
	image := "docker.io/galtbv/postgresql-12"
	return &appsv1.DeploymentSpec{
		Replicas: pointer.Int32Ptr(1),
		Selector: metav1.SetAsLabelSelector(podLabels),
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				CreationTimestamp: metav1.Time{},
				Labels:            podLabels,
			},
			Spec: corev1.PodSpec{
				InitContainers: []corev1.Container{
					{
						Name:  "pgsql-data-permission-fix",
						Image: "busybox",
						Command: []string{
							"/bin/chmod",
							"-R",
							"777",
							"/data",
						},
						VolumeMounts: []corev1.VolumeMount{
							{
								Name:      "psql-data",
								MountPath: "/data",
							},
						},
					},
				},
				Containers: []corev1.Container{
					{
						EnvFrom:                  envFrom,
						Env:                      envVars,
						Image:                    image,
						Name:                     "postgresql",
						TerminationMessagePolicy: corev1.TerminationMessageFallbackToLogsOnError,
						Ports: []corev1.ContainerPort{
							{
								ContainerPort: 5432,
								Protocol:      corev1.ProtocolTCP,
							},
						},
						VolumeMounts: []corev1.VolumeMount{
							{
								MountPath: "/var/lib/pgsql/data",
								Name:      "psql-data",
							},
						},
					},
				},
				Volumes: []corev1.Volume{
					{
						Name: "psql-data",
						VolumeSource: corev1.VolumeSource{
							PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
								ClaimName: "bux-postgresql",
							},
						},
					},
				},
			},
		},
	}
}
