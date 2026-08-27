package clock

import "time"

type Clock interface{ Now() time.Time }

type System struct{}

func (System) Now() time.Time { return time.Now().UTC() }

type Fixed struct{ Value time.Time }

func (f Fixed) Now() time.Time { return f.Value.UTC() }

func Today(c Clock) time.Time {
	now := c.Now()
	y, m, d := now.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}
