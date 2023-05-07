package znats

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

type ReqReplyCB struct {
	Callback func(*nats.Msg)
	Global   bool
}

func (c *ComponentNats) AddNatCliCmd(natcli map[string]ReqReplyCB) {
	for topic, cb := range natcli {
		reqName := c.GetReqReplyFullName(cb.Global, topic)
		if _, ok := c.natCLI[reqName]; ok {
			continue
		}
		prefix := GetResourcePrefix(c.Config.ResourcePrefixes, NoCategory, Dash)
		queueName := fmt.Sprintf("%s-%s", prefix, topic)
		c.natCLI[reqName] = cb.Callback
		_, err := c.NatsConn.QueueSubscribe(reqName, queueName, cb.Callback)
		if err != nil {
			zap.S().Error(err.Error())
		}
		zap.S().Infof("added responder for request '%s'", reqName)
	}
}

func (c *ComponentNats) GetReqReplyFullName(global bool, reqReplyName string) string {
	prefix := GetResourcePrefix(c.Config.ResourcePrefixes, NoCategory, Dash)
	if !global {
		return fmt.Sprintf("%s%s-%s", prefix, c.Config.ServiceName, reqReplyName)
	} else {
		return fmt.Sprintf("%s%s", prefix, reqReplyName)
	}
}

func (c *ComponentNats) ReplyListAvailableCliCmd(req *nats.Msg) {
	res := make(map[string][]string)
	l := make([]string, 0)
	for cmd := range c.natCLI {
		l = append(l, cmd)
	}
	res[c.Config.ServiceName] = l
	resJson, _ := json.Marshal(res)
	_ = req.Respond(resJson)
}
