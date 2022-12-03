package model

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestValidatePassword(t *testing.T) {
	testTable := []struct {
		CronEvery   int64
		LastCronRun time.Time
		res         bool
	}{
		{10, time.Now(), false},
		{10, time.Now().Add(-time.Minute * 10), true},
		{1, time.Now(), false},
		{1, time.Now().Add(-time.Minute * 1), true},
		{1440, time.Now().Add(-time.Hour * 24), true},
		{1400, time.Now().Add(-time.Hour * 24), false},
	}
	t.Parallel()
	for _, test := range testTable {
		bot := SlackBot{CronEvery: test.CronEvery, LastCronRun: &test.LastCronRun}
		res := bot.IsRun()

		assert.Equal(t, test.res, res, fmt.Sprintf("Incorect return result expect : %v, result: %v, testData : %v", test.res, res, test))
	}
}
