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
	"github.com/cloudhut/common/logging"
	"github.com/cloudhut/common/rest"
	"github.com/redpanda-data/console/backend/pkg/connect"
	"github.com/redpanda-data/console/backend/pkg/console"
	"github.com/redpanda-data/console/backend/pkg/msgpack"
	"github.com/redpanda-data/console/backend/pkg/proto"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ConsoleSpec defines the desired state of Console
type ConsoleSpec struct {
	Console console.Config  `json:"console"`
	Connect connect.Config  `json:"connect"`
	REST    rest.Config     `json:"server"`
	Logger  *logging.Config `json:"logger,omitempty"`

	MessagePack msgpack.Config `json:"messagePack"`
	Protobuf    proto.Config   `json:"protobuf"`

	// ClusterKeyRef references to Cluster Custom Resource which will
	// set kafka.Config configuration in Redpanda console
	ClusterKeyRef corev1.ObjectReference `json:"clusterKeyRef"`
	ClientID      string                 `json:"clientID,omitempty"`

	Deployment Deployment `json:"deployment"`

	Image   string `json:"image"`
	Version string `json:"version"`
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
	// Current state of the console.
	// +optional
	Conditions []ConsoleCondition `json:"conditions,omitempty"`
}

// ConsoleConditionType is a valid value for ConsoleCondition.Type
type ConsoleConditionType string

// These are valid conditions of the console.
const (
	// ConsoleAvailableConditionType indicates that all Console resources are created but does not mean ready
	ConsoleAvailableConditionType ConsoleConditionType = "ConsoleAvailable"
)

// Condition contains details for the current conditions of the resource
type ConsoleCondition struct {
	// Type is the type of the condition
	Type ConsoleConditionType `json:"type"`
	// Status is the status of the condition
	Status corev1.ConditionStatus `json:"status"`
	// Last time the condition transitioned from one status to another
	// +optional
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`
	// Unique, one-word, CamelCase reason for the condition's last transition
	// +optional
	Reason string `json:"reason,omitempty"`
	// Human-readable message indicating details about last transition
	// +optional
	Message string `json:"message,omitempty"`
}

// GetCondition return the condition of the given type
func (s *ConsoleStatus) GetCondition(cType ConsoleConditionType) *ConsoleCondition {
	for i := range s.Conditions {
		if s.Conditions[i].Type == cType {
			return &s.Conditions[i]
		}
	}
	return nil
}

// SetCondition allows setting a condition of a given type.
// In case of change in any value other than the lastTransitionTime, the lastTransitionTime
// field will be set to the current timestamp. The return value indicates if a change has happened.
func (s *ConsoleStatus) SetCondition(cType ConsoleConditionType, status corev1.ConditionStatus, reason, message string) bool {
	update := func(c *ConsoleCondition) bool {
		changed := c.Status != status || c.Reason != reason || c.Message != message
		if changed {
			c.LastTransitionTime = metav1.Now()
		}
		c.Type = cType
		c.Status = status
		c.Reason = reason
		c.Message = message
		return changed
	}
	// Try updating existing condition
	for i := range s.Conditions {
		if s.Conditions[i].Type == cType {
			return update(&s.Conditions[i])
		}
	}
	// Add a new one if missing
	newCond := ConsoleCondition{}
	update(&newCond)
	s.Conditions = append(s.Conditions, newCond)
	return true
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
