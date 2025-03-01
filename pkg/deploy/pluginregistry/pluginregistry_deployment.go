//
// Copyright (c) 2019-2021 Red Hat, Inc.
// This program and the accompanying materials are made
// available under the terms of the Eclipse Public License 2.0
// which is available at https://www.eclipse.org/legal/epl-2.0/
//
// SPDX-License-Identifier: EPL-2.0
//
// Contributors:
//   Red Hat, Inc. - initial API and implementation
//
package pluginregistry

import (
	"github.com/eclipse-che/che-operator/pkg/common/chetypes"
	"github.com/eclipse-che/che-operator/pkg/common/constants"
	defaults "github.com/eclipse-che/che-operator/pkg/common/operator-defaults"
	"github.com/eclipse-che/che-operator/pkg/common/utils"
	"github.com/eclipse-che/che-operator/pkg/deploy"
	"github.com/eclipse-che/che-operator/pkg/deploy/registry"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

func (p *PluginRegistryReconciler) getPluginRegistryDeploymentSpec(ctx *chetypes.DeployContext) *appsv1.Deployment {
	registryType := "plugin"
	registryImage := defaults.GetPluginRegistryImage(ctx.CheCluster)
	registryImagePullPolicy := corev1.PullPolicy(utils.GetPullPolicyFromDockerImage(registryImage))
	probePath := "/v3/plugins/"
	pluginImagesEnv := utils.GetEnvByRegExp("^.*plugin_registry_image.*$")

	resources := corev1.ResourceRequirements{
		Requests: corev1.ResourceList{
			corev1.ResourceMemory: resource.MustParse(constants.DefaultPluginRegistryMemoryRequest),
			corev1.ResourceCPU:    resource.MustParse(constants.DefaultPluginRegistryCpuRequest),
		},
		Limits: corev1.ResourceList{
			corev1.ResourceMemory: resource.MustParse(constants.DefaultPluginRegistryMemoryLimit),
			corev1.ResourceCPU:    resource.MustParse(constants.DefaultPluginRegistryCpuLimit),
		},
	}

	deployment := registry.GetSpecRegistryDeployment(
		ctx,
		registryType,
		registryImage,
		pluginImagesEnv,
		registryImagePullPolicy,
		resources,
		probePath)

	deploy.EnsurePodSecurityStandards(deployment, constants.DefaultSecurityContextRunAsUser, constants.DefaultSecurityContextFsGroup)
	deploy.CustomizeDeployment(deployment, ctx.CheCluster.Spec.Components.PluginRegistry.Deployment)
	return deployment
}
