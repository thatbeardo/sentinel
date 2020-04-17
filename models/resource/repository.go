package resource

// Repository is used by the service to communicate with the underlying database
type Repository interface {
	Get() (Response, error)
	GetByID(string) (Element, error)
	Create(*Input) (Element, error)
	Update(string, *Input) (Element, error)
	Delete(string) error

	UpdateOwnership(string, *Input) error
	CreateEdge(string, string) error
	DeleteEdge(string, string) error
}
