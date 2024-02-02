package znats

import (
	"github.com/nats-io/nats.go"
	"gotest.tools/assert"
	"testing"
	"time"
)

func TestNatsPing(t *testing.T) {
	nc, err := nats.Connect("nats://localhost:4222")
	if err != nil {
		t.Fatalf("Error connecting to nats: %v", err)
	}
	defer nc.Close()

	response, err := nc.Request("ping", []byte("ping"), 1000*time.Millisecond)
	if err != nil {
		t.Fatalf("error sending ping: %v", err)
	}

	if string(response.Data) != "pong" {
		t.Fatalf("wrong answer: %s", response.Data)
	}
}

func TestCreateKVStore(t *testing.T) {
	bucketNameHandle := "testKVStore"
	err := natsComponent.CreateKVStore(ConfigKVStore{
		CommonResourceConfig: CommonResourceConfig{
			Category: CategoryData,
			Prefixes: prefixes,
		},
		KVConfig: &nats.KeyValueConfig{
			Bucket: bucketNameHandle,
		},
	})
	if err != nil {
		t.Fatalf(err.Error())
	}

	// test full name is correct
	fullBucketName := natsComponent.MapKVStore[bucketNameHandle].fullName
	want := "data_TST_testnet_testKVStore"
	assert.Equal(t, fullBucketName, want)
	// check if the store exists in nats component
	if _, ok := natsComponent.MapKVStore[bucketNameHandle]; !ok {
		t.Fatalf("kv store '%s' no added in natscomponent", bucketNameHandle)
	}

	// check if the store exists in nats
	if _, err := natsComponent.JsContext.KeyValue(fullBucketName); err != nil {
		t.Fatalf("kv store '%s' was not created", fullBucketName)
	}

	// delete the store
	if err := natsComponent.DeleteKVStore(bucketNameHandle); err != nil {
		t.Fatalf(err.Error())
	}
}
