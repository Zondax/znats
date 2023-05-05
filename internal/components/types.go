package nats

import (
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

const EmptyString = ""

type CommonResourceProperties struct {
	// Category is the resource category to which this item belongs to
	Category ResourceCategory
	// NameHandle is the name of the resource
	NameHandle string
	// fullName is the full name of the resource used in the nats system
	fullName string
}

func (c *CommonResourceProperties) FullName() string {
	return c.fullName
}

type ComponentNats struct {
	// Config
	Config ConfigNats

	// Connections
	JsContext nats.JetStreamContext
	NatsConn  *nats.Conn

	// KV store map
	MapKVStore map[string]KVStoreNats

	// Object Store map
	MapObjectStore map[string]ObjectStoreNats

	// NatCLI
	natCLI map[string]func(*nats.Msg)

	// Streams
	Streams map[string]StreamNats

	// Topics map
	InputTopics  map[string]*Topic
	OutputTopics map[string]*Topic
}

type Credential struct {
	JWT  string
	Seed string
}

func (c Credential) isEmpty() bool {
	return c.JWT == EmptyString || c.Seed == EmptyString
}

const (
	Dot        = "."
	UnderScore = "_"
	Dash       = "-"
)

type ResourceCategory uint

const (
	NoCategory ResourceCategory = iota
	CategoryData
	CategorySystem
)

func (r ResourceCategory) String() string {
	switch r {
	case CategoryData:
		return "data"
	case CategorySystem:
		return "sys"
	case NoCategory:
		return ""
	default:
		zap.S().Errorf("Category '%d' unrecognized", r)
		return "unknown"
	}
}
