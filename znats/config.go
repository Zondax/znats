package znats

import (
	"fmt"
	"github.com/spf13/viper"
)

type CommonResourceConfig struct {
	// Prefixes to be added to all stream, subjects, queues, etc (optional)
	Prefixes []string
	// Category is the category of the resource
	Category ResourceCategory
}

type ConfigNats struct {
	// ServerUrl is the url of the NATS server
	ServerUrl string

	// CredentialJWT (optional)
	CredentialJWT string

	// CredentialSeed (optional)
	CredentialSeed string

	// Prefixes to be added to all stream, subjects, queues, etc
	ResourcePrefixes []string

	// ServiceName is the name of the service
	ServiceName string
}

func (c ConfigNats) SetDefaults() {
	// TODO implement me
	viper.SetDefault("foo", "bar")
}

func (c ConfigNats) Validate() error {
	if c.ServerUrl == "" {
		return fmt.Errorf("server url is required")
	}

	if c.ServiceName == "" {
		return fmt.Errorf("service name is required")
	}

	return nil
}
