package resource

// Service recieves commands from handlers and forwards them to the repository
type Service interface {
	Create(resource *Resource) (string, error)
}
