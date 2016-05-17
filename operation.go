/*
MIT License

Copyright (c) 2016 Timo Reimann

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package main

import (
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/errors"
	"k8s.io/kubernetes/pkg/api/resource"
	vapi "k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/apis/extensions"
	client "k8s.io/kubernetes/pkg/client/unversioned"
)

const namespace string = "default"

// operation represents a Kubernetes operation.
type operation interface {
	Do(c *client.Client)
}

type versionOperation struct{}

func (op *versionOperation) Do(c *client.Client) {
	info, err := c.Discovery().ServerVersion()
	if err != nil {
		logger.Fatalf("failed to retrieve server API version: %s\n", err)
	}

	logger.Printf("server API version information: %s\n", info)
}

type deployOperation struct {
	image string
	name string
	port int
}

func (op *deployOperation) Do(c *client.Client) {
	appName := op.name
	port := op.port

	// Define Deployments spec.
	deploySpec := &extensions.Deployment{
		TypeMeta: vapi.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "extensions/v1beta1",
		},
		ObjectMeta: api.ObjectMeta{
			Name: appName,
		},
		Spec: extensions.DeploymentSpec{
			Replicas: 1,
			Template: api.PodTemplateSpec{
				ObjectMeta: api.ObjectMeta{
					Name:   appName,
					Labels: map[string]string{"app": appName},
				},
				Spec: api.PodSpec{
					Containers: []api.Container{
						api.Container{
							Name:  appName,
							Image: op.image,
							Ports: []api.ContainerPort{
								api.ContainerPort{ContainerPort: port, Protocol: api.ProtocolTCP},
							},
							Resources: api.ResourceRequirements{
								Limits: api.ResourceList{
									api.ResourceCPU:    resource.MustParse("100m"),
									api.ResourceMemory: resource.MustParse("256Mi"),
								},
							},
							ImagePullPolicy: api.PullIfNotPresent,
						},
					},
					RestartPolicy: api.RestartPolicyAlways,
					DNSPolicy:     api.DNSClusterFirst,
				},
			},
		},
	}

	// Implement deployment update-or-create semantics.
	deploy := c.Extensions().Deployments(namespace)
	_, err := deploy.Update(deploySpec)
	switch {
	case err == nil:
		logger.Println("deployment controller updated")
	case !errors.IsNotFound(err):
		logger.Fatalf("failed to update deployment controller: %s", err)
	default:
		_, err = deploy.Create(deploySpec)
		if err != nil {
			logger.Fatalf("failed to create deployment controller: %s", err)
		}
		logger.Println("deployment controller created")
	}
}
