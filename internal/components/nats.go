package nats

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
	"strings"
)

func NewNatsComponent(config ConfigNats) (*ComponentNats, error) {
	err, jsContext, natsConn := newJetStreamServer(config.ServerUrl, Credential{
		JWT:  config.CredentialJWT,
		Seed: config.CredentialSeed,
	})
	if err != nil {
		return nil, err
	}

	return &ComponentNats{
		Config:         config,
		JsContext:      jsContext,
		NatsConn:       natsConn,
		MapKVStore:     make(map[string]nats.KeyValue, 0),
		Streams:        make(map[string]StreamNats, 0),
		MapObjectStore: make(map[string]nats.ObjectStore, 0),
		InputTopics:    make(map[string]*Topic, 0),
		OutputTopics:   make(map[string]*Topic, 0),
		natCLI:         make(map[string]func(*nats.Msg), 0),
	}, nil
}

func newJetStreamServer(server string, credential Credential) (error, nats.JetStreamContext, *nats.Conn) {
	var nc *nats.Conn
	var err error

	if !credential.isEmpty() {
		// Connect to NATS with provided credentials
		zap.S().Infof("Attempting to connect to nats server in %s using JWT and seed...", server)
		nc, err = nats.Connect(server, nats.UserJWTAndSeed(credential.JWT, credential.Seed))
		if err != nil {
			return err, nil, nil
		}
	} else {
		// Connect to NATS
		zap.S().Infof("Attempting to connect to nats server in %s ...", server)
		nc, err = nats.Connect(server)
		if err != nil {
			return err, nil, nil
		}
	}

	// Create JetStream Context
	js, err := nc.JetStream()
	if err != nil {
		return err, nil, nil
	}

	zap.S().Infof("Successfully connected to nats server")
	return nil, js, nc
}

func GetResourcePrefix(prefixes []string, category ResourceCategory, separator string) string {
	catString := ""
	if category != NoCategory {
		catString = fmt.Sprintf("%s%s", category.String(), separator)
	}

	prefix := strings.Join(prefixes, separator)
	return fmt.Sprintf("%s%s%s",
		catString,
		prefix,
		separator,
	)
}
