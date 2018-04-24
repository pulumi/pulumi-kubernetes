package await

import (
	"time"

	"github.com/golang/glog"
)

// DefaultTimeout creates timeout. Argument could be time.Duration, int64 or float64.
func DefaultTimeout(tx interface{}) *time.Duration {
	var td time.Duration
	switch raw := tx.(type) {
	case time.Duration:
		return &raw
	case int64:
		td = time.Duration(raw)
	case float64:
		td = time.Duration(int64(raw))
	default:
		glog.V(3).Infof("Unknown type in DefaultTimeout: %#v", tx)
	}
	return &td
}
