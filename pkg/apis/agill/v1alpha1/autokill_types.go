package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// AutoKillSpec defines the desired state of AutoKill
type AutoKillSpec struct {
	DeleteNamespaceAfter int `json:"deleteNamespaceAfter,required"`
	DeleteAssociatedHelmReleases bool `json:"deleteAssociatedHelmReleases"`
	Disable bool `json:"disable"`
	TillerNamespace string `json:"tillerNamespace,omitempty"`

}

type AutoKillStatus struct {}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AutoKill is the Schema for the autokills API
// +k8s:openapi-gen=true
type AutoKill struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AutoKillSpec   `json:"spec,omitempty"`
	Status AutoKillStatus `json:"status,omitempty"`
}

type AutoKillList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AutoKill `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AutoKill{}, &AutoKillList{})
}
