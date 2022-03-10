package controllers

import (
	"fmt"

	serverv1alpha1 "github.com/BuxOrg/bux-kube-controller/api/v1alpha1"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *BuxReconciler) ReconcileIngress(log logr.Logger) (bool, error) {
	bux := serverv1alpha1.Bux{}
	if err := r.Get(r.Context, r.NamespacedName, &bux); err != nil {
		return false, err
	}
	// Skip if domain isn't set
	if bux.Spec.Domain == "" {
		return false, nil
	}
	ingress := networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "bux",
			Namespace: r.NamespacedName.Namespace,
		},
	}
	_, err := controllerutil.CreateOrUpdate(r.Context, r.Client, &ingress, func() error {
		return r.updateIngress(&ingress, &bux)
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *BuxReconciler) ReconcileService(log logr.Logger) (bool, error) {
	bux := serverv1alpha1.Bux{}
	if err := r.Get(r.Context, r.NamespacedName, &bux); err != nil {
		return false, err
	}
	svc := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "bux",
			Namespace: r.NamespacedName.Namespace,
		},
	}
	_, err := controllerutil.CreateOrUpdate(r.Context, r.Client, &svc, func() error {
		return r.updateService(&svc, &bux)
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *BuxReconciler) updateIngress(ingress *networkingv1.Ingress, bux *serverv1alpha1.Bux) error {
	err := controllerutil.SetControllerReference(bux, ingress, r.Scheme)
	if err != nil {
		return err
	}
	if bux.Spec.ClusterIssuer != "" {
		if ingress.Annotations == nil {
			ingress.Annotations = make(map[string]string)
		}
		ingress.Annotations["cert-manager.io/cluster-issuer"] = bux.Spec.ClusterIssuer
	}
	ingress.Spec = *defaultIngressSpec(bux)
	return nil
}

func (r *BuxReconciler) updateService(svc *corev1.Service, bux *serverv1alpha1.Bux) error {
	err := controllerutil.SetControllerReference(bux, svc, r.Scheme)
	if err != nil {
		return err
	}
	svc.Spec = *defaultServiceSpec()
	return nil
}

func defaultIngressSpec(bux *serverv1alpha1.Bux) *networkingv1.IngressSpec {
	pathType := networkingv1.PathTypeImplementationSpecific
	return &networkingv1.IngressSpec{
		TLS: []networkingv1.IngressTLS{
			networkingv1.IngressTLS{
				Hosts: []string{
					fmt.Sprintf("%s.%s", bux.Namespace, bux.Spec.Domain),
				},
				SecretName: "bux-tls",
			},
		},
		Rules: []networkingv1.IngressRule{
			networkingv1.IngressRule{
				Host: fmt.Sprintf("%s.%s", bux.Namespace, bux.Spec.Domain),
				IngressRuleValue: networkingv1.IngressRuleValue{
					HTTP: &networkingv1.HTTPIngressRuleValue{
						Paths: []networkingv1.HTTPIngressPath{
							networkingv1.HTTPIngressPath{
								PathType: &pathType,
								Backend: networkingv1.IngressBackend{
									Service: &networkingv1.IngressServiceBackend{
										Name: "bux",
										Port: networkingv1.ServiceBackendPort{
											Number: int32(3003),
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func defaultServiceSpec() *corev1.ServiceSpec {
	labels := map[string]string{
		"app": "bux",
	}
	return &corev1.ServiceSpec{
		Selector: labels,
		Type:     corev1.ServiceTypeClusterIP,
		Ports: []corev1.ServicePort{
			corev1.ServicePort{
				Name:       "3003",
				Port:       int32(3003),
				TargetPort: intstr.FromInt(3003),
			},
			corev1.ServicePort{
				Name:       "443",
				Port:       int32(443),
				TargetPort: intstr.FromInt(443),
			},
			corev1.ServicePort{
				Name:       "80",
				Port:       int32(80),
				TargetPort: intstr.FromInt(80),
			},
		},
	}
}
