package timeutils

import (
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func TimeToNs(t time.Time) int64 {
	if t.IsZero() {
		return 0
	}
	return t.UnixNano()
}

func ProtobufTimestampToUnix(timestamp *timestamppb.Timestamp) int64 {
	if timestamp == nil {
		return 0
	}
	return TimeToNs(timestamp.AsTime())
}

func ProtobufTimestampToTime(timestamp *timestamppb.Timestamp) time.Time {
	if timestamp == nil {
		return time.Time{}
	}
	return timestamp.AsTime()
}
