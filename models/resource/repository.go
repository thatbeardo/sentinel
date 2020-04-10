package resource

// Repository is used by the service to communicate with the underlying database
type Repository interface {
	Get() (Response, error)
	Create(*Input) (Element, error)
	GetByID(string) (Element, error)
	Delete(string) error
}
