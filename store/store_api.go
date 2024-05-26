package store

import (
	"context"
	"rg/UnitTracker/pkg/proto"
)

type ProjectStore interface {
	//
	CreateProject(context.Context, *proto.Project) (*proto.Project, error)
	GetProject(context.Context, []int32) ([]*proto.Project, error)
	ListProjects(context.Context) ([]*proto.Project, error)
}

type UnitStore interface {
	GetUnits(context.Context, []int32) ([]*proto.Unit, error)
	AddUnit(context.Context, *proto.Unit) (*proto.Unit, error)
}

// Store is the main storage api exposed
type Store interface {
	ProjectStore() ProjectStore
	UnitStore() UnitStore
}

type Connecter interface {
	Connect(context.Context) (Store, error)
}
