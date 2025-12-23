package mock

import "time"

type Time struct {
	currentStartTime time.Time
	updatedAt        time.Time
}

func NewTime() *Time {
	return &Time{
		currentStartTime: time.Now().UTC(),
		updatedAt:        time.Now().UTC(),
	}
}

func (t *Time) SetCurrentTime(currentTime time.Time) {
	t.currentStartTime = currentTime.UTC()
	t.updatedAt = time.Now().UTC()
}

func (t *Time) Now() time.Time {
	elapsed := t.updatedAt.UTC().Sub(time.Now().UTC())
	return t.currentStartTime.UTC().Add(elapsed)
}
