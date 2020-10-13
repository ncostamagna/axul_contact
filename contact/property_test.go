package emails

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmailConstants(t *testing.T) {

	assert.EqualValues(t, toPty, "TO")
	assert.EqualValues(t, ccPty, "CC")
	assert.EqualValues(t, bccPty, "BCC")
	assert.EqualValues(t, fromPty, "FROM")
	assert.EqualValues(t, replyToPty, "REPLY_TO")
}
