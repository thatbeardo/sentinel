package repository

import (
	models "github.com/bithippie/guard-my-app/sentinel/models"
	permission "github.com/bithippie/guard-my-app/sentinel/models/permission/dto"
	"github.com/bithippie/guard-my-app/sentinel/models/permission/session"
)

// Repository exposes wrapper methods around the underlying session
type Repository interface {
	GetAllPermissionsForPolicy(string) (permission.Output, error)
	GetAllPermissionsForPolicyWithResource(string, string) (permission.Output, error)
	Create(*permission.Input, string, string) (permission.OutputDetails, error)
	Update(string, *permission.Input) (permission.OutputDetails, error)
	Delete(string) error
}

type repository struct {
	session session.Session
}

// New is a factory method to generate repository instances
func New(session session.Session) Repository {
	return &repository{
		session: session,
	}
}

func (repo *repository) GetAllPermissionsForPolicy(policyID string) (permission.Output, error) {
	return repo.session.Execute(`
		MATCH (policy:Policy { id: $policyID } )-[permission:PERMISSION]->(resource:Resource)
		RETURN {permission: permission}`,
		map[string]interface{}{
			"policyID": policyID,
		})
}

func (repo *repository) GetAllPermissionsForPolicyWithResource(policyID string, resourceID string) (permission.Output, error) {
	return repo.session.Execute(`
		MATCH (policy:Policy { id: $policyID } )-[permission:PERMISSION]->(resource:Resource { id: $resourceID })
		RETURN {permission: permission}`,
		map[string]interface{}{
			"policyID":   policyID,
			"resourceID": resourceID,
		})
}

func (repo *repository) Create(input *permission.Input, policyID string, targetID string) (response permission.OutputDetails, err error) {
	result, err := repo.session.Execute(`
		MATCH (policy: Policy), (target: Resource)
		WHERE policy.id = $policyID AND target.id = $targetID
		CREATE (policy)-[r:PERMISSION {name: $name, permitted: $permitted, id: randomUUID()}]->(target)
		RETURN {permission: r}`,
		map[string]interface{}{
			"name":      input.Data.Attributes.Name,
			"permitted": input.Data.Attributes.Permitted,
			"policyID":  policyID,
			"targetID":  targetID,
		})

	if len(result.Data) == 0 {
		return permission.OutputDetails{}, models.ErrNotFound
	}
	response = permission.OutputDetails{Data: result.Data[0]}
	return
}

func (repo *repository) Update(id string, input *permission.Input) (response permission.OutputDetails, err error) {
	result, err := repo.session.Execute(`
		MATCH(policy:Policy)-[permission:PERMISSION]->(resource:Resource)
		WHERE permission.id = $id
		SET permission.name = $name
		SET permission.permitted = $permitted
		RETURN {permission: permission}`,
		map[string]interface{}{
			"name":      input.Data.Attributes.Name,
			"permitted": input.Data.Attributes.Permitted,
			"id":        id,
		})

	if len(result.Data) == 0 {
		return permission.OutputDetails{}, models.ErrNotFound
	}
	response = permission.OutputDetails{Data: result.Data[0]}
	return
}

func (repo *repository) Delete(id string) (err error) {
	_, err = repo.session.Execute(`
		MATCH(policy:Policy)-[permission:PERMISSION]->(resource:Resource)
		WHERE permission.id = $id
		DELETE permission`,
		map[string]interface{}{
			"id": id,
		})
	return
}
