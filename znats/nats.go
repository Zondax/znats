package znats

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
	"reflect"
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
			Callback: c.ping,
			Global:   true,
		},
	})
}

func (c *ComponentNats) ping(msg *nats.Msg) {
	zap.S().Infof("Received ping request")
	var response PingResponse
	response.Name = c.Config.ServiceName
	response.InputTopics = getMapKeys(c.InputTopics)
	response.OutputTopics = getMapKeys(c.OutputTopics)
	response.CliCommands = getMapKeys(c.natCLI)
	response.Streams = getMapKeys(c.Streams)

	res, _ := json.Marshal(response)
	_ = msg.Respond(res)
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

func getMapKeys(m interface{}) []string {
	v := reflect.ValueOf(m)
	if v.Kind() != reflect.Map {
		return nil
	}

	keys := v.MapKeys()
	// convert keys to string slice
	strKeys := make([]string, len(keys))
	for i := range keys {
		strKeys[i] = keys[i].String()
	}

	return strKeys
}
