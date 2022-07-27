// Copyright 2022 Redpanda Data, Inc.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.md
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0

package v1alpha1

import (
	"time"

	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ConsoleSpec defines the desired state of Console
type ConsoleSpec struct {
	Console  ConsoleConfig `json:"console"`
	Connect  ConnectConfig `json:"connect"`
	REST     RESTConfig    `json:"server"`
	LogLevel string        `json:"logLevel,omitempty"`

	MessagePack MsgpackConfig `json:"messagePack"`
	Protobuf    ProtoConfig   `json:"protobuf"`

	// ClusterKeyRef references to Cluster Custom Resource which will
	// set kafka.Config configuration in Redpanda console
	ClusterKeyRef corev1.ObjectReference `json:"clusterKeyRef"`
	ClientID      string                 `json:"clientID,omitempty"`

	Deployment Deployment `json:"deployment"`

	Image   string `json:"image"`
	Version string `json:"version"`
}

// Config for a HTTP server
// Grabbed from https://github.com/cloudhut/common
// Added JSON tags. We can have this merged from upstream.
type RESTConfig struct {
	ServerGracefulShutdownTimeout time.Duration `yaml:"gracefulShutdownTimeout" json:"gracefulShutdownTimeout"`

	HTTPListenAddress      string        `yaml:"listenAddress" json:"listenAddress"`
	HTTPListenPort         int           `yaml:"listenPort" json:"listenPort"`
	HTTPServerReadTimeout  time.Duration `yaml:"readTimeout" json:"readTimeout"`
	HTTPServerWriteTimeout time.Duration `yaml:"writeTimeout" json:"writeTimeout"`
	HTTPServerIdleTimeout  time.Duration `yaml:"idleTimeout" json:"idleTimeout"`

	CompressionLevel int `yaml:"compressionLevel" json:"compressionLevel"`

	BasePath                        string `yaml:"basePath" json:"basePath"`
	SetBasePathFromXForwardedPrefix bool   `yaml:"setBasePathFromXForwardedPrefix" json:"setBasePathFromXForwardedPrefix"`
	StripPrefix                     bool   `yaml:"stripPrefix" json:"stripPrefix"`
}

type Deployment struct {
	Replicas                  *int32                            `json:"replicas,omitempty"`
	DeploymentStrategy        *v1.DeploymentStrategy            `json:"strategy,omitempty"`
	ProgressDeadlineSeconds   *int32                            `json:"progressDeadlineSeconds,omitempty"`
	DNSPolicy                 *corev1.DNSPolicy                 `json:"dnsPolicy,omitempty"`
	NodeSelector              map[string]string                 `json:"nodeSelector,omitempty"`
	SecurityContext           *corev1.PodSecurityContext        `json:"securityContext,omitempty"`
	ImagePullSecrets          []corev1.LocalObjectReference     `json:"imagePullSecrets,omitempty"`
	Affinity                  *corev1.Affinity                  `json:"affinity,omitempty"`
	SchedulerName             string                            `json:"schedulerName,omitempty"`
	Tolerations               []corev1.Toleration               `json:"tolerations,omitempty"`
	PriorityClassName         string                            `json:"priorityClassName,omitempty"`
	Priority                  *int32                            `json:"priority,omitempty"`
	DNSConfig                 *corev1.PodDNSConfig              `json:"dnsConfig,omitempty"`
	RuntimeClassName          *string                           `json:"runtimeClassName,omitempty"`
	Overhead                  corev1.ResourceList               `json:"overhead,omitempty"`
	TopologySpreadConstraints []corev1.TopologySpreadConstraint `json:"topologySpreadConstraints,omitempty"`
	SetHostnameAsFQDN         *bool                             `json:"setHostnameAsFQDN,omitempty"`
	ResourceRequirements      corev1.ResourceRequirements       `json:"resources,omitempty"`
	ImagePullPolicy           corev1.PullPolicy                 `json:"imagePullPolicy,omitempty"`
	LivenessProbe             *corev1.Probe                     `json:"livenessProbe,omitempty"`
	ReadinessProbe            *corev1.Probe                     `json:"readinessProbe,omitempty"`
}

// ConsoleStatus defines the observed state of Console
type ConsoleStatus struct {
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Console is the Schema for the consoles API
type Console struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ConsoleSpec   `json:"spec,omitempty"`
	Status ConsoleStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ConsoleList contains a list of Console
type ConsoleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Console `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Console{}, &ConsoleList{})
}
