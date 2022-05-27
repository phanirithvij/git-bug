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

func TestLabelChangeSerialize(t *testing.T) {
	repo := repository.NewMockRepo()

	rene, err := identity.NewIdentity(repo, "René Descartes", "rene@descartes.fr")
	require.NoError(t, err)

	unix := time.Now().Unix()
	before := NewLabelChangeOperation(rene, unix, []Label{"added"}, []Label{"removed"})

	data, err := json.Marshal(before)
	require.NoError(t, err)

	var after LabelChangeOperation
	err = json.Unmarshal(data, &after)
	require.NoError(t, err)

	// enforce creating the ID
	before.Id()

	// Replace the identity as it's not serialized
	after.Author_ = rene

	require.Equal(t, before, &after)
}

func TestLabelChange(t *testing.T) {
	repo := repository.NewMockRepo()

	rene, err := identity.NewIdentity(repo, "René Descartes", "rene@descartes.fr")
	require.NoError(t, err)

	b, _, err := Create(rene, time.Now().Unix(), "title", "")
	require.NoError(t, err)

	// LabelChange add/remove labels and create a new timeline item
	_, op, err := ChangeLabels(b, rene, time.Now().Unix(), []string{"foo", "bar"}, []string{})
	require.NoError(t, err)
	snap := b.Compile()
	require.ElementsMatch(t, []Label{"foo", "bar"}, snap.Labels)
	require.Equal(t, op, snap.Operations[1])
	// require.Equal(t, op.Id(), snap.Timeline[1].Id())
	require.Equal(t, entity.CombineIds(b.Id(), op.Id()), snap.Timeline[1].CombinedId())
	require.ElementsMatch(t, []Label{"foo", "bar"}, snap.Timeline[1].(*LabelChangeTimelineItem).Added)
	require.ElementsMatch(t, []Label{}, snap.Timeline[1].(*LabelChangeTimelineItem).Removed)

	// LabelChange add/remove labels and create a new timeline item
	_, op, err = ChangeLabels(b, rene, time.Now().Unix(), []string{"baz"}, []string{"foo"})
	require.NoError(t, err)
	snap = b.Compile()
	require.ElementsMatch(t, []Label{"bar", "baz"}, snap.Labels)
	require.Equal(t, op, snap.Operations[2])
	// require.Equal(t, op.Id(), snap.Timeline[2].Id())
	require.Equal(t, entity.CombineIds(b.Id(), op.Id()), snap.Timeline[2].CombinedId())
	require.ElementsMatch(t, []Label{"baz"}, snap.Timeline[2].(*LabelChangeTimelineItem).Added)
	require.ElementsMatch(t, []Label{"foo"}, snap.Timeline[2].(*LabelChangeTimelineItem).Removed)
}
