package nats

import (
	"fmt"
	"github.com/nats-io/nats.go"
)

func (c *ComponentNats) CreateKVStore(category ResourceCategory, config *nats.KeyValueConfig) (error, nats.KeyValue) {
	fullBucketName := c.GetFullBucketName(category, config.Bucket)
	store, err := c.JsContext.KeyValue(fullBucketName)
	if err != nil || store == nil {
		config.Bucket = fullBucketName
		store, err = c.JsContext.CreateKeyValue(config)
		if err != nil {
			return err, nil
		}
	}

	c.MapKVStore[config.Bucket] = store
	return nil, store
}

func (c *ComponentNats) GetFullBucketName(category ResourceCategory, bucket string) string {
	prefix := GetResourcePrefix(c.Config.ResourcePrefixes, category, UnderScore)
	return fmt.Sprintf("%s_%s", prefix, bucket)
}
