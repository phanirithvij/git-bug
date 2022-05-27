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

func TestAddCommentSerialize(t *testing.T) {
	repo := repository.NewMockRepo()

	rene, err := identity.NewIdentity(repo, "René Descartes", "rene@descartes.fr")
	require.NoError(t, err)

	unix := time.Now().Unix()
	before := NewAddCommentOp(rene, unix, "message", nil)

	data, err := json.Marshal(before)
	require.NoError(t, err)

	var after AddCommentOperation
	err = json.Unmarshal(data, &after)
	require.NoError(t, err)

	// enforce creating the ID
	before.Id()

	// Replace the identity as it's not serialized
	after.Author_ = rene

	require.Equal(t, before, &after)
}

func TestAddComment(t *testing.T) {
	repo := repository.NewMockRepo()

	rene, err := identity.NewIdentity(repo, "René Descartes", "rene@descartes.fr")
	require.NoError(t, err)

	b, _, err := Create(rene, time.Now().Unix(), "title", "")
	require.NoError(t, err)

	op, err := AddComment(b, rene, time.Now().Unix(), "comment1")
	require.NoError(t, err)

	// AddComment add a new comment and create a new timeline item
	snap := b.Compile()
	require.Len(t, snap.Comments, 2)
	require.Equal(t, op, snap.Operations[1])
	require.Equal(t, "comment1", snap.Comments[1].Message)
	require.Equal(t, rene, snap.Comments[1].Author)
	// require.Equal(t, op.Id(), snap.Timeline[1].Id())
	require.Equal(t, entity.CombineIds(b.Id(), op.Id()), snap.Timeline[1].CombinedId())
	require.Equal(t, "comment1", snap.Timeline[1].(*AddCommentTimelineItem).Message)
	require.Equal(t, rene, snap.Timeline[1].(*AddCommentTimelineItem).Author)

}
