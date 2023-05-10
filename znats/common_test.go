package znats

import (
	"fmt"
	"os"
	"testing"
)

var (
	natsComponent *ComponentNats
	prefixes      = []string{"TST", "testnet"}
)

func setup() {
	fmt.Println("Setting up nats component ...")
	var err error
	natsComponent, err = NewNatsComponent(ConfigNats{
		ServerUrl:        "http://0.0.0.0:4222",
		CredentialJWT:    "",
		CredentialSeed:   "",
		ResourcePrefixes: prefixes,
		ServiceName:      "test",
	})

	if err != nil {
		fmt.Printf("Error while creating JetStreamServer: %v", err)
		os.Exit(1)
	}
}

func teardown() {
	fmt.Println("Shutting down nats component ...")
	if natsComponent != nil {
		natsComponent.NatsConn.Close()
	}
}

func TestMain(m *testing.M) {
	setup()
	m.Run()
	teardown()
}
