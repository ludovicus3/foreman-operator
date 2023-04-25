/*
Copyright 2023.

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
	certmanagerv1 "github.com/jetstack/cert-manager/pkg/apis/certmanager/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ForemanSpec defines the desired state of Foreman
type ForemanSpec struct {
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Foreman Certificate Configuration"
	Certs ForemanCerts `json:"certs,omitempty"`
}

// ForemanCerts defines the configuration for generating internal certs
type ForemanCerts struct {
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Certificate Subject"
	Subject certmanagerv1.X509Subject `json:"subject:omitempty"`

	// +kubebuilder:default:="876000h"
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Certificate Duration"
	Duration metav1.Duration `json:"duration,omitempty"`

	// +kubebuilder:default:="175200h"
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="CA Certificate Duration"
	CADuration metav1.Duration `json:"caDuration,omitempty"`
}

// ForemanStatus defines the observed state of Foreman
type ForemanStatus struct {
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Foreman is the Schema for the foremen API
type Foreman struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ForemanSpec   `json:"spec,omitempty"`
	Status ForemanStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ForemanList contains a list of Foreman
type ForemanList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Foreman `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Foreman{}, &ForemanList{})
}
