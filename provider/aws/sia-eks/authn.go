//
// Copyright The Athenz Authors
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
//

package sia

import (
	"io"
	"os"

	"github.com/AthenZ/athenz/libs/go/sia/aws/options"
	"github.com/AthenZ/athenz/libs/go/sia/logutil"
)

func GetEKSPodId() string {
	podId := os.Getenv("HOSTNAME")
	if podId == "" {
		podId = "eksPod"
	}
	return podId
}

func GetEKSConfig(configFile, metaEndpoint string, useRegionalSTS bool, region string, sysLogger io.Writer) (*options.Config, *options.ConfigAccount, error) {

	config, configAccount, err := options.InitFileConfig(configFile, metaEndpoint, useRegionalSTS, region, "", sysLogger)
	if err != nil {
		logutil.LogInfo(sysLogger, "Unable to process configuration file '%s': %v\n", configFile, err)
		logutil.LogInfo(sysLogger, "Trying to determine service details from the environment variables...\n")
		config, configAccount, err = options.InitEnvConfig(config)
		if err != nil {
			logutil.LogInfo(sysLogger, "Unable to process environment settings: %v\n", err)
			// if we do not have settings in our environment, we're going
			// to use fallback to <domain>.<service>-service naming structure
			logutil.LogInfo(sysLogger, "Trying to determine service name security credentials...\n")
			configAccount, err = options.InitCredsConfig("-service", useRegionalSTS, region, sysLogger)
			if err != nil {
				logutil.LogInfo(sysLogger, "Unable to process security credentials: %v\n", err)
				logutil.LogInfo(sysLogger, "Trying to determine service name from profile arn...\n")
				configAccount, err = options.InitProfileConfig(metaEndpoint, "-service")
				if err != nil {
					logutil.LogInfo(sysLogger, "Unable to determine service name: %v\n", err)
					return config, nil, err
				}
			}
		}
	}
	return config, configAccount, nil
}
