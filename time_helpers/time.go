package time_helpers

import "time"

func ToStartOfInterval(t time.Time, ivl time.Duration) time.Time {
	return t.Add(time.Duration(-t.UnixNano() % ivl.Nanoseconds()))
}

func ToEndOfInterval(t time.Time, ivl time.Duration) time.Time {
	return ToStartOfInterval(t, ivl).Add(ivl)
}
