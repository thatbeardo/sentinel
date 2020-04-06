package views

type relationships struct {
	Parent   parent   `json:"parent"`
	Policies policies `json:"policies"`
}

type identifier struct {
	Type string `json:"type" enums:"policy, resource, grant, permission"`
	ID   string `json:"id"`
}

type parent struct {
	Data identifier `json:"data"`
}

type policies struct {
	Data []identifier `json:"data"`
}
