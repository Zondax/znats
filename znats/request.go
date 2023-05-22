package znats

import (
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
	"time"
)

func (c *ComponentNats) SendRequest(topic *Topic, data []byte, reqRetry int, reqTimeout time.Duration, waitInterval time.Duration) (error, []byte) {
	var out []byte
	for i := 0; i < reqRetry; i++ {
		response, err := c.NatsConn.Request(topic.FullRoute(), data, reqTimeout)

		if err != nil {
			if err == nats.ErrTimeout {
				zap.S().Errorf("Request to topic '%s' timeout, retrying...", topic.FullRoute())
				time.Sleep(waitInterval)
				continue
			} else {
				return err, nil
			}
		}

		out = response.Data
		break
	}

	return nil, out
}

func (c *ComponentNats) WaitTopicAndSendRequest(topic *Topic, data []byte, reqRetry int, reqTimeout time.Duration, sleepInterval time.Duration) (error, []byte) {
	c.WaitForTopicToExist(topic)
	return c.SendRequest(topic, data, reqRetry, reqTimeout, sleepInterval)
}
