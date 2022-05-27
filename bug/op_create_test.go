package bug

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/MichaelMure/git-bug/entity"
	"github.com/MichaelMure/git-bug/identity"
	"github.com/MichaelMure/git-bug/repository"
)

func TestCreate(t *testing.T) {
	repo := repository.NewMockRepo()

	rene, err := identity.NewIdentity(repo, "René Descartes", "rene@descartes.fr")
	require.NoError(t, err)

	b, op, err := Create(rene, time.Now().Unix(), "title", "message")
	require.NoError(t, err)

	require.Equal(t, "title", op.Title)
	require.Equal(t, "message", op.Message)

	// Create generate the initial operation and create a new timeline item
	snap := b.Compile()
	require.Equal(t, OpenStatus, snap.Status)
	require.Equal(t, rene, snap.Author)
	require.Equal(t, "title", snap.Title)
	require.Len(t, snap.Operations, 1)
	require.Equal(t, op, snap.Operations[0])

	require.Len(t, snap.Timeline, 1)
	require.Equal(t, entity.CombineIds(b.Id(), op.Id()), snap.Timeline[0].CombinedId())
	require.Equal(t, rene, snap.Timeline[0].(*CreateTimelineItem).Author)
	require.Equal(t, "message", snap.Timeline[0].(*CreateTimelineItem).Message)
}

func TestCreateSerialize(t *testing.T) {
	repo := repository.NewMockRepo()

	rene, err := identity.NewIdentity(repo, "René Descartes", "rene@descartes.fr")
	require.NoError(t, err)

	unix := time.Now().Unix()
	before := NewCreateOp(rene, unix, "title", "message", nil)

	data, err := json.Marshal(before)
	require.NoError(t, err)

	var after CreateOperation
	err = json.Unmarshal(data, &after)
	require.NoError(t, err)

	// enforce creating the ID
	before.Id()

	// Replace the identity as it's not serialized
	after.Author_ = rene

	require.Equal(t, before, &after)
}
