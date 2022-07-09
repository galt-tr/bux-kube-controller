package controllers

import (
	"encoding/json"
	"fmt"
	"time"

	serverv1alpha1 "github.com/BuxOrg/bux-kube-controller/api/v1alpha1"
	"github.com/BuxOrg/bux-server/config"
	"github.com/BuxOrg/bux/datastore"
	"github.com/BuxOrg/bux/taskmanager"
	"github.com/go-logr/logr"
	"github.com/mrz1836/go-cachestore"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// ReconcileConfig will reconcile configuration
func (r *BuxReconciler) ReconcileConfig(log logr.Logger) (bool, error) {
	bux := serverv1alpha1.Bux{}
	if err := r.Get(r.Context, r.NamespacedName, &bux); err != nil {
		return false, err
	}
	cm := corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "bux-config",
			Namespace: r.NamespacedName.Namespace,
			Labels:    r.getAppLabels(),
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

// updateBuxConfigMap will update the config
func (r *BuxReconciler) updateBuxConfigMap(configMap *corev1.ConfigMap, bux *serverv1alpha1.Bux) error {
	err := controllerutil.SetControllerReference(bux, configMap, r.Scheme)
	if err != nil {
		return err
	}
	configuration := defaultBuxConfig()
	if bux.Spec.Configuration != nil && bux.Spec.Configuration.AdminXpub != "" {
		configuration.Authentication.AdminKey = bux.Spec.Configuration.AdminXpub
	}
	if bux.Spec.Domain != "" {
		configuration.Paymail.Domains[0] = fmt.Sprintf("%s.%s", bux.Namespace, bux.Spec.Domain)
	}

	if bux.Spec.Configuration.Paymail != nil {
		configuration.Paymail.Enabled = bux.Spec.Configuration.Paymail.Enabled
		configuration.Paymail.DefaultNote = bux.Spec.Configuration.Paymail.DefaultNote
		configuration.Paymail.DefaultFromPaymail = bux.Spec.Configuration.Paymail.DefaultFromPaymail
		configuration.Paymail.DomainValidationEnabled = bux.Spec.Configuration.Paymail.DomainValidationEnabled
		configuration.Paymail.SenderValidationEnabled = bux.Spec.Configuration.Paymail.SenderValidationEnabled
	}

	var data []byte
	if data, err = json.Marshal(configuration); err != nil {
		return err
	}
	configMap.Data = map[string]string{
		"development.json": string(data),
	}
	return nil
}

// defaultBuxConfig is the default configuration
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
