package common

import "time"

func ToStartOfInterval(t time.Time, ivl time.Duration) time.Time {
	return t.Add(-(time.Duration(t.Nanosecond()) % time.Duration(ivl.Nanoseconds())))
}

func ToEndOfInterval(t time.Time, ivl time.Duration) time.Time {
	return t.
		Add(-(time.Duration(t.Nanosecond()) % time.Duration(ivl.Nanoseconds()))).
		Add(ivl)
}
