package clock

import "time"

type Clock interface {
	Now() time.Time
}

var _ Clock = (*clockImpl)(nil)

type clockImpl struct{}

func New() Clock {
	return &clockImpl{}
}

func (c *clockImpl) Now() time.Time {
	return time.Now()
}

type ClockMock struct {
	data []time.Time
}

var _ Clock = (*ClockMock)(nil)

func NewMock(data []time.Time) Clock {
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
