package service

import (
	"github.com/zondax/znats/internal/conf"
	"go.uber.org/zap"
	"time"
)

func Start(config *conf.Config) {
	// TODO implement me
	for {
		zap.L().Info("Implement me!")
		time.Sleep(5 * time.Second)
	}
}
