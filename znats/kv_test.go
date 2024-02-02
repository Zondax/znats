package znats

import (
	"github.com/nats-io/nats.go"
	"gotest.tools/assert"
	"testing"
)

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
	assert.Equal(t, 1, 1)
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
