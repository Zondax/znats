package nats

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

type ConfigKVStore struct {
	CommonResourceConfig
	KVConfig *nats.KeyValueConfig
}

type KVStoreNats struct {
	CommonResourceProperties
	Store nats.KeyValue
}

func (cfg ConfigKVStore) Validate() bool {
	if cfg.KVConfig == nil {
		zap.S().Errorf("kv config is nil")
		return false
	}

	if cfg.KVConfig.Bucket == EmptyString {
		zap.S().Errorf("kv config bucket is empty")
		return false
	}

	return true
}

func (c *ComponentNats) CreateKVStore(config ConfigKVStore) error {
	if !config.Validate() {
		return fmt.Errorf("invalid kv store config")
	}

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

func (c *ComponentNats) DeleteKVStore(nameHandle string) error {
	// check if the store exists in nats component
	if store, ok := c.MapKVStore[nameHandle]; ok {

		// delete from nats
		if err := c.JsContext.DeleteKeyValue(store.FullName()); err != nil {
			return err
		}

		// delete from map
		delete(c.MapKVStore, nameHandle)
	} else {
		return fmt.Errorf("kv store '%s' does not exist in nats component", nameHandle)
	}

	return nil
}

func GetFullBucketName(config ConfigKVStore) string {
	prefix := GetResourcePrefix(config.Prefixes, config.Category, UnderScore)
	return fmt.Sprintf("%s%s", prefix, config.KVConfig.Bucket)
}
