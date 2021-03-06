// Copyright © 2020 Banzai Cloud
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package kubesecret_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/banzaicloud/pipeline/internal/secret/kubesecret"
)

func TestCreateKubeSecret(t *testing.T) {
	tests := map[string]struct {
		kubeSecret        v1.Secret
		kubeSecretRequest kubesecret.KubeSecretRequest
	}{
		"simple secret": {
			v1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "secret",
					Namespace: "namespace",
				},
				StringData: map[string]string{
					"key": "value",
				},
			},
			kubesecret.KubeSecretRequest{
				Name:      "secret",
				Namespace: "namespace",
				Type:      "generic",
				Values: map[string]string{
					"key": "value",
				},
			},
		},
		"secret with opaque fields": {
			v1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "secret",
					Namespace: "namespace",
				},
				StringData: map[string]string{
					".htpasswd": "blah",
				},
			},
			kubesecret.KubeSecretRequest{
				Name:      "secret",
				Namespace: "namespace",
				Type:      "htpasswd",
				Values: map[string]string{
					"username":  "user",
					"password":  "pass",
					".htpasswd": "blah",
				},
			},
		},
		"secret with spec source": {
			v1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "secret",
					Namespace: "namespace",
				},
				StringData: map[string]string{
					"tls.crt": "tlscert",
					"tls.key": "tlskey",
				},
			},
			kubesecret.KubeSecretRequest{
				Name:      "secret",
				Namespace: "namespace",
				Type:      "generic",
				Values: map[string]string{
					"clientCert": "tlscert",
					"clientKey":  "tlskey",
				},
				Spec: map[string]kubesecret.KubeSecretSpecItem{
					"tls.crt": {
						Source: "clientCert",
					},
					"tls.key": {
						Source: "clientKey",
					},
				},
			},
		},
		"secret with spec source mapping": {
			v1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "secret",
					Namespace: "namespace",
				},
				StringData: map[string]string{
					"docker.json": "{\"docker_password\":\"password\",\"docker_username\":\"username\"}",
				},
			},
			kubesecret.KubeSecretRequest{
				Name:      "secret",
				Namespace: "namespace",
				Type:      "generic",
				Values: map[string]string{
					"username": "username",
					"password": "password",
				},
				Spec: map[string]kubesecret.KubeSecretSpecItem{
					"docker.json": {
						SourceMap: map[string]string{
							"docker_username": "username",
							"docker_password": "password",
						},
					},
				},
			},
		},
		"secret with empty spec item": {
			v1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "secret",
					Namespace: "namespace",
				},
				StringData: map[string]string{
					"google.json": "{\"password\":\"password\",\"username\":\"username\"}",
				},
			},
			kubesecret.KubeSecretRequest{
				Name:      "secret",
				Namespace: "namespace",
				Type:      "generic",
				Values: map[string]string{
					"username": "username",
					"password": "password",
				},
				Spec: map[string]kubesecret.KubeSecretSpecItem{
					"google.json": {},
				},
			},
		},
		"secret with plain values": {
			v1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "secret",
					Namespace: "namespace",
				},
				StringData: map[string]string{
					"key": "value",
				},
			},
			kubesecret.KubeSecretRequest{
				Name:      "secret",
				Namespace: "namespace",
				Spec: map[string]kubesecret.KubeSecretSpecItem{
					"key": {
						Value: "value",
					},
				},
			},
		},
		"secret with plain value and empty spec items": {
			v1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "secret",
					Namespace: "namespace",
				},
				StringData: map[string]string{
					"key":         "value",
					"google.json": "{\"password\":\"password\",\"username\":\"username\"}",
				},
			},
			kubesecret.KubeSecretRequest{
				Name:      "secret",
				Namespace: "namespace",
				Type:      "generic",
				Values: map[string]string{
					"username": "username",
					"password": "password",
				},
				Spec: map[string]kubesecret.KubeSecretSpecItem{
					"key": {
						Value: "value",
					},
					"google.json": {},
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			kubeSecret, err := kubesecret.CreateKubeSecret(test.kubeSecretRequest)
			require.NoError(t, err)

			assert.Equal(t, test.kubeSecret, kubeSecret)
		})
	}
}
