package models

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func I18nStringError(err error) diag.Diagnostic {
	return diag.NewErrorDiagnostic(
		"I18nString conversion error",
		"Got error encoding I18nString: "+err.Error(),
	)
}

func TimestampError(err error) diag.Diagnostic {
	return diag.NewErrorDiagnostic(
		"Timestamp conversion error",
		"Got error encoding Timestamp: "+err.Error(),
	)
}

func UnknownServiceError(m string) diag.Diagnostic {
	return diag.NewErrorDiagnostic(
		"Service decoding error",
		fmt.Sprintf("Got an unknown Service type: %v", m),
	)
}

func ServiceDecodeError(t string, err error) diag.Diagnostic {
	return diag.NewErrorDiagnostic(
		"Service decoding error",
		"Got error decoding Service object as "+t+": "+err.Error(),
	)
}

func UnknownNativeServiceError(m string) diag.Diagnostic {
	return diag.NewErrorDiagnostic(
		"Native Service decoding error",
		fmt.Sprintf("Got an unknown Native Service method: %v", m),
	)
}

func NativeServiceDecodeError(m string, err error) diag.Diagnostic {
	return diag.NewErrorDiagnostic(
		"Native Service decoding error",
		"Got error decoding Native Service object as "+m+": "+err.Error(),
	)
}
