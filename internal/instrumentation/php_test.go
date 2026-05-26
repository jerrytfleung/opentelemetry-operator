// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package instrumentation

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"

	"github.com/open-telemetry/opentelemetry-operator/apis/v1alpha1"
)

func TestInjectPhpSDK(t *testing.T) {
	tests := []struct {
		name string
		v1alpha1.Php
		pod              corev1.Pod
		platform         string
		expected         corev1.Pod
		err              error
		inst             v1alpha1.Instrumentation
		simulateDefaults bool
	}{
		{
			name: "PYTHONPATH not defined",
			Php:  v1alpha1.Php{Image: "foo/bar:1"},
			pod: corev1.Pod{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{},
					},
				},
			},
			platform: "glibc",
			expected: corev1.Pod{
				Spec: corev1.PodSpec{
					Volumes: []corev1.Volume{
						{
							Name: phpVolumeName,
							VolumeSource: corev1.VolumeSource{
								EmptyDir: &corev1.EmptyDirVolumeSource{
									SizeLimit: &defaultVolumeLimitSize,
								},
							},
						},
					},
					InitContainers: []corev1.Container{
						{
							Name:    "opentelemetry-auto-instrumentation-php",
							Image:   "foo/bar:1",
							Command: []string{"cp", "-r", "/autoinstrumentation/.", "/otel-auto-instrumentation-php"},
							VolumeMounts: []corev1.VolumeMount{{
								Name:      "opentelemetry-auto-instrumentation-php",
								MountPath: "/otel-auto-instrumentation-php",
							}},
						},
					},
					Containers: []corev1.Container{
						{
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "opentelemetry-auto-instrumentation-php",
									MountPath: "/otel-auto-instrumentation-php",
								},
							},
							Env: []corev1.EnvVar{
								{
									Name:  "PYTHONPATH",
									Value: fmt.Sprintf("%s:%s", "/otel-auto-instrumentation-php/opentelemetry/instrumentation/auto_instrumentation", "/otel-auto-instrumentation-php"),
								},
								{
									Name:  "OTEL_EXPORTER_OTLP_PROTOCOL",
									Value: "http/protobuf",
								},
								{
									Name:  "OTEL_TRACES_EXPORTER",
									Value: "otlp",
								},
								{
									Name:  "OTEL_METRICS_EXPORTER",
									Value: "otlp",
								},
								{
									Name:  "OTEL_LOGS_EXPORTER",
									Value: "otlp",
								},
							},
						},
					},
				},
			},
			err: nil,
		},
		{
			name: "spec.env overrides defaults",
			Php:  v1alpha1.Php{Image: "foo/bar:1"},
			pod: corev1.Pod{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{}},
				},
			},
			platform:         "glibc",
			inst:             v1alpha1.Instrumentation{Spec: v1alpha1.InstrumentationSpec{Env: []corev1.EnvVar{{Name: "OTEL_METRICS_EXPORTER", Value: "none"}}}},
			simulateDefaults: true,
			expected: corev1.Pod{
				Spec: corev1.PodSpec{
					Volumes: []corev1.Volume{{
						Name: phpVolumeName,
						VolumeSource: corev1.VolumeSource{
							EmptyDir: &corev1.EmptyDirVolumeSource{SizeLimit: &defaultVolumeLimitSize},
						},
					}},
					InitContainers: []corev1.Container{{
						Name:         "opentelemetry-auto-instrumentation-php",
						Image:        "foo/bar:1",
						Command:      []string{"cp", "-r", "/autoinstrumentation/.", "/otel-auto-instrumentation-php"},
						VolumeMounts: []corev1.VolumeMount{{Name: phpVolumeName, MountPath: "/otel-auto-instrumentation-php"}},
					}},
					Containers: []corev1.Container{{
						VolumeMounts: []corev1.VolumeMount{{Name: phpVolumeName, MountPath: "/otel-auto-instrumentation-php"}},
						Env: []corev1.EnvVar{
							{
								Name: "OTEL_NODE_IP",
								ValueFrom: &corev1.EnvVarSource{
									FieldRef: &corev1.ObjectFieldSelector{FieldPath: "status.hostIP"},
								},
							},
							{
								Name: "OTEL_POD_IP",
								ValueFrom: &corev1.EnvVarSource{
									FieldRef: &corev1.ObjectFieldSelector{FieldPath: "status.podIP"},
								},
							},
							{Name: "PYTHONPATH", Value: fmt.Sprintf("%s:%s", "/otel-auto-instrumentation-php/opentelemetry/instrumentation/auto_instrumentation", "/otel-auto-instrumentation-php")},
							{Name: "OTEL_METRICS_EXPORTER", Value: "none"},
							{Name: "OTEL_EXPORTER_OTLP_PROTOCOL", Value: "http/protobuf"},
							{Name: "OTEL_TRACES_EXPORTER", Value: "otlp"},
							{Name: "OTEL_LOGS_EXPORTER", Value: "otlp"},
						},
					}},
				},
			},
			err: nil,
		},
		{
			name: "defaults applied when no spec.env",
			Php:  v1alpha1.Php{Image: "foo/bar:1"},
			pod: corev1.Pod{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{}},
				},
			},
			platform:         "glibc",
			inst:             v1alpha1.Instrumentation{},
			simulateDefaults: true,
			expected: corev1.Pod{
				Spec: corev1.PodSpec{
					Volumes: []corev1.Volume{{
						Name: phpVolumeName,
						VolumeSource: corev1.VolumeSource{
							EmptyDir: &corev1.EmptyDirVolumeSource{SizeLimit: &defaultVolumeLimitSize},
						},
					}},
					InitContainers: []corev1.Container{{
						Name:         "opentelemetry-auto-instrumentation-php",
						Image:        "foo/bar:1",
						Command:      []string{"cp", "-r", "/autoinstrumentation/.", "/otel-auto-instrumentation-php"},
						VolumeMounts: []corev1.VolumeMount{{Name: phpVolumeName, MountPath: "/otel-auto-instrumentation-php"}},
					}},
					Containers: []corev1.Container{{
						VolumeMounts: []corev1.VolumeMount{{Name: phpVolumeName, MountPath: "/otel-auto-instrumentation-php"}},
						Env: []corev1.EnvVar{
							{
								Name: "OTEL_NODE_IP",
								ValueFrom: &corev1.EnvVarSource{
									FieldRef: &corev1.ObjectFieldSelector{FieldPath: "status.hostIP"},
								},
							},
							{
								Name: "OTEL_POD_IP",
								ValueFrom: &corev1.EnvVarSource{
									FieldRef: &corev1.ObjectFieldSelector{FieldPath: "status.podIP"},
								},
							},
							{Name: "PYTHONPATH", Value: fmt.Sprintf("%s:%s", "/otel-auto-instrumentation-php/opentelemetry/instrumentation/auto_instrumentation", "/otel-auto-instrumentation-php")},
							{Name: "OTEL_EXPORTER_OTLP_PROTOCOL", Value: "http/protobuf"},
							{Name: "OTEL_TRACES_EXPORTER", Value: "otlp"},
							{Name: "OTEL_METRICS_EXPORTER", Value: "otlp"},
							{Name: "OTEL_LOGS_EXPORTER", Value: "otlp"},
						},
					}},
				},
			},
			err: nil,
		},
		{
			name: "PYTHONPATH defined",
			Php:  v1alpha1.Php{Image: "foo/bar:1", Resources: testResourceRequirements},
			pod: corev1.Pod{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Env: []corev1.EnvVar{
								{
									Name:  "PYTHONPATH",
									Value: "/foo:/bar",
								},
							},
						},
					},
				},
			},
			platform: "glibc",
			expected: corev1.Pod{
				Spec: corev1.PodSpec{
					Volumes: []corev1.Volume{
						{
							Name: "opentelemetry-auto-instrumentation-php",
							VolumeSource: corev1.VolumeSource{
								EmptyDir: &corev1.EmptyDirVolumeSource{
									SizeLimit: &defaultVolumeLimitSize,
								},
							},
						},
					},
					InitContainers: []corev1.Container{
						{
							Name:    "opentelemetry-auto-instrumentation-php",
							Image:   "foo/bar:1",
							Command: []string{"cp", "-r", "/autoinstrumentation/.", "/otel-auto-instrumentation-php"},
							VolumeMounts: []corev1.VolumeMount{{
								Name:      "opentelemetry-auto-instrumentation-php",
								MountPath: "/otel-auto-instrumentation-php",
							}},
							Resources: testResourceRequirements,
						},
					},
					Containers: []corev1.Container{
						{
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "opentelemetry-auto-instrumentation-php",
									MountPath: "/otel-auto-instrumentation-php",
								},
							},
							Env: []corev1.EnvVar{
								{
									Name:  "PYTHONPATH",
									Value: fmt.Sprintf("%s:%s:%s", "/otel-auto-instrumentation-php/opentelemetry/instrumentation/auto_instrumentation", "/foo:/bar", "/otel-auto-instrumentation-php"),
								},
								{
									Name:  "OTEL_EXPORTER_OTLP_PROTOCOL",
									Value: "http/protobuf",
								},
								{
									Name:  "OTEL_TRACES_EXPORTER",
									Value: "otlp",
								},
								{
									Name:  "OTEL_METRICS_EXPORTER",
									Value: "otlp",
								},
								{
									Name:  "OTEL_LOGS_EXPORTER",
									Value: "otlp",
								},
							},
						},
					},
				},
			},
			err: nil,
		},
		{
			name: "OTEL_TRACES_EXPORTER defined",
			Php:  v1alpha1.Php{Image: "foo/bar:1"},
			pod: corev1.Pod{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Env: []corev1.EnvVar{
								{
									Name:  "OTEL_TRACES_EXPORTER",
									Value: "zipkin",
								},
							},
						},
					},
				},
			},
			platform: "glibc",
			expected: corev1.Pod{
				Spec: corev1.PodSpec{
					Volumes: []corev1.Volume{
						{
							Name: phpVolumeName,
							VolumeSource: corev1.VolumeSource{
								EmptyDir: &corev1.EmptyDirVolumeSource{
									SizeLimit: &defaultVolumeLimitSize,
								},
							},
						},
					},
					InitContainers: []corev1.Container{
						{
							Name:    "opentelemetry-auto-instrumentation-php",
							Image:   "foo/bar:1",
							Command: []string{"cp", "-r", "/autoinstrumentation/.", "/otel-auto-instrumentation-php"},
							VolumeMounts: []corev1.VolumeMount{{
								Name:      "opentelemetry-auto-instrumentation-php",
								MountPath: "/otel-auto-instrumentation-php",
							}},
						},
					},
					Containers: []corev1.Container{
						{
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "opentelemetry-auto-instrumentation-php",
									MountPath: "/otel-auto-instrumentation-php",
								},
							},
							Env: []corev1.EnvVar{
								{
									Name:  "OTEL_TRACES_EXPORTER",
									Value: "zipkin",
								},
								{
									Name:  "PYTHONPATH",
									Value: fmt.Sprintf("%s:%s", "/otel-auto-instrumentation-php/opentelemetry/instrumentation/auto_instrumentation", "/otel-auto-instrumentation-php"),
								},
								{
									Name:  "OTEL_EXPORTER_OTLP_PROTOCOL",
									Value: "http/protobuf",
								},
								{
									Name:  "OTEL_METRICS_EXPORTER",
									Value: "otlp",
								},
								{
									Name:  "OTEL_LOGS_EXPORTER",
									Value: "otlp",
								},
							},
						},
					},
				},
			},
			err: nil,
		},
		{
			name: "OTEL_METRICS_EXPORTER defined",
			Php:  v1alpha1.Php{Image: "foo/bar:1"},
			pod: corev1.Pod{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Env: []corev1.EnvVar{
								{
									Name:  "OTEL_METRICS_EXPORTER",
									Value: "somebackend",
								},
							},
						},
					},
				},
			},
			platform: "glibc",
			expected: corev1.Pod{
				Spec: corev1.PodSpec{
					Volumes: []corev1.Volume{
						{
							Name: "opentelemetry-auto-instrumentation-php",
							VolumeSource: corev1.VolumeSource{
								EmptyDir: &corev1.EmptyDirVolumeSource{
									SizeLimit: &defaultVolumeLimitSize,
								},
							},
						},
					},
					InitContainers: []corev1.Container{
						{
							Name:    "opentelemetry-auto-instrumentation-php",
							Image:   "foo/bar:1",
							Command: []string{"cp", "-r", "/autoinstrumentation/.", "/otel-auto-instrumentation-php"},
							VolumeMounts: []corev1.VolumeMount{{
								Name:      "opentelemetry-auto-instrumentation-php",
								MountPath: "/otel-auto-instrumentation-php",
							}},
						},
					},
					Containers: []corev1.Container{
						{
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "opentelemetry-auto-instrumentation-php",
									MountPath: "/otel-auto-instrumentation-php",
								},
							},
							Env: []corev1.EnvVar{
								{
									Name:  "OTEL_METRICS_EXPORTER",
									Value: "somebackend",
								},
								{
									Name:  "PYTHONPATH",
									Value: fmt.Sprintf("%s:%s", "/otel-auto-instrumentation-php/opentelemetry/instrumentation/auto_instrumentation", "/otel-auto-instrumentation-php"),
								},
								{
									Name:  "OTEL_EXPORTER_OTLP_PROTOCOL",
									Value: "http/protobuf",
								},
								{
									Name:  "OTEL_TRACES_EXPORTER",
									Value: "otlp",
								},
								{
									Name:  "OTEL_LOGS_EXPORTER",
									Value: "otlp",
								},
							},
						},
					},
				},
			},
			err: nil,
		},
		{
			name: "OTEL_LOGS_EXPORTER defined",
			Php:  v1alpha1.Php{Image: "foo/bar:1"},
			pod: corev1.Pod{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Env: []corev1.EnvVar{
								{
									Name:  "OTEL_LOGS_EXPORTER",
									Value: "somebackend",
								},
							},
						},
					},
				},
			},
			expected: corev1.Pod{
				Spec: corev1.PodSpec{
					Volumes: []corev1.Volume{
						{
							Name: "opentelemetry-auto-instrumentation-php",
							VolumeSource: corev1.VolumeSource{
								EmptyDir: &corev1.EmptyDirVolumeSource{
									SizeLimit: &defaultVolumeLimitSize,
								},
							},
						},
					},
					InitContainers: []corev1.Container{
						{
							Name:    "opentelemetry-auto-instrumentation-php",
							Image:   "foo/bar:1",
							Command: []string{"cp", "-r", "/autoinstrumentation/.", "/otel-auto-instrumentation-php"},
							VolumeMounts: []corev1.VolumeMount{{
								Name:      "opentelemetry-auto-instrumentation-php",
								MountPath: "/otel-auto-instrumentation-php",
							}},
						},
					},
					Containers: []corev1.Container{
						{
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "opentelemetry-auto-instrumentation-php",
									MountPath: "/otel-auto-instrumentation-php",
								},
							},
							Env: []corev1.EnvVar{
								{
									Name:  "OTEL_LOGS_EXPORTER",
									Value: "somebackend",
								},
								{
									Name:  "PYTHONPATH",
									Value: fmt.Sprintf("%s:%s", "/otel-auto-instrumentation-php/opentelemetry/instrumentation/auto_instrumentation", "/otel-auto-instrumentation-php"),
								},
								{
									Name:  "OTEL_EXPORTER_OTLP_PROTOCOL",
									Value: "http/protobuf",
								},
								{
									Name:  "OTEL_TRACES_EXPORTER",
									Value: "otlp",
								},
								{
									Name:  "OTEL_METRICS_EXPORTER",
									Value: "otlp",
								},
							},
						},
					},
				},
			},
			err: nil,
		},
		{
			name: "OTEL_EXPORTER_OTLP_PROTOCOL defined",
			Php:  v1alpha1.Php{Image: "foo/bar:1"},
			pod: corev1.Pod{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Env: []corev1.EnvVar{
								{
									Name:  "OTEL_EXPORTER_OTLP_PROTOCOL",
									Value: "somebackend",
								},
							},
						},
					},
				},
			},
			platform: "glibc",
			expected: corev1.Pod{
				Spec: corev1.PodSpec{
					Volumes: []corev1.Volume{
						{
							Name: "opentelemetry-auto-instrumentation-php",
							VolumeSource: corev1.VolumeSource{
								EmptyDir: &corev1.EmptyDirVolumeSource{
									SizeLimit: &defaultVolumeLimitSize,
								},
							},
						},
					},
					InitContainers: []corev1.Container{
						{
							Name:    "opentelemetry-auto-instrumentation-php",
							Image:   "foo/bar:1",
							Command: []string{"cp", "-r", "/autoinstrumentation/.", "/otel-auto-instrumentation-php"},
							VolumeMounts: []corev1.VolumeMount{{
								Name:      "opentelemetry-auto-instrumentation-php",
								MountPath: "/otel-auto-instrumentation-php",
							}},
						},
					},
					Containers: []corev1.Container{
						{
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "opentelemetry-auto-instrumentation-php",
									MountPath: "/otel-auto-instrumentation-php",
								},
							},
							Env: []corev1.EnvVar{
								{
									Name:  "OTEL_EXPORTER_OTLP_PROTOCOL",
									Value: "somebackend",
								},
								{
									Name:  "PYTHONPATH",
									Value: fmt.Sprintf("%s:%s", "/otel-auto-instrumentation-php/opentelemetry/instrumentation/auto_instrumentation", "/otel-auto-instrumentation-php"),
								},
								{
									Name:  "OTEL_TRACES_EXPORTER",
									Value: "otlp",
								},
								{
									Name:  "OTEL_METRICS_EXPORTER",
									Value: "otlp",
								},
								{
									Name:  "OTEL_LOGS_EXPORTER",
									Value: "otlp",
								},
							},
						},
					},
				},
			},
			err: nil,
		},
		{
			name: "PYTHONPATH defined as ValueFrom",
			Php:  v1alpha1.Php{Image: "foo/bar:1"},
			pod: corev1.Pod{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Env: []corev1.EnvVar{
								{
									Name:      "PYTHONPATH",
									ValueFrom: &corev1.EnvVarSource{},
								},
							},
						},
					},
				},
			},
			platform: "glibc",
			expected: corev1.Pod{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Env: []corev1.EnvVar{
								{
									Name:      "PYTHONPATH",
									ValueFrom: &corev1.EnvVarSource{},
								},
							},
						},
					},
				},
			},
			err: fmt.Errorf("the container defines env var value via ValueFrom, envVar: %s", envPythonPath),
		},
		{
			name: "musl platform defined",
			Php:  v1alpha1.Php{Image: "foo/bar:1"},
			pod: corev1.Pod{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{},
					},
				},
			},
			platform: "musl",
			expected: corev1.Pod{
				Spec: corev1.PodSpec{
					Volumes: []corev1.Volume{
						{
							Name: phpVolumeName,
							VolumeSource: corev1.VolumeSource{
								EmptyDir: &corev1.EmptyDirVolumeSource{
									SizeLimit: &defaultVolumeLimitSize,
								},
							},
						},
					},
					InitContainers: []corev1.Container{
						{
							Name:    "opentelemetry-auto-instrumentation-php",
							Image:   "foo/bar:1",
							Command: []string{"cp", "-r", "/autoinstrumentation-musl/.", "/otel-auto-instrumentation-php"},
							VolumeMounts: []corev1.VolumeMount{{
								Name:      "opentelemetry-auto-instrumentation-php",
								MountPath: "/otel-auto-instrumentation-php",
							}},
						},
					},
					Containers: []corev1.Container{
						{
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "opentelemetry-auto-instrumentation-php",
									MountPath: "/otel-auto-instrumentation-php",
								},
							},
							Env: []corev1.EnvVar{
								{
									Name:  "PYTHONPATH",
									Value: fmt.Sprintf("%s:%s", "/otel-auto-instrumentation-php/opentelemetry/instrumentation/auto_instrumentation", "/otel-auto-instrumentation-php"),
								},
								{
									Name:  "OTEL_EXPORTER_OTLP_PROTOCOL",
									Value: "http/protobuf",
								},
								{
									Name:  "OTEL_TRACES_EXPORTER",
									Value: "otlp",
								},
								{
									Name:  "OTEL_METRICS_EXPORTER",
									Value: "otlp",
								},
								{
									Name:  "OTEL_LOGS_EXPORTER",
									Value: "otlp",
								},
							},
						},
					},
				},
			},
			err: nil,
		},
		{
			name: "platform not defined",
			Php:  v1alpha1.Php{Image: "foo/bar:1"},
			pod: corev1.Pod{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{},
					},
				},
			},
			platform: "",
			expected: corev1.Pod{
				Spec: corev1.PodSpec{
					Volumes: []corev1.Volume{
						{
							Name: phpVolumeName,
							VolumeSource: corev1.VolumeSource{
								EmptyDir: &corev1.EmptyDirVolumeSource{
									SizeLimit: &defaultVolumeLimitSize,
								},
							},
						},
					},
					InitContainers: []corev1.Container{
						{
							Name:    "opentelemetry-auto-instrumentation-php",
							Image:   "foo/bar:1",
							Command: []string{"cp", "-r", "/autoinstrumentation/.", "/otel-auto-instrumentation-php"},
							VolumeMounts: []corev1.VolumeMount{{
								Name:      "opentelemetry-auto-instrumentation-php",
								MountPath: "/otel-auto-instrumentation-php",
							}},
						},
					},
					Containers: []corev1.Container{
						{
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "opentelemetry-auto-instrumentation-php",
									MountPath: "/otel-auto-instrumentation-php",
								},
							},
							Env: []corev1.EnvVar{
								{
									Name:  "PYTHONPATH",
									Value: fmt.Sprintf("%s:%s", "/otel-auto-instrumentation-php/opentelemetry/instrumentation/auto_instrumentation", "/otel-auto-instrumentation-php"),
								},
								{
									Name:  "OTEL_EXPORTER_OTLP_PROTOCOL",
									Value: "http/protobuf",
								},
								{
									Name:  "OTEL_TRACES_EXPORTER",
									Value: "otlp",
								},
								{
									Name:  "OTEL_METRICS_EXPORTER",
									Value: "otlp",
								},
								{
									Name:  "OTEL_LOGS_EXPORTER",
									Value: "otlp",
								},
							},
						},
					},
				},
			},
			err: nil,
		},
		{
			name: "platform not supported",
			Php:  v1alpha1.Php{Image: "foo/bar:1"},
			pod: corev1.Pod{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{},
					},
				},
			},
			platform: "not-supported",
			expected: corev1.Pod{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{},
					},
				},
			},
			err: errors.New("provided instrumentation.opentelemetry.io/otel-php-platform annotation value 'not-supported' is not supported"),
		},
		{
			name: "inject into init container",
			Php:  v1alpha1.Php{Image: "foo/bar:1"},
			pod: corev1.Pod{
				Spec: corev1.PodSpec{
					InitContainers: []corev1.Container{
						{
							Name: "my-init",
						},
					},
				},
			},
			platform: "glibc",
			expected: corev1.Pod{
				Spec: corev1.PodSpec{
					Volumes: []corev1.Volume{
						{
							Name: phpVolumeName,
							VolumeSource: corev1.VolumeSource{
								EmptyDir: &corev1.EmptyDirVolumeSource{
									SizeLimit: &defaultVolumeLimitSize,
								},
							},
						},
					},
					InitContainers: []corev1.Container{
						{
							Name:    "opentelemetry-auto-instrumentation-php",
							Image:   "foo/bar:1",
							Command: []string{"cp", "-r", "/autoinstrumentation/.", "/otel-auto-instrumentation-php"},
							VolumeMounts: []corev1.VolumeMount{{
								Name:      "opentelemetry-auto-instrumentation-php",
								MountPath: "/otel-auto-instrumentation-php",
							}},
						},
						{
							Name: "my-init",
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "opentelemetry-auto-instrumentation-php",
									MountPath: "/otel-auto-instrumentation-php",
								},
							},
							Env: []corev1.EnvVar{
								{
									Name:  "PYTHONPATH",
									Value: fmt.Sprintf("%s:%s", "/otel-auto-instrumentation-php/opentelemetry/instrumentation/auto_instrumentation", "/otel-auto-instrumentation-php"),
								},
								{
									Name:  "OTEL_EXPORTER_OTLP_PROTOCOL",
									Value: "http/protobuf",
								},
								{
									Name:  "OTEL_TRACES_EXPORTER",
									Value: "otlp",
								},
								{
									Name:  "OTEL_METRICS_EXPORTER",
									Value: "otlp",
								},
								{
									Name:  "OTEL_LOGS_EXPORTER",
									Value: "otlp",
								},
							},
						},
					},
				},
			},
			err: nil,
		},
	}

	injector := sdkInjector{}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			pod := test.pod

			// Collect all containers (regular first, then init)
			containers := allContainers(&pod)

			err := injectPhpSDK(test.Php, &pod, containers, test.platform, v1alpha1.InstrumentationSpec{})
			if err != nil {
				assert.Equal(t, test.expected, pod)
				assert.Equal(t, test.err, err)
				return
			}

			for i := range pod.Spec.Containers {
				if test.simulateDefaults {
					injector.injectCommonEnvVar(test.inst, &pod.Spec.Containers[i])
				}
				injector.injectDefaultPhpEnvVars(&pod.Spec.Containers[i])
			}
			for i := range pod.Spec.InitContainers {
				// Skip the instrumentation init container we added
				if pod.Spec.InitContainers[i].Name == phpInitContainerName {
					continue
				}
				if test.simulateDefaults {
					injector.injectCommonEnvVar(test.inst, &pod.Spec.InitContainers[i])
				}
				injector.injectDefaultPhpEnvVars(&pod.Spec.InitContainers[i])
			}
			assert.Equal(t, test.expected, pod)
			assert.Equal(t, test.err, err)
		})
	}
}

//func allContainers(pod *corev1.Pod) []*corev1.Container {
//	// Collect all containers (regular first, then init)
//	var containers []*corev1.Container
//	for i := range pod.Spec.Containers {
//		containers = append(containers, &pod.Spec.Containers[i])
//	}
//	for i := range pod.Spec.InitContainers {
//		containers = append(containers, &pod.Spec.InitContainers[i])
//	}
//	return containers
//}
