package znats

import (
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

func (c *ComponentNats) PublishAsync(topicName string, msg []byte, opts ...nats.PubOpt) (nats.PubAckFuture, error) {
	// Get output topic
	topic, err := c.GetOutputTopic(topicName)
	if err != nil {
		zap.S().Errorf("failed to get output topic: %s", topicName)
		return nil, err
	}

	ackFuture, err := c.JsContext.PublishAsync(topic.FullRoute(), msg, opts...)
	if err != nil {
		zap.S().Errorf("error on PublishAsync for topic '%s': %s", topic.FullRoute(), err.Error())
		return nil, err
	}
	return ackFuture, nil
}

func (c *ComponentNats) Publish(topicName string, msg []byte, opts ...nats.PubOpt) (*nats.PubAck, error) {
	// Get output topic
	topic, err := c.GetOutputTopic(topicName)
	if err != nil {
		zap.S().Errorf("failed to get output topic: %s", topicName)
		return nil, err
	}

	pub, err := c.JsContext.Publish(topic.FullRoute(), msg, opts...)
	if err != nil {
		zap.S().Errorf("error on Publish for topic '%s': %s", topic.FullRoute(), err.Error())
		return nil, err
	}
	return pub, nil
}

func (c *ComponentNats) PublishMsg(msg *nats.Msg, opts ...nats.PubOpt) (*nats.PubAck, error) {
	pub, err := c.JsContext.PublishMsg(msg, opts...)
	if err != nil {
		zap.S().Errorf("error on PublishMsg for topic '%s': %s", msg.Subject, err)
		return nil, err
	}
	return pub, nil
}
