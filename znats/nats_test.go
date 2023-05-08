package znats

import (
	"github.com/nats-io/nats.go"
	"gotest.tools/assert"
	"testing"
)

func TestPingService(t *testing.T) {
	// test ping with default service config
	res := natsComponent.Ping()
	want := PingResponse{
		Name:         "test",
		InputTopics:  []string{},
		OutputTopics: []string{},
		CliCommands:  []string{"TST-testnet-ping", "TST-testnet-list-commands"},
		Streams:      []string{},
	}

	assert.DeepEqual(t, res, want)

	// add a stream and compare
	streamName := "myteststream"
	err := natsComponent.CreateStream(ConfigStream{
		CommonResourceConfig: CommonResourceConfig{
			Category: CategoryData,
			Prefixes: prefixes,
		},
		NatsStreamConfig: &nats.StreamConfig{
			Name: streamName,
		},
	})

	assert.NilError(t, err)

	res = natsComponent.Ping()
	want.Streams = []string{"data_TST_testnet_myteststream"}
	assert.DeepEqual(t, res, want)

	// add an input topic and compare
	inputTopicName := "mytestinputtopic"
	inputTopic := NewTopic(&TopicConfig{
		CommonResourceConfig: CommonResourceConfig{
			Category: CategoryData,
			Prefixes: prefixes,
		},
		Name:    inputTopicName,
		Subject: "testinputsubject",
	})

	natsComponent.AddInputTopic(inputTopic)

	res = natsComponent.Ping()
	want.InputTopics = []string{"data.TST.testnet.testinputsubject"}
	assert.DeepEqual(t, res, want)

	// add an output topic and compare
	outputTopicName := "mytestoutputtopic"
	outputTopic := NewTopic(&TopicConfig{
		CommonResourceConfig: CommonResourceConfig{
			Category: CategoryData,
			Prefixes: prefixes,
		},
		Name:    outputTopicName,
		Subject: "testoutputsubject",
	})

	natsComponent.AddOutputTopic(outputTopic)

	res = natsComponent.Ping()
	want.OutputTopics = []string{"data.TST.testnet.testoutputsubject"}
	assert.DeepEqual(t, res, want)
}
