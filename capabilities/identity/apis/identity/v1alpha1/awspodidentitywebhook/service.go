/*
Copyright 2025.

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

package awspodidentitywebhook

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/nukleros/operator-builder-tools/pkg/controller/workload"

	identityv1alpha1 "github.com/klearwave/operators-capability/capabilities/identity/apis/identity/v1alpha1"
	"github.com/klearwave/operators-capability/capabilities/identity/apis/identity/v1alpha1/awspodidentitywebhook/mutate"
)

// +kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete

// CreateServiceNamespaceAwsPodIdentityWebhook creates the Service resource with name aws-pod-identity-webhook.
func CreateServiceNamespaceAwsPodIdentityWebhook(
	parent *identityv1alpha1.AWSPodIdentityWebhook,
	reconciler workload.Reconciler,
	req *workload.Request,
) ([]client.Object, error) {

	var resourceObj = &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Service",
			"metadata": map[string]interface{}{
				"name":      "aws-pod-identity-webhook",
				"namespace": parent.Spec.Namespace, //  controlled by field: namespace
				"annotations": map[string]interface{}{
					"prometheus.io/port":   "9443",
					"prometheus.io/scheme": "https",
					"prometheus.io/scrape": "true",
				},
			},
			"spec": map[string]interface{}{
				"ports": []interface{}{
					map[string]interface{}{
						"port":       443,
						"targetPort": 9443,
					},
				},
				"selector": map[string]interface{}{
					"app":                                    "pod-identity-webhook",
					"platform.klearwave.io/capability":       "identity",
					"platform.klearwave.io/component":        "aws-pod-identity-webhook",
					"platform.klearwave.io/operator-version": "unstable",
					"platform.klearwave.io/platform-version": "unstable",
				},
			},
		},
	}

	return mutate.MutateServiceNamespaceAwsPodIdentityWebhook(resourceObj, parent, reconciler, req)
}
