package store

type ProjectStore interface{
}
type UnitStore interface{
}

// Store is the main 
type Store interface {
   ProjectStore()
   UnitStore()
}
