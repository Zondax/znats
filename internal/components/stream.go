package nats

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

const (
	StreamMain = "main"
)

type StreamNats struct {
	CommonResourceProperties
	Info *nats.StreamInfo
}

type ConfigStream struct {
	CommonResourceConfig
	NatsStreamConfig *nats.StreamConfig
}

func (c *ComponentNats) CreateStream(config ConfigStream) error {
	// Append prefix to stream name and subjects
	nameHandle := config.NatsStreamConfig.Name
	config.NatsStreamConfig.Name = GetStreamFullName(config)
	if len(config.NatsStreamConfig.Subjects) == 0 {
		// Add default subject
		config.NatsStreamConfig.Subjects = []string{">"}
	}

	fixedSubjects := make([]string, 0, len(config.NatsStreamConfig.Subjects))
	for _, subject := range config.NatsStreamConfig.Subjects {
		// create a Topic for each subject
		prefixes := config.Prefixes
		prefixes = append(prefixes, nameHandle)
		t := NewTopic(&TopicConfig{
			CommonResourceConfig: CommonResourceConfig{
				Category: config.Category,
				Prefixes: prefixes,
			},
			Subject: subject,
		})

		fixedSubjects = append(fixedSubjects, t.FullRoute())
	}
	config.NatsStreamConfig.Subjects = fixedSubjects

	// Check if stream already exists
	if exists := c.StreamExists(config.NatsStreamConfig.Name); !exists {
		streamInfo, err := c.JsContext.AddStream(config.NatsStreamConfig)
		if err != nil {
			zap.S().Errorf("could not create stream '%s': %s", config.NatsStreamConfig.Name, err.Error())
			return err
		}

		// Add the stream to the streams map
		c.Streams[nameHandle] = StreamNats{
			CommonResourceProperties: CommonResourceProperties{
				NameHandle: nameHandle,
				Category:   config.Category,
				fullName:   config.NatsStreamConfig.Name,
			},
			Info: streamInfo,
		}
		zap.S().Infof("created stream '%s' with subjects %v", config.NatsStreamConfig.Name, config.NatsStreamConfig.Subjects)
	}

	return nil
}

func (c *ComponentNats) CreateStreams(configs []ConfigStream) error {
	for _, config := range configs {
		err := c.CreateStream(config)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *ComponentNats) StreamExists(fullName string) bool {
	info, err := c.JsContext.StreamInfo(fullName)
	if err != nil {
		return false
	}

	return info != nil
}

func (c *ComponentNats) PurgeStream(fullName string) error {
	if ok := c.StreamExists(fullName); ok {
		err := c.JsContext.PurgeStream(fullName)
		if err != nil {
			zap.S().Errorf("could not purge stream '%s': %v", fullName, err)
			return err
		}
		zap.S().Infof("purged stream '%s'", fullName)
	}
	return nil
}

func (c *ComponentNats) DeleteStream(nameHandle string) error {
	// delete from nats server
	if stream, ok := c.Streams[nameHandle]; ok {
		fullStreamName := stream.Info.Config.Name
		if ok := c.StreamExists(fullStreamName); ok {
			err := c.JsContext.DeleteStream(fullStreamName)
			if err != nil {
				zap.S().Errorf("could not delete stream '%s': %v", fullStreamName, err)
				return err
			}

			delete(c.Streams, nameHandle)
			zap.S().Infof("deleted stream '%s' with full name '%s'", nameHandle, fullStreamName)
		}
	}

	return nil
}

func GetStreamFullName(config ConfigStream) string {
	prefix := GetResourcePrefix(config.Prefixes, config.Category, UnderScore)
	return fmt.Sprintf("%s%s", prefix, config.NatsStreamConfig.Name)
}
