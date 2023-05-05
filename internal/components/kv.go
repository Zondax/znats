package nats

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

type ConfigKeyValueStore struct {
	CommonResourceConfig
	KVConfig *nats.KeyValueConfig
}

type KVStoreNats struct {
	CommonResourceProperties
	Store nats.KeyValue
}

func (c *ComponentNats) CreateKVStore(config ConfigKeyValueStore) error {
	nameHandle := config.KVConfig.Bucket
	fullBucketName := GetFullBucketName(config)

	// check if bucket exists
	store, err := c.JsContext.KeyValue(fullBucketName)
	if err != nil || store == nil {
		config.KVConfig.Bucket = fullBucketName
		store, err = c.JsContext.CreateKeyValue(config.KVConfig)
		if err != nil {
			zap.S().Errorf("could not create kv store '%s': %s", fullBucketName, err.Error())
			return err
		}
	} else {
		zap.S().Debugf("kv store '%s' already exists", fullBucketName)
	}

	kvStore := KVStoreNats{
		CommonResourceProperties: CommonResourceProperties{
			NameHandle: nameHandle,
			Category:   config.Category,
			fullName:   fullBucketName,
		},
		Store: store,
	}

	// Add to map
	c.MapKVStore[kvStore.NameHandle] = kvStore

	return nil
}

func GetFullBucketName(config ConfigKeyValueStore) string {
	prefix := GetResourcePrefix(config.Prefixes, config.Category, UnderScore)
	return fmt.Sprintf("%s_%s", prefix, config.KVConfig.Bucket)
}
