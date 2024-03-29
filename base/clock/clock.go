package clock

import "time"

type Clock interface {
	Now() time.Time
}

var _ Clock = (*ClockImpl)(nil)

type ClockImpl struct{}

func New() *ClockImpl {
	return &ClockImpl{}
}

func (c *ClockImpl) Now() time.Time {
	return time.Now()
}

type ClockMock struct {
	data []time.Time
}

var _ Clock = (*ClockMock)(nil)

func NewMock(data []time.Time) *ClockMock {
	return &ClockMock{data: data}
}

func (c *ClockMock) Now() time.Time {
	if len(c.data) == 0 {
		panic("clock mock data is empty")
	}

	now := c.data[0]
	c.data = c.data[1:]
	return now
}
