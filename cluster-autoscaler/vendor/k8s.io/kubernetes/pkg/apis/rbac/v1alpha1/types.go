/*
Copyright 2016 The Kubernetes Authors All rights reserved.

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
	"k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/api/v1"
	"k8s.io/kubernetes/pkg/runtime"
)

// Authorization is calculated against
// 1. evaluation of ClusterRoleBindings - short circuit on match
// 2. evaluation of RoleBindings in the namespace requested - short circuit on match
// 3. deny by default

// PolicyRule holds information that describes a policy rule, but does not contain information
// about who the rule applies to or which namespace the rule applies to.
type PolicyRule struct {
	// Verbs is a list of Verbs that apply to ALL the ResourceKinds and AttributeRestrictions contained in this rule.  VerbAll represents all kinds.
	Verbs []string `json:"verbs" protobuf:"bytes,1,rep,name=verbs"`
	// AttributeRestrictions will vary depending on what the Authorizer/AuthorizationAttributeBuilder pair supports.
	// If the Authorizer does not recognize how to handle the AttributeRestrictions, the Authorizer should report an error.
	AttributeRestrictions runtime.RawExtension `json:"attributeRestrictions,omitempty" protobuf:"bytes,2,opt,name=attributeRestrictions"`
	// APIGroups is the name of the APIGroup that contains the resources.  If multiple API groups are specified, any action requested against one of
	// the enumerated resources in any API group will be allowed.
	APIGroups []string `json:"apiGroups" protobuf:"bytes,3,rep,name=apiGroups"`
	// Resources is a list of resources this rule applies to.  ResourceAll represents all resources.
	Resources []string `json:"resources" protobuf:"bytes,4,rep,name=resources"`
	// ResourceNames is an optional white list of names that the rule applies to.  An empty set means that everything is allowed.
	ResourceNames []string `json:"resourceNames,omitempty" protobuf:"bytes,5,rep,name=resourceNames"`
	// NonResourceURLsSlice is a set of partial urls that a user should have access to.  *s are allowed, but only as the full, final step in the path
	// This name is intentionally different than the internal type so that the DefaultConvert works nicely and because the ordering may be different.
	// Since non-resource URLs are not namespaced, this field is only applicable for ClusterRoles referenced from a ClusterRoleBinding.
	NonResourceURLs []string `json:"nonResourceURLs,omitempty" protobuf:"bytes,6,rep,name=nonResourceURLs"`
}

// Subject contains a reference to the object or user identities a role binding applies to.  This can either hold a direct API object reference,
// or a value for non-objects such as user and group names.
type Subject struct {
	// Kind of object being referenced. Values defined by this API group are "User", "Group", and "ServiceAccount".
	// If the Authorizer does not recognized the kind value, the Authorizer should report an error.
	Kind string `json:"kind" protobuf:"bytes,1,opt,name=kind"`
	// APIVersion holds the API group and version of the referenced object. For non-object references such as "Group" and "User" this is
	// expected to be API version of this API group. For example "rbac/v1alpha1".
	APIVersion string `json:"apiVersion" protobuf:"bytes,2,opt.name=apiVersion"`
	// Name of the object being referenced.
	Name string `json:"name" protobuf:"bytes,3,opt,name=name"`
	// Namespace of the referenced object.  If the object kind is non-namespace, such as "User" or "Group", and this value is not empty
	// the Authorizer should report an error.
	Namespace string `json:"namespace,omitempty" protobuf:"bytes,4,opt,name=namespace"`
}

// +genclient=true

// Role is a namespaced, logical grouping of PolicyRules that can be referenced as a unit by a RoleBinding.
type Role struct {
	unversioned.TypeMeta `json:",inline"`
	// Standard object's metadata.
	v1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Rules holds all the PolicyRules for this Role
	Rules []PolicyRule `json:"rules" protobuf:"bytes,2,rep,name=rules"`
}

// +genclient=true

// RoleBinding references a role, but does not contain it.  It can reference a Role in the same namespace or a ClusterRole in the global namespace.
// It adds who information via Subjects and namespace information by which namespace it exists in.  RoleBindings in a given
// namespace only have effect in that namespace.
type RoleBinding struct {
	unversioned.TypeMeta `json:",inline"`
	// Standard object's metadata.
	v1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Subjects holds references to the objects the role applies to.
	Subjects []Subject `json:"subject" protobuf:"bytes,2,rep,name=subjects"`

	// RoleRef can reference a Role in the current namespace or a ClusterRole in the global namespace.
	// If the RoleRef cannot be resolved, the Authorizer must return an error.
	RoleRef v1.ObjectReference `json:"roleRef" protobuf:"bytes,3,opt,name=roleRef"`
}

// RoleBindingList is a collection of RoleBindings
type RoleBindingList struct {
	unversioned.TypeMeta `json:",inline"`
	// Standard object's metadata.
	unversioned.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Items is a list of RoleBindings
	Items []RoleBinding `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// RoleList is a collection of Roles
type RoleList struct {
	unversioned.TypeMeta `json:",inline"`
	// Standard object's metadata.
	unversioned.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Items is a list of Roles
	Items []Role `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// +genclient=true,nonNamespaced=true

// ClusterRole is a cluster level, logical grouping of PolicyRules that can be referenced as a unit by a RoleBinding or ClusterRoleBinding.
type ClusterRole struct {
	unversioned.TypeMeta `json:",inline"`
	// Standard object's metadata.
	v1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Rules holds all the PolicyRules for this ClusterRole
	Rules []PolicyRule `json:"rules" protobuf:"bytes,2,rep,name=rules"`
}

// +genclient=true,nonNamespaced=true

// ClusterRoleBinding references a ClusterRole, but not contain it.  It can reference a ClusterRole in the global namespace,
// and adds who information via Subject.
type ClusterRoleBinding struct {
	unversioned.TypeMeta `json:",inline"`
	// Standard object's metadata.
	v1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Subjects holds references to the objects the role applies to.
	Subjects []Subject `json:"subject" protobuf:"bytes,2,rep,name=subjects"`

	// RoleRef can only reference a ClusterRole in the global namespace.
	// If the RoleRef cannot be resolved, the Authorizer must return an error.
	RoleRef v1.ObjectReference `json:"roleRef" protobuf:"bytes,3,opt,name=roleRef"`
}

// ClusterRoleBindingList is a collection of ClusterRoleBindings
type ClusterRoleBindingList struct {
	unversioned.TypeMeta `json:",inline"`
	// Standard object's metadata.
	unversioned.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Items is a list of ClusterRoleBindings
	Items []ClusterRoleBinding `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// ClusterRoleList is a collection of ClusterRoles
type ClusterRoleList struct {
	unversioned.TypeMeta `json:",inline"`
	// Standard object's metadata.
	unversioned.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Items is a list of ClusterRoles
	Items []ClusterRole `json:"items" protobuf:"bytes,2,rep,name=items"`
}
