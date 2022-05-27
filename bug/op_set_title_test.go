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

func TestSetTitleSerialize(t *testing.T) {
	repo := repository.NewMockRepo()

	rene, err := identity.NewIdentity(repo, "René Descartes", "rene@descartes.fr")
	require.NoError(t, err)

	unix := time.Now().Unix()
	before := NewSetTitleOp(rene, unix, "title", "was")

	data, err := json.Marshal(before)
	require.NoError(t, err)

	var after SetTitleOperation
	err = json.Unmarshal(data, &after)
	require.NoError(t, err)

	// enforce creating the ID
	before.Id()

	// Replace the identity as it's not serialized
	after.Author_ = rene

	require.Equal(t, before, &after)
}

func TestSetTitle(t *testing.T) {
	repo := repository.NewMockRepo()

	rene, err := identity.NewIdentity(repo, "René Descartes", "rene@descartes.fr")
	require.NoError(t, err)

	b, _, err := Create(rene, time.Now().Unix(), "original title", "")
	require.NoError(t, err)

	op, err := SetTitle(b, rene, time.Now().Unix(), "title1")
	require.NoError(t, err)

	// SetTitle change the title and create a new timeline item
	snap := b.Compile()
	require.Equal(t, "title1", snap.Title)
	require.Equal(t, op, snap.Operations[1])
	// require.Equal(t, op.Id(), snap.Timeline[1].Id())
	require.Equal(t, entity.CombineIds(b.Id(), op.Id()), snap.Timeline[1].CombinedId())
	require.Equal(t, "title1", snap.Timeline[1].(*SetTitleTimelineItem).Title)
	require.Equal(t, "original title", snap.Timeline[1].(*SetTitleTimelineItem).Was)

	op, err = SetTitle(b, rene, time.Now().Unix(), "title2")
	require.NoError(t, err)

	// SetTitle change the title and create a new timeline item
	snap = b.Compile()
	require.Equal(t, "title2", b.Compile().Title)
	require.Equal(t, op, snap.Operations[2])
	// require.Equal(t, op.Id(), snap.Timeline[2].Id())
	require.Equal(t, entity.CombineIds(b.Id(), op.Id()), snap.Timeline[2].CombinedId())
	require.Equal(t, "title2", snap.Timeline[2].(*SetTitleTimelineItem).Title)
	require.Equal(t, "title1", snap.Timeline[2].(*SetTitleTimelineItem).Was)
}
