package match

import (
	"fmt"

	"github.com/golang/mock/gomock"
)

// sliceSize matcher for gomock that checks size of slice.
type sliceSize struct {
	size int
}

// SliceSize returns matcher that checks size of slice.
func SliceSize(size int) gomock.Matcher {
	return &sliceSize{size: size}
}

func (m *sliceSize) Matches(x interface{}) bool {
	s, ok := x.([]byte)
	if !ok {
		return false
	}
	return len(s) == m.size
}

func (m *sliceSize) String() string {
	return fmt.Sprintf("size is %d", m.size)
}
