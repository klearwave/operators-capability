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
	v1alpha1certificates "github.com/klearwave/operators-capability/capabilities/certificates/apis/certificates/v1alpha1"
	v1alpha1certmanager "github.com/klearwave/operators-capability/capabilities/certificates/apis/certificates/v1alpha1/certmanager"
)

// Code generated by operator-builder. DO NOT EDIT.

// CertManagerLatestGroupVersion returns the latest group version object associated with this
// particular kind.
var CertManagerLatestGroupVersion = v1alpha1certificates.GroupVersion

// CertManagerLatestSample returns the latest sample manifest associated with this
// particular kind.
var CertManagerLatestSample = v1alpha1certmanager.Sample(false)
