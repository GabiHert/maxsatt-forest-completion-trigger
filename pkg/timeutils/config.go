package timeutils

import (
	"sync"
	"time"
)

var (
	instance *timeConfig
	once     sync.Once
)

type CurrentTimeProvider func() time.Time

type timeConfig struct {
	CurrentTimeProvider
}

func GetTimeConfig() *timeConfig {
	once.Do(func() {
		instance = &timeConfig{
			CurrentTimeProvider: time.Now,
		}
	})
	return instance
}

func (tc *timeConfig) Now() time.Time {
	return tc.CurrentTimeProvider()
}

func (tc *timeConfig) SetCurrentTimeProvider(provider CurrentTimeProvider) {
	tc.CurrentTimeProvider = provider
}
