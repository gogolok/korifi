//go:build !ignore_autogenerated

/*
Copyright 2021.

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

// Code generated by controller-gen. DO NOT EDIT.

package services

import (
	"k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InputParameterSchema) DeepCopyInto(out *InputParameterSchema) {
	*out = *in
	if in.Parameters != nil {
		in, out := &in.Parameters, &out.Parameters
		*out = new(runtime.RawExtension)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InputParameterSchema.
func (in *InputParameterSchema) DeepCopy() *InputParameterSchema {
	if in == nil {
		return nil
	}
	out := new(InputParameterSchema)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MaintenanceInfo) DeepCopyInto(out *MaintenanceInfo) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MaintenanceInfo.
func (in *MaintenanceInfo) DeepCopy() *MaintenanceInfo {
	if in == nil {
		return nil
	}
	out := new(MaintenanceInfo)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServiceBindingSchema) DeepCopyInto(out *ServiceBindingSchema) {
	*out = *in
	in.Create.DeepCopyInto(&out.Create)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ServiceBindingSchema.
func (in *ServiceBindingSchema) DeepCopy() *ServiceBindingSchema {
	if in == nil {
		return nil
	}
	out := new(ServiceBindingSchema)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServiceBroker) DeepCopyInto(out *ServiceBroker) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ServiceBroker.
func (in *ServiceBroker) DeepCopy() *ServiceBroker {
	if in == nil {
		return nil
	}
	out := new(ServiceBroker)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServiceBrokerCatalog) DeepCopyInto(out *ServiceBrokerCatalog) {
	*out = *in
	if in.Metadata != nil {
		in, out := &in.Metadata, &out.Metadata
		*out = new(runtime.RawExtension)
		(*in).DeepCopyInto(*out)
	}
	out.Features = in.Features
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ServiceBrokerCatalog.
func (in *ServiceBrokerCatalog) DeepCopy() *ServiceBrokerCatalog {
	if in == nil {
		return nil
	}
	out := new(ServiceBrokerCatalog)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServiceInstanceSchema) DeepCopyInto(out *ServiceInstanceSchema) {
	*out = *in
	in.Create.DeepCopyInto(&out.Create)
	in.Update.DeepCopyInto(&out.Update)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ServiceInstanceSchema.
func (in *ServiceInstanceSchema) DeepCopy() *ServiceInstanceSchema {
	if in == nil {
		return nil
	}
	out := new(ServiceInstanceSchema)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServiceOffering) DeepCopyInto(out *ServiceOffering) {
	*out = *in
	if in.Tags != nil {
		in, out := &in.Tags, &out.Tags
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Requires != nil {
		in, out := &in.Requires, &out.Requires
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.DocumentationURL != nil {
		in, out := &in.DocumentationURL, &out.DocumentationURL
		*out = new(string)
		**out = **in
	}
	in.BrokerCatalog.DeepCopyInto(&out.BrokerCatalog)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ServiceOffering.
func (in *ServiceOffering) DeepCopy() *ServiceOffering {
	if in == nil {
		return nil
	}
	out := new(ServiceOffering)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServicePlan) DeepCopyInto(out *ServicePlan) {
	*out = *in
	in.BrokerCatalog.DeepCopyInto(&out.BrokerCatalog)
	in.Schemas.DeepCopyInto(&out.Schemas)
	out.MaintenanceInfo = in.MaintenanceInfo
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ServicePlan.
func (in *ServicePlan) DeepCopy() *ServicePlan {
	if in == nil {
		return nil
	}
	out := new(ServicePlan)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServicePlanBrokerCatalog) DeepCopyInto(out *ServicePlanBrokerCatalog) {
	*out = *in
	if in.Metadata != nil {
		in, out := &in.Metadata, &out.Metadata
		*out = new(runtime.RawExtension)
		(*in).DeepCopyInto(*out)
	}
	out.Features = in.Features
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ServicePlanBrokerCatalog.
func (in *ServicePlanBrokerCatalog) DeepCopy() *ServicePlanBrokerCatalog {
	if in == nil {
		return nil
	}
	out := new(ServicePlanBrokerCatalog)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServicePlanSchemas) DeepCopyInto(out *ServicePlanSchemas) {
	*out = *in
	in.ServiceInstance.DeepCopyInto(&out.ServiceInstance)
	in.ServiceBinding.DeepCopyInto(&out.ServiceBinding)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ServicePlanSchemas.
func (in *ServicePlanSchemas) DeepCopy() *ServicePlanSchemas {
	if in == nil {
		return nil
	}
	out := new(ServicePlanSchemas)
	in.DeepCopyInto(out)
	return out
}
