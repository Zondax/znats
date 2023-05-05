package nats

import (
	"gotest.tools/assert"
	"testing"
)

func TestBuildTopic(t *testing.T) {
	topic := NewTopic(&TopicConfig{
		CommonResourceConfig: CommonResourceConfig{
			Category: CategoryData,
			Prefixes: []string{"pfxA", "pfxB", "pfxC"},
		},
		Name:    "test",
		Subject: "subject",
	})

	want := "data.pfxA.pfxB.pfxC.subject"
	assert.Equal(t, topic.FullRoute(), want)

	topic = NewTopic(&TopicConfig{
		CommonResourceConfig: CommonResourceConfig{
			Category: CategoryData,
			Prefixes: []string{"pfxA", "pfxB", "pfxC"},
		},
		Name:       "test",
		Subject:    "subject",
		SubSubject: "subsubject",
	})

	want = "data.pfxA.pfxB.pfxC.subject.subsubject"
	assert.Equal(t, topic.FullRoute(), want)

	topic = NewTopic(&TopicConfig{
		CommonResourceConfig: CommonResourceConfig{
			Category: NoCategory,
			Prefixes: []string{"pfxA", "pfxB", "pfxC"},
		},
		Name:       "test",
		Subject:    "subject",
		SubSubject: "subsubject",
	})

	want = "pfxA.pfxB.pfxC.subject.subsubject"
	assert.Equal(t, topic.FullRoute(), want)
}
