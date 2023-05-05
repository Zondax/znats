package nats

import "github.com/spf13/viper"

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
	// TODO implement me
	return nil
}
