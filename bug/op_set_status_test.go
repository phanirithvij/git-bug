package bug

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/MichaelMure/git-bug/entity"
	"github.com/stretchr/testify/require"

	"github.com/MichaelMure/git-bug/identity"
	"github.com/MichaelMure/git-bug/repository"
)

func TestSetStatusSerialize(t *testing.T) {
	repo := repository.NewMockRepo()

	rene, err := identity.NewIdentity(repo, "René Descartes", "rene@descartes.fr")
	require.NoError(t, err)

	unix := time.Now().Unix()
	before := NewSetStatusOp(rene, unix, ClosedStatus)

	data, err := json.Marshal(before)
	require.NoError(t, err)

	var after SetStatusOperation
	err = json.Unmarshal(data, &after)
	require.NoError(t, err)

	// enforce creating the ID
	before.Id()

	// Replace the identity as it's not serialized
	after.Author_ = rene

	require.Equal(t, before, &after)
}

func TestSetStatus(t *testing.T) {
	repo := repository.NewMockRepo()

	rene, err := identity.NewIdentity(repo, "René Descartes", "rene@descartes.fr")
	require.NoError(t, err)

	b, _, err := Create(rene, time.Now().Unix(), "title", "")
	require.NoError(t, err)

	op, err := Close(b, rene, time.Now().Unix())
	require.NoError(t, err)

	// SetStatus change the status and create a new timeline item
	snap := b.Compile()
	require.Equal(t, ClosedStatus, snap.Status)
	require.Equal(t, op, snap.Operations[1])
	// require.Equal(t, op.Id(), snap.Timeline[1].Id())
	require.Equal(t, entity.CombineIds(b.Id(), op.Id()), snap.Timeline[1].CombinedId())
	require.Equal(t, ClosedStatus, snap.Timeline[1].(*SetStatusTimelineItem).Status)

	op, err = Open(b, rene, time.Now().Unix())
	require.NoError(t, err)

	// SetStatus change the status and create a new timeline item
	snap = b.Compile()
	require.Equal(t, OpenStatus, snap.Status)
	require.Equal(t, op, snap.Operations[2])
	// require.Equal(t, op.Id(), snap.Timeline[2].Id())
	require.Equal(t, entity.CombineIds(b.Id(), op.Id()), snap.Timeline[2].CombinedId())
	require.Equal(t, OpenStatus, snap.Timeline[2].(*SetStatusTimelineItem).Status)
}
