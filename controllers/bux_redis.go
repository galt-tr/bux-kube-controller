package controllers

import (
	serverv1alpha1 "github.com/BuxOrg/bux-kube-controller/api/v1alpha1"
	"github.com/go-logr/logr"
	redisv1beta1 "github.com/murray-distributed-technologies/redis-operator/api/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// ReconcileRedis is for redis
func (r *BuxReconciler) ReconcileRedis(log logr.Logger) (bool, error) {
	bux := serverv1alpha1.Bux{}
	if err := r.Get(r.Context, r.NamespacedName, &bux); err != nil {
		return false, err
	}
	redis := redisv1beta1.Redis{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "redis-standalone",
			Namespace: r.NamespacedName.Namespace,
			Labels:    r.getAppLabels(),
		},
	}
	_, err := controllerutil.CreateOrUpdate(r.Context, r.Client, &redis, func() error {
		return r.updateRedis(&redis, &bux)
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *BuxReconciler) updateRedis(redis *redisv1beta1.Redis, _ *serverv1alpha1.Bux) error {
	redis.Spec = *defaultRedisSpec()
	return nil
}

func defaultRedisSpec() *redisv1beta1.RedisSpec {
	return &redisv1beta1.RedisSpec{
		KubernetesConfig: redisv1beta1.KubernetesConfig{
			Image:           "quay.io/opstree/redis:v6.2.5",
			ImagePullPolicy: corev1.PullAlways,
			Resources: &corev1.ResourceRequirements{
				Limits: corev1.ResourceList{
					"cpu":    resource.MustParse("101m"),
					"memory": resource.MustParse("128Mi"),
				},
				Requests: corev1.ResourceList{
					"cpu":    resource.MustParse("101m"),
					"memory": resource.MustParse("128Mi"),
				},
			},
		},
		RedisExporter: &redisv1beta1.RedisExporter{
			Image:   "quay.io/opstree/redis-exporter:1.0",
			Enabled: false,
		},
		Storage: &redisv1beta1.Storage{
			VolumeClaimTemplate: corev1.PersistentVolumeClaim{
				Spec: corev1.PersistentVolumeClaimSpec{
					AccessModes: []corev1.PersistentVolumeAccessMode{
						corev1.ReadWriteOnce,
					},
					Resources: corev1.ResourceRequirements{
						Requests: corev1.ResourceList{
							"storage": resource.MustParse("1Gi"),
						},
					},
				},
			},
		},
	}
}
