package controllers

import (
	"encoding/json"
	"fmt"
	"time"

	serverv1alpha1 "github.com/BuxOrg/bux-kube-controller/api/v1alpha1"
	"github.com/BuxOrg/bux-server/config"
	"github.com/BuxOrg/bux/cachestore"
	"github.com/BuxOrg/bux/datastore"
	"github.com/BuxOrg/bux/taskmanager"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *BuxReconciler) ReconcileConfig(log logr.Logger) (bool, error) {
	bux := serverv1alpha1.Bux{}
	if err := r.Get(r.Context, r.NamespacedName, &bux); err != nil {
		return false, err
	}
	cm := corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "bux-config",
			Namespace: r.NamespacedName.Namespace,
		},
	}
	_, err := controllerutil.CreateOrUpdate(r.Context, r.Client, &cm, func() error {
		return r.updateBuxConfigMap(&cm, &bux)
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *BuxReconciler) updateBuxConfigMap(configMap *corev1.ConfigMap, bux *serverv1alpha1.Bux) error {
	err := controllerutil.SetControllerReference(bux, configMap, r.Scheme)
	if err != nil {
		return err
	}
	config := defaultBuxConfig()
	if bux.Spec.Configuration != nil && bux.Spec.Configuration.AdminXpub != "" {
		config.Authentication.AdminKey = bux.Spec.Configuration.AdminXpub
	}
	if bux.Spec.Domain != "" {
		config.Paymail.Domains[0] = fmt.Sprintf("%s.%s", bux.Namespace, bux.Spec.Domain)
	}
	data, err := json.Marshal(config)
	if err != nil {
		return err
	}
	configMap.Data = map[string]string{
		"development.json": string(data),
	}
	return nil
}

func defaultBuxConfig() *config.AppConfig {
	return &config.AppConfig{
		Debug:          true,
		DebugProfiling: false,
		DisableITC:     false,
		Environment:    "development",
		GDPRCompliance: false,
		Authentication: &config.AuthenticationConfig{
			AdminKey:        "12345",
			RequireSigning:  false,
			Scheme:          "xpub",
			SigningDisabled: false,
		},
		Cachestore: &config.CachestoreConfig{
			Engine: cachestore.Redis,
		},
		Datastore: &config.DatastoreConfig{
			AutoMigrate: true,
			Engine:      datastore.PostgreSQL,
			Debug:       true,
			TablePrefix: "bux",
		},
		GraphQL: &config.GraphqlConfig{
			Enabled:    true,
			ServerPath: "/graphql",
		},
		Mongo: &datastore.MongoDBConfig{},
		NewRelic: &config.NewRelicConfig{
			DomainName: "domain.com",
			Enabled:    false,
			LicenseKey: "BOGUS-LICENSE-KEY-1234567890987654321234",
		},
		Paymail: &config.PaymailConfig{
			Enabled:            true,
			DefaultFromPaymail: "from@domain.com",
			DefaultNote:        "bux address resolution",
			Domains: []string{
				"domain.com",
			},
			SenderValidationEnabled: false,
		},
		Redis: &config.RedisConfig{
			DependencyMode:        true,
			MaxActiveConnections:  0,
			MaxConnectionLifetime: time.Second * 10,
			MaxIdleConnections:    10,
			MaxIdleTimeout:        time.Second * 10,
			URL:                   "redis://redis-standalone:6379",
			UseTLS:                false,
		},
		Ristretto: &config.RistrettoConfig{},
		Server: &config.ServerConfig{
			IdleTimeout:  60 * time.Second,
			Port:         "3003",
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
		},
		SQL: &datastore.SQLConfig{
			Host:                      "bux-datastore",
			Name:                      "bux",
			Password:                  "postgres",
			Port:                      "5432",
			Replica:                   false,
			SkipInitializeWithVersion: true,
			TimeZone:                  "UTC",
			TxTimeout:                 10 * time.Second,
			User:                      "bux",
		},
		TaskManager: &config.TaskManagerConfig{
			Engine:    taskmanager.TaskQ,
			Factory:   taskmanager.FactoryRedis,
			QueueName: "development_queue",
		},
	}
}
