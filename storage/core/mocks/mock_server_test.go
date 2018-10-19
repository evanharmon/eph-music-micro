package mocks_test

import (
	"errors"
	"testing"

	"github.com/evanharmon/eph-music-micro/storage/core/mocks"
	"github.com/golang/mock/gomock"
)

func TestListen(t *testing.T) {
	t.Run("Non-int port should return error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mock := mocks.NewMockServerService(ctrl)
		mock.EXPECT().Listen(
			gomock.Any(),
		).Return(nil, errors.New("invalid port"))
		testListenService(t, mock)
	})
}

func testListenService(t *testing.T, s *mocks.MockServerService) {
	if _, err := s.Listen(0); err == nil {
		t.Errorf("Non-int port should return error")
	}
}
