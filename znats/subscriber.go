package znats

import (
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
	"strings"
	"time"
)

type SubscriberNats struct {
	Topic        *Topic
	Subscription *nats.Subscription
}

func (c *ComponentNats) AddQueuedSubscriber(topicName string, cb nats.MsgHandler, opts ...nats.SubOpt) (error, SubscriberNats) {
	// Get input topic
	err, topic := c.GetInputTopic(topicName)
	if err != nil {
		zap.S().Errorf("failed to get input topic: %s", topicName)
		return err, SubscriberNats{}
	}

	queue := GetSubscriberQueueName(topic)
	consumer := GetSubscriberConsumerName(topic)
	topicFullName := topic.FullRoute()

	zap.S().Infof("Creating queued subscriber: \n Topic: %s \n Queue: %s \n Consumer: %s",
		topicFullName, queue, consumer)

	subscription, err := c.JsContext.QueueSubscribe(topicFullName, queue, cb, opts...)

	return err, SubscriberNats{
		Topic:        topic,
		Subscription: subscription,
	}
}

func (c *ComponentNats) AddQueuedPullSubscriber(topicName string, cb nats.MsgHandler, pullTime time.Duration, opts ...nats.SubOpt) (error, SubscriberNats) {
	// Get input topic
	err, topic := c.GetInputTopic(topicName)
	if err != nil {
		zap.S().Errorf("failed to get input topic: %s", topicName)
		return err, SubscriberNats{}
	}

	queue := GetSubscriberQueueName(topic)
	consumer := GetSubscriberConsumerName(topic)
	topicFullName := topic.FullRoute()

	zap.S().Infof("Creating queued pull subscriber: \n Topic: %s \n Queue: %s \n Consumer: %s",
		topicFullName, queue, consumer)

	subscription, err := c.JsContext.PullSubscribe(topicFullName, queue, opts...)
	if err != nil {
		return err, SubscriberNats{}
	}

	go func() {
		for {
			fetch, errFetch := subscription.Fetch(1)
			// nats returns ErrTimeout if no messages are available, so we skip printing that error
			if errFetch != nil && errFetch != nats.ErrTimeout {
				zap.S().Errorf("Error fetching messages: %s", errFetch.Error())
				time.Sleep(pullTime)
				continue
			}

			if len(fetch) == 0 {
				// b.LogDebug(fmt.Sprintf("No messages available for topic %s. Sleeping ...", topicFullName))
				time.Sleep(pullTime)
				continue
			}

			for _, msg := range fetch {
				cb(msg)
			}
		}
	}()

	return nil, SubscriberNats{
		Topic:        topic,
		Subscription: subscription,
	}
}

func (c *ComponentNats) AddSubscriber(topicName string, cb func(*nats.Msg), opts ...nats.SubOpt) (error, SubscriberNats) {
	// Get input topic
	err, topic := c.GetInputTopic(topicName)
	if err != nil {
		zap.S().Errorf("failed to get input topic: %s", topicName)
		return err, SubscriberNats{}
	}

	topicFullName := topic.FullRoute()
	zap.S().Infof("Creating subscriber: \n Topic: %s \n",
		topicFullName)

	subscription, err := c.JsContext.Subscribe(topicFullName, func(msg *nats.Msg) {
		cb(msg)
	}, opts...)

	return err, SubscriberNats{
		Topic:        topic,
		Subscription: subscription,
	}
}

func GetSubscriberConsumerName(topic *Topic) string {
	return strings.ReplaceAll(topic.FullRoute(), Dot, UnderScore)
}

func GetSubscriberQueueName(topic *Topic) string {
	return strings.ReplaceAll(topic.FullRoute(), Dot, UnderScore)
}
