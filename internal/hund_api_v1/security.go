package hundApiV1

import (
	"github.com/deepmap/oapi-codegen/v2/pkg/securityprovider"
)

func WithSecurity(token string) (ClientOption, error) {
	bearerProvider, err := securityprovider.NewSecurityProviderBearerToken(token)

	if err != nil {
		return nil, err
	}

	return WithRequestEditorFn(bearerProvider.Intercept), nil
}
