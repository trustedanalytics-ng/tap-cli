package cli

import (
	"errors"
	"github.com/trustedanalytics/tapng-cli/api"
)


func convert(serviceName, planName string) (string, string, error) {

	catalog, err := api.ConnectionConfig.ConsoleServiceApi.GetCatalog()
	if err != nil {
		return "", "", err
	}

	for _, service := range catalog {

		if service.Entity.Label == serviceName {
			for _, plan := range service.Entity.ServicePlans {

				if plan.Entity.Name == planName {
					return service.Entity.UniqueId, plan.Entity.UniqueId, nil
				}
			}
			return "", "", errors.New("cannot find plan: '" + planName + "' for service: '" + serviceName + "'")
		}
	}

	return "", "", errors.New("cannot find service: '" + serviceName + "'")


}
