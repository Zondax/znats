package znats

import (
	"encoding/json"
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

	c := &ComponentNats{
		Config:         config,
		JsContext:      jsContext,
		NatsConn:       natsConn,
		MapKVStore:     make(map[string]KVStoreNats, 0),
		Streams:        make(map[string]StreamNats, 0),
		MapObjectStore: make(map[string]ObjectStoreNats, 0),
		InputTopics:    make(map[string]*Topic, 0),
		OutputTopics:   make(map[string]*Topic, 0),
		natCLI:         make(map[string]func(*nats.Msg), 0),
	}

	c.addDefaultCliFunctions()
	return c, nil
}

func (c *ComponentNats) Shutdown() error {
	err := c.NatsConn.Drain()
	if err != nil {
		zap.S().Errorf("Error while draining nats connection: %s", err.Error())
		return err
	}
	c.NatsConn.Close()
	zap.S().Infof("Successfully shutdown nats connection")
	return nil
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

func (c *ComponentNats) addDefaultCliFunctions() {
	c.AddNatCliCmd(map[string]ReqReplyCB{
		"ping": {
			Callback: c.replyPing,
			Global:   true,
		},
		"list-commands": {
			Callback: c.replyListAvailableCliCmd,
			Global:   true,
		},
	})
}

func (c *ComponentNats) replyPing(msg *nats.Msg) {
	zap.S().Infof("Received Ping request")
	pingData := c.Ping()
	res, _ := json.Marshal(pingData)
	_ = msg.Respond(res)
}

func (c *ComponentNats) Ping() PingResponse {
	return PingResponse{
		Name:         c.Config.ServiceName,
		InputTopics:  c.GetInputTopicsFullNames(),
		OutputTopics: c.GetOutputTopicsFullNames(),
		CliCommands:  c.GetCliCmdFullNames(),
		Streams:      c.GetStreamsFullNames(),
	}
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
