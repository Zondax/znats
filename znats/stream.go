package znats

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
	"sort"
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

func (cfg ConfigStream) Validate() bool {
	if cfg.NatsStreamConfig == nil {
		zap.S().Errorf("stream config is nil")
		return false
	}

	if cfg.NatsStreamConfig.Name == EmptyString {
		zap.S().Errorf("stream config name is empty")
		return false
	}

	return true
}

func (c *ComponentNats) CreateStream(config ConfigStream) error {
	if !config.Validate() {
		return fmt.Errorf("invalid stream config")
	}

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
		// stream does not exist, create it
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
	} else {
		// stream exists, add it to the streams map
		streamInfo, err := c.JsContext.StreamInfo(config.NatsStreamConfig.Name)
		if err != nil {
			zap.S().Errorf("could not get stream info for '%s': %s", config.NatsStreamConfig.Name, err.Error())
			return err
		}

		c.Streams[nameHandle] = StreamNats{
			CommonResourceProperties: CommonResourceProperties{
				NameHandle: nameHandle,
				Category:   config.Category,
				fullName:   config.NatsStreamConfig.Name,
			},
			Info: streamInfo,
		}
		zap.S().Infof("added already existing stream '%s'", config.NatsStreamConfig.Name)
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

func (c *ComponentNats) GetStreamsFullNames() []string {
	streams := make([]string, 0, len(c.Streams))
	for _, stream := range c.Streams {
		streams = append(streams, stream.FullName())
	}

	sort.Strings(streams)
	return streams
}

func GetStreamFullName(config ConfigStream) string {
	prefix := GetResourcePrefix(config.Prefixes, config.Category, UnderScore)
	return fmt.Sprintf("%s%s", prefix, config.NatsStreamConfig.Name)
}
