package zlog

import (
	"Kama-Chat/initialize/zlog"
	"go.uber.org/zap"
	"testing"
)

func TestInfo(t *testing.T) {
	zlog.Info("this is a info", zap.String("name", "apylee"))
}
