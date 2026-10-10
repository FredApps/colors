package space_test

import (
	"github.com/thmalt/colors/space"
	"testing"
)

func TestInfoReturnsNilForEveryInvalidSpace(t *testing.T) {
	for i := 0; i < 256; i++ {
		s := space.Space(i)
		info := s.Info()
		if !s.IsValid() {
			if info != nil {
				t.Errorf("invalid space %d Info() = %T; want nil", i, info)
			}
			continue
		}
		if info == nil {
			t.Errorf("valid space %d has no metadata", i)
			continue
		}
		if info.Name() == "" || info.ChannelCount() != s.ChannelCount() {
			t.Errorf("valid space %d has inconsistent metadata", i)
		}
	}
}
