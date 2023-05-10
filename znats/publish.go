package znats

import (
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

func (c *ComponentNats) PublishAsync(topicName string, msg []byte, opts ...nats.PubOpt) (error, nats.PubAckFuture) {
	// Get output topic
	err, topic := c.GetOutputTopic(topicName)
	if err != nil {
		zap.S().Errorf("failed to get output topic: %s", topicName)
		return err, nil
	}

	ackFuture, err := c.JsContext.PublishAsync(topic.FullRoute(), msg, opts...)
	if err != nil {
		zap.S().Errorf("error on PublishAsync for topic '%s': %s", topic.FullRoute(), err.Error())
		return err, nil
	}
	return nil, ackFuture
}

func (c *ComponentNats) Publish(topicName string, msg []byte, opts ...nats.PubOpt) (error, *nats.PubAck) {
	// Get output topic
	err, topic := c.GetOutputTopic(topicName)
	if err != nil {
		zap.S().Errorf("failed to get output topic: %s", topicName)
		return err, nil
	}

	pub, err := c.JsContext.Publish(topic.FullRoute(), msg, opts...)
	if err != nil {
		zap.S().Errorf("error on Publish for topic '%s': %s", topic.FullRoute(), err.Error())
		return err, nil
	}
	return nil, pub
}

func (c *ComponentNats) PublishMsg(msg *nats.Msg, opts ...nats.PubOpt) (error, *nats.PubAck) {
	pub, err := c.JsContext.PublishMsg(msg, opts...)
	if err != nil {
		zap.S().Errorf("error on PublishMsg for topic '%s': %s", msg.Subject, err)
		return err, nil
	}
	return nil, pub
}
