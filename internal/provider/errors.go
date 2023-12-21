package provider

import "github.com/hashicorp/terraform-plugin-framework/diag"

func WatchdogServiceError(err error) diag.Diagnostic {
	return diag.NewErrorDiagnostic(
		"Watchdog Service conversion error",
		"Got error encoding watchdog service: "+err.Error(),
	)
}

func MetricProviderServiceError(err error) diag.Diagnostic {
	return diag.NewErrorDiagnostic(
		"MetricProvider Service conversion error",
		"Got error encoding metric_provider service: "+err.Error(),
	)
}

func IssueAndUpdateDestructionWarning() diag.Diagnostic {
	return diag.NewWarningDiagnostic(
		"Issue/Update Destruction Considerations",
		"Hund recommends leaving behind resolved Issues and their Updates, so that "+
			"the history remains shown on your status page. If you no longer require "+
			"Terraform to manage this Issue and its Updates, consider simply removing "+
			"these resources from your Terraform state, or set the `archive_on_destroy` "+
			"attribute to true, to squelch this warning and prevent removal of the "+
			"Issue/Update from your status page history.",
	)
}
