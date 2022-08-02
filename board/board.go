package board

import (
	"github.com/MichaelMure/git-bug/entity"
	"github.com/MichaelMure/git-bug/entity/dag"
	"github.com/MichaelMure/git-bug/identity"
	"github.com/MichaelMure/git-bug/repository"
)

var _ entity.Interface = &Board{}

// 1: original format
const formatVersion = 1

var def = dag.Definition{
	Typename:             "board",
	Namespace:            "boards",
	OperationUnmarshaler: operationUnmarshaller,
	FormatVersion:        formatVersion,
}

var ClockLoader = dag.ClockLoader(def)

// Board holds the data of a project board.
type Board struct {
	*dag.Entity
}

// NewBoard create a new Board
func NewBoard() *Board {
	return &Board{
		Entity: dag.New(def),
	}
}

// Read will read a board from a repository
func Read(repo repository.ClockedRepo, id entity.Id) (*Board, error) {
	return ReadWithResolver(repo, identity.NewSimpleResolver(repo), id)
}

// ReadWithResolver will read a board from its Id, with a custom identity.Resolver
func ReadWithResolver(repo repository.ClockedRepo, identityResolver identity.Resolver, id entity.Id) (*Board, error) {
	e, err := dag.Read(def, repo, identityResolver, id)
	if err != nil {
		return nil, err
	}
	return &Board{Entity: e}, nil
}
