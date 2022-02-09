package v1alpha1

// +kubebuilder:validation:Enum=Pending;Ready;Rotating;Terminating;Error
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
	VaultBackend      Backend = "vault"
)
