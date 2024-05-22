package store

import (
	"context"
	"rg/UnitTracker/pkg/proto"
)

type ProjectStore interface {
	//
	CreateProject(context.Context, *proto.Project) (*proto.Project, error)
	GetProject(context.Context, int) (*proto.Project, error)
	// ListProjects(P, error) ([]P, error)
}

type UnitStore interface {
	//
	// GetUnitsForProject(string) ([]U, error) // preview : get this and return first 5 elements
	// AddUnitToProject(U) (error) // maybe return something as ACK
}

// Store is the main storage api exposed
type Store interface {
	ProjectStore() ProjectStore
	UnitStore() UnitStore
}

type Connecter interface {
	Connect(context.Context) (Store, error)
}
