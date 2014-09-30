package common

type RegistrySync struct {
	Registry map[string][]string
}

func NewRegistrySync(registry map[string][]string) *RegistrySync {
	return &RegistrySync{registry}
}
