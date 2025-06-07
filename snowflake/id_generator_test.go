package snowflake

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestId(t *testing.T) {
	id1 := Id()
	id2 := Id()

	assert.Truef(t, id2 > id1, "generate random incremented values")
}

func TestIdString(t *testing.T) {
	id1 := IdString()
	id2 := IdString()

	assert.Truef(t, id2 > id1, "generate random incremented values")
}
