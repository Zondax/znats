package nats

import (
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
	"time"
)

const RequestRetries = 5

func (c *ComponentNats) SendRequest(topic *Topic, data []byte, timeout time.Duration) (error, []byte) {
	var out []byte
	for i := 0; i < RequestRetries; i++ {
		response, err := c.NatsConn.Request(topic.FullRoute(), data, timeout)

		if err != nil {
			if err == nats.ErrTimeout {
				zap.S().Errorf("Request to topic '%s' timeout, retrying...", topic.FullRoute())
				time.Sleep(5 * time.Second)
				continue
			} else {
				return err, nil
			}
		}

		out = response.Data
	}

	return nil, out
}

func (c *ComponentNats) WaitTopicAndSendRequest(topic *Topic, data []byte, timeout time.Duration) (error, []byte) {
	c.WaitForTopicToExist(topic)
	return c.SendRequest(topic, data, timeout)
}
