// Copyright (c) 2019 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package infrastructure

import (
	"time"

	"github.com/gardener/gardener-extension-provider-packet/pkg/imagevector"
	"github.com/gardener/gardener-extension-provider-packet/pkg/packet"

	"github.com/gardener/gardener/extensions/pkg/controller/common"
	"github.com/gardener/gardener/extensions/pkg/controller/infrastructure"
	"github.com/gardener/gardener/extensions/pkg/terraformer"
	glogger "github.com/gardener/gardener/pkg/logger"
	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type actuator struct {
	logger logr.Logger

	common.ChartRendererContext
}

// NewActuator creates a new Actuator that updates the status of the handled Infrastructure resources.
func NewActuator() infrastructure.Actuator {
	return &actuator{
		logger: log.Log.WithName("infrastructure-actuator"),
	}
}

// Helper functions

func (a *actuator) newTerraformer(purpose, namespace, name string) (terraformer.Terraformer, error) {
	tf, err := terraformer.NewForConfig(glogger.NewLogger("info"), a.RESTConfig(), purpose, namespace, name, imagevector.TerraformerImage())
	if err != nil {
		return nil, err
	}

	return tf.
		SetTerminationGracePeriodSeconds(630).
		SetDeadlineCleaning(5 * time.Minute).
		SetDeadlinePod(15 * time.Minute), nil
}

func generateTerraformInfraVariablesEnvironment(credentials *packet.Credentials) map[string]string {
	return map[string]string{
		"TF_VAR_PACKET_API_KEY": string(credentials.APIToken),
	}
}
