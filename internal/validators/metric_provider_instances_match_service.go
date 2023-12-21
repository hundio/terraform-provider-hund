package validators

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hundio/terraform-provider-hund/internal/models"
)

func MetricProviderInstancesMatchService() metricProviderInstancesMatchService {
	return metricProviderInstancesMatchService{expected: models.MetricProviderServiceInstances()}
}

type metricProviderInstancesMatchService struct {
	expected map[string][]string
}

var _ resource.ConfigValidator = &metricProviderInstancesMatchService{}

func (v metricProviderInstancesMatchService) Description(ctx context.Context) string {
	return "Validate a MetricProvider's given instances match that of the given service."
}

func (v metricProviderInstancesMatchService) MarkdownDescription(ctx context.Context) string {
	return "Validate a MetricProvider's given instances match that of the given service."
}

func (v metricProviderInstancesMatchService) ValidateResource(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var instances types.Map
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("instances"), &instances)...)

	var service types.Object
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("service"), &service)...)

	if instances.IsNull() || service.IsNull() {
		return
	}

	serviceType := models.MetricProviderServiceType(service)
	expected, ok := v.expected[serviceType]

	if !ok || expected == nil {
		return
	}

	instanceMap := instances.Elements()

	for _, metric := range expected {
		_, ok := instanceMap[metric]

		if !ok {
			resp.Diagnostics.AddAttributeError(
				path.Root("instances"),
				"Missing instance from MetricProvider",
				fmt.Sprintf("An expected instance for this MetricProvider's service type "+
					"(%[1]v) was missing from the configuration: %[2]v", serviceType, metric),
			)
		}
	}

	if len(instanceMap) > len(expected) {
		resp.Diagnostics.AddAttributeError(
			path.Root("instances"),
			"Extraneous instance in MetricProvider",
			fmt.Sprintf("An unexpected instance for this MetricProvider's service type "+
				"(%[1]v) was found in the configuration. The expected instances are: %[2]v", serviceType, expected),
		)
	}
}
