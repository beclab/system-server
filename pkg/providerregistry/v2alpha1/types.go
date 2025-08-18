package v2alpha1

type ProviderRegisterRequest struct {
	AppName      string     `json:"app_name"`
	AppNamespace string     `json:"app_namespace"`
	Providers    []Provider `json:"providers"`
}

type Provider struct {
	Domain  string   `json:"domain"`
	Service string   `json:"service"`
	Paths   []string `json:"paths"`
	Verbs   []string `json:"verbs"`
}
