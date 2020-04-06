package resource

// Repository is used by the service to communicate with the underlying database
type Repository interface {
	Get() ([]*Resource, error)
	Create(*Resource) (*Resource, error)
}
