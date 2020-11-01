package utils_test

import (
	"testing"

	utils "./"
	"github.com/stretchr/testify/assert"
)

func TestGetTimeNow(t *testing.T) {
	value := utils.GetTimeNow()
	assert.Equal(t, value, value)
}
