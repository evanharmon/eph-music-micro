package core_test

import (
	"testing"

	"github.com/evanharmon/eph-music-micro/storage/core"
)

func TestNewServerGRPC(t *testing.T) {
	t.Run("Non-int port should return error", func(t *testing.T) {
		cfg := core.ServerGRPCConfig{Port: 0}
		if _, err := core.NewServerGRPC(cfg); err == nil {
			t.Errorf("Invalid Port 0 should throw error")
		}
	})

}
