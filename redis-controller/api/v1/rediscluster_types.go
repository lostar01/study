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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	PENDING               = "Pending"
	RUNNING               = "Running"
	DELETING              = "Deleting"
	DELETED               = "Deleted"
	CREATING              = "Creating"
	UPDATING              = "Updating"
	FAILED                = "Failed"
	COMPLETED             = "Completed"
	UNKNOWN               = "Unknown"
	STOPPED               = "Stopped"
	STARTING              = "Starting"
	STOPPING              = "Stopping"
	RESTARTING            = "Restarting"
	RedisClusterFinalizer = "finalizer.rediscluster.com"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// RedisclusterSpec defines the desired state of Rediscluster
type RedisclusterSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of Rediscluster. Edit rediscluster_types.go to remove/update
	Name       string `json:"name,omitempty"`
	Image      string `json:"image,omitempty"`
	Replicas   int32  `json:"replicas,omitempty"`
	MemorySize string `json:"memorySize,omitempty"`
}

// RedisclusterStatus defines the observed state of Rediscluster
type RedisclusterStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Phase          string                  `json:"phase,omitempty"`
	Conditions     []RedisclusterCondition `json:"conditions,omitempty"`
	AdditionalInfo string                  `json:"additionalInfo,omitempty"`
}

type RedisclusterCondition struct {
	Type               string      `json:"type,omitempty"`
	Status             string      `json:"status,omitempty"`
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Rediscluster is the Schema for the redisclusters API
type Rediscluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RedisclusterSpec   `json:"spec,omitempty"`
	Status RedisclusterStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// RedisclusterList contains a list of Rediscluster
type RedisclusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Rediscluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Rediscluster{}, &RedisclusterList{})
}
