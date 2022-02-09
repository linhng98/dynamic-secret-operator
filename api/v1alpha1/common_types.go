package v1alpha1

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// +kubebuilder:validation:Enum=Pending;Ready;Rotating;Terminating;error
type Phase string
const (
	// PhasePending means a custom secret resource has just been created and is not yet ready
	PhasePending Phase = "Pending"

	// PhaseReady means a custom secret resource is ready and up to date
	PhaseReady Phase = "Ready"

	// PhaseUpdating means a custom secret resource is in the process of rotating to a new desired state (spec)
	PhaseRotating Phase = "Rotating"

	// PhaseTerminating means a custom secret resource is in the process of being removed
	PhaseTerminating Phase = "Terminating"

	// PhaseError means an error occured with custom resource management
	PhaseError Phase = "Error"
)

// +kubebuilder:validation:Enum=kubernetes;vault
type Backend string

const (
  KubernetesBackend Backend = "kubernetes"
  VaultBackend Backend = "vault" 
)
