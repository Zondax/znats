package nats

import (
	"fmt"
	"github.com/nats-io/nats.go"
)

type ConfigObjectStore struct {
	CommonResourceConfig
	ObjectStoreConfig *nats.ObjectStoreConfig
}

type ObjectStoreNats struct {
	CommonResourceProperties
	Store nats.ObjectStore
}

func (c *ComponentNats) CreateObjectStore(config ConfigObjectStore) error {
	nameHandle := config.ObjectStoreConfig.Bucket
	fullBucketName := GetFullObjectStoreName(config)

	objectStore := ObjectStoreNats{
		CommonResourceProperties: CommonResourceProperties{
			NameHandle: nameHandle,
			Category:   config.Category,
			fullName:   fullBucketName,
		},
	}

	// check if bucket exists
	store, err := c.JsContext.ObjectStore(fullBucketName)
	if err == nil {
		objectStore.Store = store
		c.MapObjectStore[nameHandle] = objectStore
		return nil
	}

	// store doesn't exist, create it
	config.ObjectStoreConfig.Bucket = fullBucketName
	store, err = c.JsContext.CreateObjectStore(config.ObjectStoreConfig)
	if err != nil {
		return err
	}
	objectStore.Store = store
	c.MapObjectStore[nameHandle] = objectStore

	return nil
}

func GetFullObjectStoreName(config ConfigObjectStore) string {
	prefix := GetResourcePrefix(config.Prefixes, config.Category, UnderScore)
	return fmt.Sprintf("%s_%s", prefix, config.ObjectStoreConfig.Bucket)
}
