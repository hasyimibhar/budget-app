package budgeting

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCategory_Equal(t *testing.T) {
	assert := assert.New(t)

	c1 := NewCategory("food")
	c2 := NewCategory("bills")

	assert.True(c1.Equal(c1))
	assert.True(c2.Equal(c2))
	assert.False(c1.Equal(c2))
	assert.False(c1.Equal(nil))
	assert.False(c2.Equal(nil))

	c2.Name = "food"
	assert.False(c1.Equal(c2))

	c3 := c2
	c3.Name = "entertainment"
	assert.True(c2.Equal(c3))
	assert.EqualValues("entertainment", c2.Name)
}
