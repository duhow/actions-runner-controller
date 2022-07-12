/*
Copyright 2020 The actions-runner-controller authors.

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
	"errors"
	"fmt"

	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/util/validation/field"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// RunnerSpec defines the desired state of Runner
type RunnerSpec struct {
	RunnerConfig  `json:",inline"`
	RunnerPodSpec `json:",inline"`
}

type RunnerConfig struct {
	// Name of the GitHub Enterprise to register this runner.
	// Don't use if using Organization or Repository.
	// +optional
	// +kubebuilder:validation:Pattern=`^[^/]+$`
	Enterprise string `json:"enterprise,omitempty"`

	// Name of the Organization to register this runner.
	// Don't use if using Enterprise or Repository.
	// +optional
	// +kubebuilder:validation:Pattern=`^[^/]+$`
	Organization string `json:"organization,omitempty"`

	// Name of the Repository to register this runner.
	// The runner will be able to execute jobs from this repository only.
	// Syntax of this value is `USER/REPO`.
	// +optional
	// +kubebuilder:validation:Pattern=`^[^/]+/[^/]+$`
	Repository string `json:"repository,omitempty"`

	// List of strings with labels to apply to this runner,
	// in order to distinguish from other kind of self-hosted runners.
	// +optional
	Labels []string `json:"labels,omitempty"`

	// Name of the Runner Group to add.
	// Runner groups can be used to limit which repositories are able to use the 
	// GitHub Runner at an organization level. Runner groups have to be created 
	// in GitHub first before they can be referenced.
	// +optional
	Group string `json:"group,omitempty"`

	// Ephemeral flag restarts the runner after running jobs, to ensure a 
	// clean and reproducible build environment.
	// If set to false, the runner is persistent across jobs
	// and doesn't automatically restart.
	// Defaults to true.
	// +optional
	Ephemeral *bool `json:"ephemeral,omitempty"`

	// Docker image name to use.
	// +optional
	Image string `json:"image"`

	// You can customise this setting allowing you to change the 
	// default working directory location. Defaults to /runner/_work
	// +optional
	WorkDir string `json:"workDir,omitempty"`

	// If set to true, no Docker sidecar container is deployed in the runner pod,
	// but docker can be used within the runner container instead.
	// The image summerwind/actions-runner-dind is used by default.
	// If set to false, Docker support is provided by a sidecar container deployed in the runner pod.
	// Defaults to false.
	// +optional
	DockerdWithinRunnerContainer *bool `json:"dockerdWithinRunnerContainer,omitempty"`

	// If set to true, a privileged docker sidecar container is included in the runner pod.
	// If set to false, a docker sidecar container is not included in the runner pod and you can't use docker.
	// Defaults to true.
	// +optional
	DockerEnabled *bool `json:"dockerEnabled,omitempty"`

	// If your network card MTU is smaller than Docker's default 1500, you might encounter Docker networking issues.
	// To fix these issues, you should setup Docker MTU smaller than or equal to that on the outgoing network card.
	// More information: https://mlohr.com/docker-mtu/
	// +optional
	DockerMTU *int64 `json:"dockerMTU,omitempty"`

	// Since Docker Hub applies rate-limit configuration for free plans, to avoid disruptions in your CI/CD pipelines,
	// you might want to setup an external or on-premises Docker registry mirror.
	// You can provide the URL of the registry mirror.
	// More information: https://docs.docker.com/docker-hub/download-rate-limit/
	// +optional
	DockerRegistryMirror *string `json:"dockerRegistryMirror,omitempty"`

	// Total amount of local storage resources required for runner volume mount.
	// The default limit is undefined.
	// +optional
	VolumeSizeLimit *resource.Quantity `json:"volumeSizeLimit,omitempty"`

	// Optional storage medium type of runner volume mount.
	// More info: https://kubernetes.io/docs/concepts/storage/volumes/#emptydir
	// Allowed value is: Memory.
	// +optional
	// +kubebuilder:validation:Enum=Memory
	VolumeStorageMedium *string `json:"volumeStorageMedium,omitempty"`

	// By default, jobs that use a container will run in Docker, which requires privileged mode.
	// If you set the container mode, you can instead invoke these jobs using a Kubernetes implementation
	// without requiring to execute in privileged mode.
	// Allowed value is: kubernetes.
	// Requires defining WorkVolumeClaimTemplate.
	// More information: https://github.com/actions-runner-controller/actions-runner-controller/pull/1546
	// +optional
	// +kubebuilder:validation:Enum=kubernetes
	ContainerMode string `json:"containerMode,omitempty"`

	GitHubAPICredentialsFrom *GitHubAPICredentialsFrom `json:"githubAPICredentialsFrom,omitempty"`
}

type GitHubAPICredentialsFrom struct {
	SecretRef SecretReference `json:"secretRef,omitempty"`
}

type SecretReference struct {
	Name string `json:"name"`
}

// RunnerPodSpec defines the desired pod spec fields of the runner pod
type RunnerPodSpec struct {
	// +optional
	DockerdContainerResources corev1.ResourceRequirements `json:"dockerdContainerResources,omitempty"`

	// +optional
	DockerVolumeMounts []corev1.VolumeMount `json:"dockerVolumeMounts,omitempty"`

	// +optional
	DockerEnv []corev1.EnvVar `json:"dockerEnv,omitempty"`

	// +optional
	Containers []corev1.Container `json:"containers,omitempty"`

	// Image pull policy. One of Always, Never, IfNotPresent. Defaults to Always
	// +optional
	ImagePullPolicy corev1.PullPolicy `json:"imagePullPolicy,omitempty"`

	// +optional
	Env []corev1.EnvVar `json:"env,omitempty"`

	// +optional
	EnvFrom []corev1.EnvFromSource `json:"envFrom,omitempty"`

	// +optional
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`

	// +optional
	VolumeMounts []corev1.VolumeMount `json:"volumeMounts,omitempty"`

	// +optional
	Volumes []corev1.Volume `json:"volumes,omitempty"`

	// EnableServiceLinks indicates whether information about services should be
	// injected into pod's environment variables, matching the syntax of Docker
	// links. Defaults to true.
	// +optional
	EnableServiceLinks *bool `json:"enableServiceLinks,omitempty"`

	// +optional
	InitContainers []corev1.Container `json:"initContainers,omitempty"`

	// +optional
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`

	// ServiceAccountName is the name of the ServiceAccount to use to run this pod.
	// More info: https://kubernetes.io/docs/tasks/configure-pod-container/configure-service-account/
	// +optional
	ServiceAccountName string `json:"serviceAccountName,omitempty"`

	// AutomountServiceAccountToken indicates whether a service account token
	// should be automatically mounted.
	// +optional
	AutomountServiceAccountToken *bool `json:"automountServiceAccountToken,omitempty"`

	// +optional
	SidecarContainers []corev1.Container `json:"sidecarContainers,omitempty"`

	// +optional
	SecurityContext *corev1.PodSecurityContext `json:"securityContext,omitempty"`

	// ImagePullSecrets is an optional list of references to secrets in the same
	// namespace to use for pulling any of the images used by this Runner
	// +optional
	ImagePullSecrets []corev1.LocalObjectReference `json:"imagePullSecrets,omitempty"`

	// If specified, the pod's scheduling constraints.
	// +optional
	Affinity *corev1.Affinity `json:"affinity,omitempty"`

	// If specified, the pod's tolerations.
	// +optional
	Tolerations []corev1.Toleration `json:"tolerations,omitempty"`

	// If specified, indicates the pod's priority. "system-node-critical" and
	// "system-cluster-critical" are two special keywords which indicate the
	// highest priorities with the former being the highest priority. Any other
	// name must be defined by creating a PriorityClass object with that name. If
	// not specified, the pod priority will be default or zero if there is no
	// default.
	// +optional
	PriorityClassName string `json:"priorityClassName,omitempty"`

	// Optional duration in seconds the pod needs to terminate gracefully. May be
	// decreased in delete request. Value must be non-negative integer. The value
	// zero indicates stop immediately via the kill signal (no opportunity to shut
	// down). If this value is nil, the default grace period will be used instead.
	// +optional
	TerminationGracePeriodSeconds *int64 `json:"terminationGracePeriodSeconds,omitempty"`

	// List of ephemeral containers run in this pod. Ephemeral containers may be
	// run in an existing pod to perform user-initiated actions such as debugging.
	// +optional
	EphemeralContainers []corev1.EphemeralContainer `json:"ephemeralContainers,omitempty"`

	// HostAliases is an optional list of hosts and IPs that will be injected into
	// the pod's hosts file if specified. This is only valid for non-hostNetwork pods.
	// +optional
	HostAliases []corev1.HostAlias `json:"hostAliases,omitempty"`

	// TopologySpreadConstraints describes how a group of pods ought to spread
	// across topology domains. Scheduler will schedule pods in a way which abides
	// by the constraints. All topologySpreadConstraints are ANDed.
	// +optional
	TopologySpreadConstraints []corev1.TopologySpreadConstraint `json:"topologySpreadConstraints,omitempty"`

	// RuntimeClassName is the container runtime configuration that containers should run under.
	// More info: https://kubernetes.io/docs/concepts/containers/runtime-class
	// +optional
	RuntimeClassName *string `json:"runtimeClassName,omitempty"`

	// Specifies the DNS parameters of a pod. Parameters specified here will be
	// merged to the generated DNS configuration based on DNSPolicy.
	// +optional
	DnsConfig *corev1.PodDNSConfig `json:"dnsConfig,omitempty"`

	// +optional
	WorkVolumeClaimTemplate *WorkVolumeClaimTemplate `json:"workVolumeClaimTemplate,omitempty"`
}

func (rs *RunnerSpec) Validate(rootPath *field.Path) field.ErrorList {
	var (
		errList field.ErrorList
		err     error
	)

	err = rs.validateRepository()
	if err != nil {
		errList = append(errList, field.Invalid(rootPath.Child("repository"), rs.Repository, err.Error()))
	}

	err = rs.validateWorkVolumeClaimTemplate()
	if err != nil {
		errList = append(errList, field.Invalid(rootPath.Child("workVolumeClaimTemplate"), rs.WorkVolumeClaimTemplate, err.Error()))
	}

	return errList
}

// ValidateRepository validates repository field.
func (rs *RunnerSpec) validateRepository() error {
	// Enterprise, Organization and repository are both exclusive.
	foundCount := 0
	if len(rs.Organization) > 0 {
		foundCount += 1
	}
	if len(rs.Repository) > 0 {
		foundCount += 1
	}
	if len(rs.Enterprise) > 0 {
		foundCount += 1
	}
	if foundCount == 0 {
		return errors.New("Spec needs enterprise, organization or repository")
	}
	if foundCount > 1 {
		return errors.New("Spec cannot have many fields defined enterprise, organization and repository")
	}

	return nil
}

func (rs *RunnerSpec) validateWorkVolumeClaimTemplate() error {
	if rs.ContainerMode != "kubernetes" {
		return nil
	}

	if rs.WorkVolumeClaimTemplate == nil {
		return errors.New("Spec.ContainerMode: kubernetes must have workVolumeClaimTemplate field specified")
	}

	return rs.WorkVolumeClaimTemplate.validate()
}

// RunnerStatus defines the observed state of Runner
type RunnerStatus struct {
	// Turns true only if the runner pod is ready.
	// +optional
	Ready bool `json:"ready"`
	// +optional
	Registration RunnerStatusRegistration `json:"registration"`
	// +optional
	Phase string `json:"phase,omitempty"`
	// +optional
	Reason string `json:"reason,omitempty"`
	// +optional
	Message string `json:"message,omitempty"`
	// +optional
	// +nullable
	LastRegistrationCheckTime *metav1.Time `json:"lastRegistrationCheckTime,omitempty"`
}

// RunnerStatusRegistration contains runner registration status
type RunnerStatusRegistration struct {
	Enterprise   string      `json:"enterprise,omitempty"`
	Organization string      `json:"organization,omitempty"`
	Repository   string      `json:"repository,omitempty"`
	Labels       []string    `json:"labels,omitempty"`
	Token        string      `json:"token"`
	ExpiresAt    metav1.Time `json:"expiresAt"`
}

type WorkVolumeClaimTemplate struct {
	StorageClassName string                              `json:"storageClassName"`
	AccessModes      []corev1.PersistentVolumeAccessMode `json:"accessModes"`
	Resources        corev1.ResourceRequirements         `json:"resources"`
}

func (w *WorkVolumeClaimTemplate) validate() error {
	if w.AccessModes == nil || len(w.AccessModes) == 0 {
		return errors.New("Access mode should have at least one mode specified")
	}

	for _, accessMode := range w.AccessModes {
		switch accessMode {
		case corev1.ReadWriteOnce, corev1.ReadWriteMany:
		default:
			return fmt.Errorf("Access mode %v is not supported", accessMode)
		}
	}
	return nil
}

func (w *WorkVolumeClaimTemplate) V1Volume() corev1.Volume {
	return corev1.Volume{
		Name: "work",
		VolumeSource: corev1.VolumeSource{
			Ephemeral: &corev1.EphemeralVolumeSource{
				VolumeClaimTemplate: &corev1.PersistentVolumeClaimTemplate{
					Spec: corev1.PersistentVolumeClaimSpec{
						AccessModes:      w.AccessModes,
						StorageClassName: &w.StorageClassName,
						Resources:        w.Resources,
					},
				},
			},
		},
	}
}

func (w *WorkVolumeClaimTemplate) V1VolumeMount(mountPath string) corev1.VolumeMount {
	return corev1.VolumeMount{
		MountPath: mountPath,
		Name:      "work",
	}
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:JSONPath=".spec.enterprise",name=Enterprise,type=string
// +kubebuilder:printcolumn:JSONPath=".spec.organization",name=Organization,type=string
// +kubebuilder:printcolumn:JSONPath=".spec.repository",name=Repository,type=string
// +kubebuilder:printcolumn:JSONPath=".spec.labels",name=Labels,type=string
// +kubebuilder:printcolumn:JSONPath=".status.phase",name=Status,type=string
// +kubebuilder:printcolumn:JSONPath=".status.message",name=Message,type=string
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// Runner is the Schema for the runners API
type Runner struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RunnerSpec   `json:"spec,omitempty"`
	Status RunnerStatus `json:"status,omitempty"`
}

func (r Runner) IsRegisterable() bool {
	if r.Status.Registration.Repository != r.Spec.Repository {
		return false
	}

	if r.Status.Registration.Token == "" {
		return false
	}

	now := metav1.Now()
	if r.Status.Registration.ExpiresAt.Before(&now) {
		return false
	}

	return true
}

// +kubebuilder:object:root=true

// RunnerList contains a list of Runner
type RunnerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Runner `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Runner{}, &RunnerList{})
}
