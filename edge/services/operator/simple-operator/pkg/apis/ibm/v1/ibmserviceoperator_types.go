package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// IBMserviceOperatorSpec defines the desired state of IBMserviceOperator
type IBMserviceOperatorSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	Size int32 `json:"size"`
}

// IBMserviceOperatorStatus defines the observed state of IBMserviceOperator
type IBMserviceOperatorStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	Nodes []string `json:"nodes"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// IBMserviceOperator is the Schema for the ibmserviceoperators API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=ibmserviceoperators,scope=Namespaced
type IBMserviceOperator struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   IBMserviceOperatorSpec   `json:"spec,omitempty"`
	Status IBMserviceOperatorStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// IBMserviceOperatorList contains a list of IBMserviceOperator
type IBMserviceOperatorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []IBMserviceOperator `json:"items"`
}

func init() {
	SchemeBuilder.Register(&IBMserviceOperator{}, &IBMserviceOperatorList{})
}
