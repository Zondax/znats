package znats

import (
	nats2 "github.com/nats-io/nats.go"
	"gotest.tools/assert"
	"testing"
)

func TestCreateStream(t *testing.T) {
	// create the stream
	streamName := "myteststream"
	err := natsComponent.CreateStream(ConfigStream{
		CommonResourceConfig: CommonResourceConfig{
			Category: CategoryData,
			Prefixes: prefixes,
		},
		NatsStreamConfig: &nats2.StreamConfig{
			Name: streamName,
		},
	})
	if err != nil {
		t.Fatalf("Error while creating stream: %v", err)
	}

	// test full name is correct
	fullStreamName := natsComponent.Streams[streamName].fullName
	want := "data_TST_testnet_myteststream"
	assert.Equal(t, fullStreamName, want)

	// check if stream is stored in nats component
	if _, ok := natsComponent.Streams[streamName]; !ok {
		t.Fatalf("Stream '%s' not found in nats component", streamName)
	}

	// check if stream exists
	stream := natsComponent.Streams[streamName]
	if exists := natsComponent.StreamExists(stream.Info.Config.Name); !exists {
		t.Fatalf("Stream '%s' does not exist in nats server", streamName)
	}

	// delete stream
	err = natsComponent.DeleteStream(streamName)
	if err != nil {
		t.Fatalf("Error while deleting stream: %v", err)
	}

	// check if stream is deleted
	if exists := natsComponent.StreamExists(stream.Info.Config.Name); exists {
		t.Fatalf("Stream '%s' still exists in nats server", streamName)
	}
}
