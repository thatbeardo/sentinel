package resource

// Repository is used by the service to communicate with the underlying database
type Repository interface {
	Get() (Response, error)
	Create(*Input) (Response, error)
}
