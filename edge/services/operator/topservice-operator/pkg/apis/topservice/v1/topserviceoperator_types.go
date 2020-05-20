package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// TopServiceOperatorSpec defines the desired state of TopServiceOperator
type TopServiceOperatorSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	Size int32 `json:"size"`
}

// TopServiceOperatorStatus defines the observed state of TopServiceOperator
type TopServiceOperatorStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	Nodes []string `json:"nodes"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// TopServiceOperator is the Schema for the topserviceoperators API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=topserviceoperators,scope=Namespaced
type TopServiceOperator struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TopServiceOperatorSpec   `json:"spec,omitempty"`
	Status TopServiceOperatorStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// TopServiceOperatorList contains a list of TopServiceOperator
type TopServiceOperatorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TopServiceOperator `json:"items"`
}

func init() {
	SchemeBuilder.Register(&TopServiceOperator{}, &TopServiceOperatorList{})
}
