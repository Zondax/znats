package nats

import (
	"github.com/nats-io/nats.go"
)

func (c *ComponentNats) CreateObjectStore(category ResourceCategory, config *nats.ObjectStoreConfig) (error, nats.ObjectStore) {
	fullBucketName := c.GetFullBucketName(category, config.Bucket)
	// check if bucket exists
	store, err := c.JsContext.ObjectStore(fullBucketName)
	if err != nil || store == nil {
		config.Bucket = fullBucketName
		store, err = c.JsContext.CreateObjectStore(config)
		if err != nil {
			return err, nil
		}
		return nil, store
	}

	c.MapObjectStore[config.Bucket] = store
	return nil, store
}
