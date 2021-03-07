package diff

import "github.com/getkin/kin-openapi/openapi3"

// OAuthFlowDiff is a diff between oauth flow objects: https://swagger.io/specification/#oauth-flow-object
type OAuthFlowDiff struct {
	Added                bool            `json:"added,omitempty"`
	Deleted              bool            `json:"deleted,omitempty"`
	ExtensionsDiff       *ExtensionsDiff `json:"extensions,omitempty"`
	AuthorizationURLDiff *ValueDiff      `json:"authorizationURL,omitempty"`
	TokenURLDiff         *ValueDiff      `json:"tokenURL,omitempty"`
	RefreshURLDiff       *ValueDiff      `json:"refresh,omitempty"`
	ScopesDiff           *StringMapDiff  `json:"scopes,omitempty"`
}

func (diff OAuthFlowDiff) empty() bool {
	return diff == OAuthFlowDiff{}
}

func getOAuthFlowDiff(config *Config, flow1, flow2 *openapi3.OAuthFlow) *OAuthFlowDiff {

	if flow1 == nil && flow2 == nil {
		return nil
	}

	result := OAuthFlowDiff{}

	if flow1 == nil && flow2 != nil {
		result.Added = true
		return &result
	}

	if flow1 != nil && flow2 == nil {
		result.Deleted = true
		return &result
	}

	result.ExtensionsDiff = getExtensionsDiff(config, flow1.ExtensionProps, flow2.ExtensionProps)
	result.AuthorizationURLDiff = getValueDiff(flow1.AuthorizationURL, flow2.AuthorizationURL)
	result.TokenURLDiff = getValueDiff(flow1.TokenURL, flow2.TokenURL)
	result.RefreshURLDiff = getValueDiff(flow1.RefreshURL, flow2.RefreshURL)
	result.ScopesDiff = getStringMapDiff(flow1.Scopes, flow2.Scopes)

	if result.empty() {
		return nil
	}

	return &result
}