/*
Copyright 2018 The Kubernetes Authors.

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
	v1alpha1 "github.com/giantswarm/apiextensions/pkg/apis/provider/v1alpha1"
	scheme "github.com/giantswarm/apiextensions/pkg/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// AzureConfigsGetter has a method to return a AzureConfigInterface.
// A group's client should implement this interface.
type AzureConfigsGetter interface {
	AzureConfigs(namespace string) AzureConfigInterface
}

// AzureConfigInterface has methods to work with AzureConfig resources.
type AzureConfigInterface interface {
	Create(*v1alpha1.AzureConfig) (*v1alpha1.AzureConfig, error)
	Update(*v1alpha1.AzureConfig) (*v1alpha1.AzureConfig, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.AzureConfig, error)
	List(opts v1.ListOptions) (*v1alpha1.AzureConfigList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.AzureConfig, err error)
	AzureConfigExpansion
}

// azureConfigs implements AzureConfigInterface
type azureConfigs struct {
	client rest.Interface
	ns     string
}

// newAzureConfigs returns a AzureConfigs
func newAzureConfigs(c *ProviderV1alpha1Client, namespace string) *azureConfigs {
	return &azureConfigs{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the azureConfig, and returns the corresponding azureConfig object, and an error if there is any.
func (c *azureConfigs) Get(name string, options v1.GetOptions) (result *v1alpha1.AzureConfig, err error) {
	result = &v1alpha1.AzureConfig{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("azureconfigs").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of AzureConfigs that match those selectors.
func (c *azureConfigs) List(opts v1.ListOptions) (result *v1alpha1.AzureConfigList, err error) {
	result = &v1alpha1.AzureConfigList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("azureconfigs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested azureConfigs.
func (c *azureConfigs) Watch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("azureconfigs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a azureConfig and creates it.  Returns the server's representation of the azureConfig, and an error, if there is any.
func (c *azureConfigs) Create(azureConfig *v1alpha1.AzureConfig) (result *v1alpha1.AzureConfig, err error) {
	result = &v1alpha1.AzureConfig{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("azureconfigs").
		Body(azureConfig).
		Do().
		Into(result)
	return
}

// Update takes the representation of a azureConfig and updates it. Returns the server's representation of the azureConfig, and an error, if there is any.
func (c *azureConfigs) Update(azureConfig *v1alpha1.AzureConfig) (result *v1alpha1.AzureConfig, err error) {
	result = &v1alpha1.AzureConfig{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("azureconfigs").
		Name(azureConfig.Name).
		Body(azureConfig).
		Do().
		Into(result)
	return
}

// Delete takes name of the azureConfig and deletes it. Returns an error if one occurs.
func (c *azureConfigs) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("azureconfigs").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *azureConfigs) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("azureconfigs").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched azureConfig.
func (c *azureConfigs) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.AzureConfig, err error) {
	result = &v1alpha1.AzureConfig{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("azureconfigs").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
