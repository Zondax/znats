package znats

import (
	"fmt"
	"go.uber.org/zap"
	"time"
)

type TopicConfig struct {
	CommonResourceConfig
	// Name is the name human friendly name
	Name string
	// Subject is the subject of the topic
	Subject string
	// SubSubject is the sub-subject of the topic (optional)
	SubSubject string
}

type Topic struct {
	CommonResourceProperties
}

func NewTopic(config *TopicConfig) *Topic {
	t := &Topic{
		CommonResourceProperties{
			Category:   config.Category,
			NameHandle: config.Name,
		},
	}

	t.build(config)
	return t
}

func (t *Topic) FullRoute() string {
	return t.CommonResourceProperties.FullName()
}

func (t *Topic) Name() string {
	return t.CommonResourceProperties.NameHandle
}

func (c *ComponentNats) WaitForTopicToExist(topic *Topic) {
	connected := false
	for !connected {
		stream, err := c.JsContext.StreamNameBySubject(topic.FullRoute())
		if err != nil || stream == "" {
			zap.S().Infof("waiting for topic to be created: '%s' ...", topic.FullRoute())
			time.Sleep(5 * time.Second)
			continue
		}
		connected = true
	}
}

func (c *ComponentNats) AddInputTopic(topic *Topic) {
	if topic == nil {
		zap.S().Error("trying to add nil input topic")
		return
	}

	c.InputTopics[topic.CommonResourceProperties.NameHandle] = topic
}

func (c *ComponentNats) AddOutputTopic(topic *Topic) {
	if topic == nil {
		zap.S().Error("trying to add nil output topic")
		return
	}

	c.OutputTopics[topic.CommonResourceProperties.NameHandle] = topic
}

func (c *ComponentNats) GetInputTopic(name string) (error, *Topic) {
	if topic, ok := c.InputTopics[name]; ok {
		return nil, topic
	} else {
		err := fmt.Errorf("topic '%s' not found", name)
		zap.S().Errorf(err.Error())
		return err, nil
	}
}

func (c *ComponentNats) GetInputTopicsFullNames() []string {
	topics := make([]string, 0, len(c.InputTopics))
	for _, topic := range c.InputTopics {
		topics = append(topics, topic.FullRoute())
	}
	return topics
}

func (c *ComponentNats) GetOutputTopic(name string) (error, *Topic) {
	if topic, ok := c.OutputTopics[name]; ok {
		return nil, topic
	} else {
		err := fmt.Errorf("topic '%s' not found", name)
		zap.S().Errorf(err.Error())
		return err, nil
	}
}

func (c *ComponentNats) GetOutputTopicsFullNames() []string {
	topics := make([]string, 0, len(c.OutputTopics))
	for _, topic := range c.OutputTopics {
		topics = append(topics, topic.FullRoute())
	}
	return topics
}

// build builds the full route for the topic.
// Format: <category>.<prefix_1>...<prefix_n>.<subject>.<subsubject>
func (t *Topic) build(config *TopicConfig) {
	p := GetResourcePrefix(config.Prefixes, config.Category, Dot)
	t.CommonResourceProperties.fullName = fmt.Sprintf("%s%s", p, config.Subject)
	if config.SubSubject != "" {
		t.CommonResourceProperties.fullName += fmt.Sprintf(".%s", config.SubSubject)
	}

	// Set default name if not set
	if t.CommonResourceProperties.NameHandle == "" {
		t.CommonResourceProperties.NameHandle = t.CommonResourceProperties.fullName
	}
}
