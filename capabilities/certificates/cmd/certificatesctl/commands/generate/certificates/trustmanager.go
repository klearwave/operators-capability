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

package certificates

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"sigs.k8s.io/controller-runtime/pkg/client"

	// common imports for subcommands
	cmdgenerate "github.com/klearwave/operators-capability/capabilities/certificates/cmd/certificatesctl/commands/generate"
	// specific imports for workloads
	v1alpha1trustmanager "github.com/klearwave/operators-capability/capabilities/certificates/apis/certificates/v1alpha1/trustmanager"
	// +kubebuilder:scaffold:operator-builder:imports
)

// NewTrustManagerSubCommand creates a new command and adds it to its
// parent command.
func NewTrustManagerSubCommand(parentCommand *cobra.Command) {
	generateCmd := &cmdgenerate.GenerateSubCommand{
		Name:         "generate",
		Description:  "generate child resource manifests from a workload's custom resource",
		SubCommandOf: parentCommand,
		GenerateFunc: GenerateTrustManager, UseWorkloadManifest: true,
		WorkloadKind: "TrustManager",
	}

	generateCmd.Setup()
}

// GenerateTrustManager runs the logic to generate child resources for a
// TrustManager workload.
func GenerateTrustManager(g *cmdgenerate.GenerateSubCommand) error {
	var apiVersion string

	workloadFilename, _ := filepath.Abs(g.WorkloadManifest)
	workloadFile, err := os.ReadFile(workloadFilename)
	if err != nil {
		return fmt.Errorf("failed to open workload file %s, %w", workloadFile, err)
	}

	var workload map[string]interface{}

	if err := yaml.Unmarshal(workloadFile, &workload); err != nil {
		return fmt.Errorf("failed to unmarshal yaml into workload, %w", err)
	}

	workloadGroupVersion := strings.Split(workload["apiVersion"].(string), "/")
	workloadAPIVersion := workloadGroupVersion[len(workloadGroupVersion)-1]

	apiVersion = workloadAPIVersion

	// generate a map of all versions to generate functions for each api version created
	type generateFunc func([]byte) ([]client.Object, error)
	generateFuncMap := map[string]generateFunc{
		"v1alpha1": v1alpha1trustmanager.GenerateForCLI,
		// +kubebuilder:scaffold:operator-builder:versionmap
	}

	generate := generateFuncMap[apiVersion]
	resourceObjects, err := generate(workloadFile)
	if err != nil {
		return fmt.Errorf("unable to retrieve resources; %w", err)
	}

	e := json.NewYAMLSerializer(json.DefaultMetaFactory, nil, nil)

	outputStream := os.Stdout

	for _, o := range resourceObjects {
		if _, err := outputStream.WriteString("---\n"); err != nil {
			return fmt.Errorf("failed to write output, %w", err)
		}

		if err := e.Encode(o, os.Stdout); err != nil {
			return fmt.Errorf("failed to write output, %w", err)
		}
	}

	return nil
}
