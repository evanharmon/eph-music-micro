package core_test

import (
	"testing"

	"github.com/evanharmon/eph-music-micro/storage/core"
)

func TestNewProviderGRPC(t *testing.T) {
	t.Run("Non-int port should return error", func(t *testing.T) {
		cfg := core.ProviderGRPCConfig{Port: 0}
		if _, err := core.NewProviderGRPC(cfg); err == nil {
			t.Errorf("Invalid Port 0 should throw error")
		}
	})

}
