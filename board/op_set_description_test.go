package board

import (
	"testing"

	"github.com/MichaelMure/git-bug/entity/dag"
	"github.com/MichaelMure/git-bug/identity"
)

func TestSetDescriptionSerialize(t *testing.T) {
	dag.SerializeRoundTripTest(t, func(author identity.Interface, unixTime int64) *SetDescriptionOperation {
		return NewSetDescriptionOp(author, unixTime, "description", "was")
	})
}
