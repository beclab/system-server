package v2alpha1

type RegisterResp struct {
}

type PermissionRequire struct {
	ProviderName      string  `json:"provider_name"`
	ProviderNamespace string  `json:"provider_namespace"`
	ServiceAccount    *string `json:"service_account,omitempty"`
}

type PermissionRegister struct {
	App   string              `json:"app"`
	AppID string              `json:"appid"`
	Perm  []PermissionRequire `json:"perm"`
}
