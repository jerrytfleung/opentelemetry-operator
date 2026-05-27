// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package instrumentation

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"

	"github.com/open-telemetry/opentelemetry-operator/apis/v1alpha1"
)

const (
	phpInstrMountPath = "/otel-auto-instrumentation-php"

	phpIniScanDirEnvVarName = "PHP_INI_SCAN_DIR"
	// https://www.php.net/manual/en/configuration.file.php//configuration.file.scan
	//       If a blank directory is given in PHP_INI_SCAN_DIR, PHP will also scan the directory given at compile time via --with-config-file-scan-dir.
	//       PHP_INI_SCAN_DIR=:/usr/local/etc/php.d php
	//                        ^ separator after empty string
	//           PHP will load all files in /etc/php.d/*.ini, then /usr/local/etc/php.d/*.ini as configuration files.
	phpIniScanDirEnvVarValue = ":/" + phpInstrMountPath + "/php_ini_scan_dir"

	otelPhpAutoloadEnabledrEnvVarName  = "OTEL_PHP_AUTOLOAD_ENABLED"
	otelPhpAutoloadEnabledrEnvVarValue = "true"

	glibcLinuxPhpAutoInstrumentationSrc = "/autoinstrumentation/."
	muslLinuxPhpAutoInstrumentationSrc  = "/autoinstrumentation-musl/."

	phpInitContainerName = initContainerName + "-php"
	phpVolumeName        = volumeName + "-php"
)

func phpPlatformSrc(platform string) (string, error) {
	// Validate platform
	switch platform {
	case "", glibcLinux:
		return glibcLinuxPhpAutoInstrumentationSrc, nil
	case muslLinux:
		return muslLinuxPhpAutoInstrumentationSrc, nil
	default:
		return "", fmt.Errorf("provided instrumentation.opentelemetry.io/otel-php-platform annotation value '%s' is not supported", platform)
	}
}

func injectPhpSDKToContainer(phpSpec v1alpha1.Php, container *corev1.Container, platform string) error {
	volume := instrVolume(phpSpec.VolumeClaimTemplate, phpVolumeName, phpSpec.VolumeSizeLimit)

	err := validateContainerEnv(container.Env, phpIniScanDirEnvVarName, otelPhpAutoloadEnabledrEnvVarName)
	if err != nil {
		return err
	}

	_, err = phpPlatformSrc(platform)
	if err != nil {
		return err
	}

	// inject Php instrumentation spec env vars.
	container.Env = appendIfNotSet(container.Env, phpSpec.Env...)

	container.VolumeMounts = append(container.VolumeMounts, corev1.VolumeMount{
		Name:      volume.Name,
		MountPath: phpInstrMountPath,
	})

	return nil
}

func injectPhpSDKToPod(phpSpec v1alpha1.Php, pod corev1.Pod, firstContainerName, platform string, instSpec v1alpha1.InstrumentationSpec) corev1.Pod {
	volume := instrVolume(phpSpec.VolumeClaimTemplate, phpVolumeName, phpSpec.VolumeSizeLimit)

	// This has been validated already
	autoInstrumentationSrc, _ := phpPlatformSrc(platform)

	// We just inject Volumes and init containers for the first processed container.
	if isInitContainerMissing(pod, phpInitContainerName) {
		pod.Spec.Volumes = append(pod.Spec.Volumes, volume)

		initContainer := corev1.Container{
			Name:      phpInitContainerName,
			Image:     phpSpec.Image,
			Command:   []string{"cp", "-r", autoInstrumentationSrc, phpInstrMountPath},
			Resources: phpSpec.Resources,
			VolumeMounts: []corev1.VolumeMount{{
				Name:      volume.Name,
				MountPath: phpInstrMountPath,
			}},
			ImagePullPolicy: instSpec.ImagePullPolicy,
		}

		pod.Spec.InitContainers = insertInitContainer(&pod, initContainer, firstContainerName)
	}
	return pod
}

// injectPhpSDK injects PHP instrumentation into the specified containers.
// Containers must point into the provided pod and be ordered with init containers first.
func injectPhpSDK(phpSpec v1alpha1.Php, pod *corev1.Pod, containers []*corev1.Container, platform string, instSpec v1alpha1.InstrumentationSpec) error {
	for _, container := range containers {
		if err := injectPhpSDKToContainer(phpSpec, container, platform); err != nil {
			return err
		}
	}
	if len(containers) > 0 {
		*pod = injectPhpSDKToPod(phpSpec, *pod, containers[0].Name, platform, instSpec)
	}
	return nil
}

func getDefaultPhpEnvVars() []corev1.EnvVar {
	return []corev1.EnvVar{
		{
			Name:  phpIniScanDirEnvVarName,
			Value: phpIniScanDirEnvVarValue,
		},
		{
			Name:  otelPhpAutoloadEnabledrEnvVarName,
			Value: otelPhpAutoloadEnabledrEnvVarValue,
		},
	}
}
