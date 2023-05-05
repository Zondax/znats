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

	// check if bucket exists
	store, err := c.JsContext.ObjectStore(fullBucketName)
	if err != nil || store == nil {
		config.ObjectStoreConfig.Bucket = fullBucketName
		store, err = c.JsContext.CreateObjectStore(config.ObjectStoreConfig)
		if err != nil {
			return err
		}
		return nil
	}

	objectStore := ObjectStoreNats{
		CommonResourceProperties: CommonResourceProperties{
			NameHandle: nameHandle,
			Category:   config.Category,
			fullName:   fullBucketName,
		},
		Store: store,
	}

	c.MapObjectStore[nameHandle] = objectStore
	return nil
}

func GetFullObjectStoreName(config ConfigObjectStore) string {
	prefix := GetResourcePrefix(config.Prefixes, config.Category, UnderScore)
	return fmt.Sprintf("%s_%s", prefix, config.ObjectStoreConfig.Bucket)
}
