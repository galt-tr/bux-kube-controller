/*
Copyright 2022 Dylan Murray.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ConditionReconciled is reconciled
const ConditionReconciled = "Reconciled"

// ReconciledReasonComplete is complete
const ReconciledReasonComplete = "Complete"

// ReconciledReasonError is an error
const ReconciledReasonError = "Error"

// ReconcileCompleteMessage is when the reconile is complete
const ReconcileCompleteMessage = "Reconcile complete"

// TODO: this should just be the bux config type, but its missing DeepCopy
// funcs or something like that idk:
// https://github.com/operator-framework/operator-sdk/issues/612

// BuxConfig is the BUX configuration
type BuxConfig struct {
	EnablePaymail  bool   `json:"enablePaymail"`
	AdminXpub      string `json:"adminXpub"`
	RequireSigning bool   `json:"requireSigning"`
	AutoMigrate    bool   `json:"autoMigrate"`
	//Agent          *AgentConfig `json:"agent"`
}

// AgentConfig is the bux agent configuration
type AgentConfig struct {
	URL string `json:"url"`
}

// BuxSpec defines the desired state of Bux
type BuxSpec struct {
	Configuration *BuxConfig `json:"configuration"`
	Domain        string     `json:"domain"`
	ClusterIssuer string     `json:"clusterIssuer"`
	Console       bool       `json:"console"`
}

// BuxStatus defines the observed state of Bux
type BuxStatus struct {
	Conditions []metav1.Condition `json:"conditions,omitempty"`
	Route      string             `json:"route,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Bux is the Schema for the buxes API
type Bux struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BuxSpec   `json:"spec,omitempty"`
	Status BuxStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// BuxList contains a list of Bux
type BuxList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Bux `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Bux{}, &BuxList{})
}
