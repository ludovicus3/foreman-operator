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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ForemanSpec defines the desired state of Foreman
type ForemanSpec struct {
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
