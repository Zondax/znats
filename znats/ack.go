package znats

import (
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
	"time"
)

func AckWithLog(msg *nats.Msg, opts ...nats.AckOpt) {
	err := msg.Ack(opts...)
	if err != nil {
		zap.S().Errorf("[ERROR] on acknowledging message: %s", err.Error())
	}
}

func AckSyncWithLog(msg *nats.Msg, opts ...nats.AckOpt) {
	err := msg.AckSync(opts...)
	if err != nil {
		zap.S().Errorf("[ERROR] on acknowledging (sync) message: %s", err.Error())
	}
}

func NakWithLog(msg *nats.Msg, opts ...nats.AckOpt) {
	err := msg.Nak(opts...)
	if err != nil {
		zap.S().Errorf("[ERROR] on not acknowledging message: %s", err.Error())
	}
}

func NakDelayWithLog(msg *nats.Msg, delay time.Duration, opts ...nats.AckOpt) {
	err := msg.NakWithDelay(delay, opts...)
	if err != nil {
		zap.S().Errorf("[ERROR] on delayed not acknowledging message: %s", err.Error())
	}
}

func NakDelayWithRetry(msg *nats.Msg, delay time.Duration, maxRetries uint64, opts ...nats.AckOpt) {
	msgMetadata, err := msg.Metadata()
	if err != nil {
		AckWithLog(msg, opts...)
	}

	if msgMetadata.NumDelivered > maxRetries {
		// give up... do not redeliver this message anymore
		zap.S().Debugf("Giving up on message after %d retries: stream '%s', cunsumer '%s'",
			maxRetries, msgMetadata.Stream, msgMetadata.Consumer)
		AckWithLog(msg, opts...)
	} else {
		NakDelayWithLog(msg, delay, opts...)
	}
}
