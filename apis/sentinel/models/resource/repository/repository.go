package repository

import (
	models "github.com/bithippie/guard-my-app/apis/sentinel/models"
	policy "github.com/bithippie/guard-my-app/apis/sentinel/models/policy/dto"
	policies "github.com/bithippie/guard-my-app/apis/sentinel/models/policy/session"
	resource "github.com/bithippie/guard-my-app/apis/sentinel/models/resource/dto"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/resource/session"
)

// Repository is used by the service to communicate with the underlying database
type Repository interface {
	Get(string) (resource.Output, error)
	GetByID(string) (resource.OutputDetails, error)

	AttachResourceToExistingParent(*resource.Input) (resource.OutputDetails, error)
	AttachResourceToTenantPolicy(string, *resource.Input) (resource.OutputDetails, error)
	CreateTenantResource(*resource.Input) (resource.OutputDetails, error)

	AssociatePolicy(string, *policy.Input) (policy.OutputDetails, error)
	GetAllAssociatedPolicies(string) (policy.Output, error)
	Update(resource.Details, *resource.Input) (resource.OutputDetails, error)
	Delete(string) error
}

type repository struct {
	session       session.Session
	policySession policies.Session
}

// Get retrieves all the resources present in the graph
func (repo *repository) Get(tenantID string) (resource.Output, error) {
	return repo.session.Execute(`
		MATCH(n:Resource{source_id:$tenant_id})<-[:GRANTED_TO]-(:Policy)-[:PERMISSION]->(root:Resource)<-[:OWNED_BY*0..]-(child:Resource)
		OPTIONAL MATCH (child)-[:OWNED_BY]->(parent: Resource)
		OPTIONAL MATCH (policy: Policy)-[:GRANTED_TO]->(child)
		RETURN {child: child, parent: parent, policy: COLLECT(policy)}`,
		map[string]interface{}{
			"tenant_id": tenantID,
		})
}

// GetByID function adds a resource node
func (repo *repository) GetByID(id string) (resource.OutputDetails, error) {
	results, err := repo.session.Execute(`
		MATCH(child:Resource{id: $id})
		OPTIONAL MATCH (child: Resource)-[:OWNED_BY]->(parent: Resource)
		OPTIONAL MATCH (policy: Policy)-[:GRANTED_TO]->(child)
		RETURN {child: child, parent: parent, policy: COLLECT(policy)}`,
		map[string]interface{}{
			"id": id,
		})
	if len(results.Data) == 0 {
		return resource.OutputDetails{}, models.ErrNotFound
	}
	return resource.OutputDetails{Data: results.Data[0]}, err
}

// Create function adds a node to the graph

func (repo *repository) AttachResourceToExistingParent(input *resource.Input) (resource.OutputDetails, error) {
	results, err := repo.session.Execute(`
		MATCH (parent:Resource{id:$parent_id})
		WITH parent
		CREATE (child:Resource{name:$name, source_id:$source_id, id:randomUUID()})-[:OWNED_BY]->(parent)
		RETURN {child:child, parent:parent}`,
		map[string]interface{}{
			"name":      input.Data.Attributes.Name,
			"source_id": input.Data.Attributes.SourceID,
			"parent_id": input.Data.Relationships.Parent.Data.ID,
		})
	if len(results.Data) == 0 {
		return resource.OutputDetails{}, models.ErrNotFound
	}
	return resource.OutputDetails{Data: results.Data[0]}, err
}

// Create function adds a node to the graph
func (repo *repository) AttachResourceToTenantPolicy(tenantPolicy string, input *resource.Input) (resource.OutputDetails, error) {
	results, err := repo.session.Execute(`
		MATCH(tenant:Resource{source_id: $tenant_id})<-[:GRANTED_TO]-(policy:Policy)
		WITH policy
		CREATE (policy)-[:PERMISSION]->(child:Resource{name:$name, source_id:$source_id, id:randomUUID()})
		RETURN {child:child}`,
		map[string]interface{}{
			"name":      input.Data.Attributes.Name,
			"source_id": input.Data.Attributes.SourceID,
			"tenant_id": tenantPolicy,
		})
	if len(results.Data) == 0 {
		return resource.OutputDetails{}, models.ErrNotFound
	}
	return resource.OutputDetails{Data: results.Data[0]}, err
}

// CreateTenantResource is called only when guard my app is requesting to create a resource.
func (repo *repository) CreateTenantResource(input *resource.Input) (resource.OutputDetails, error) {
	results, err := repo.session.Execute(`
		CREATE (child:Resource{name:$name, source_id:$source_id, id: randomUUID()})
		RETURN {child:child}`,
		map[string]interface{}{
			"name":      input.Data.Attributes.Name,
			"source_id": input.Data.Attributes.SourceID,
		})
	if len(results.Data) == 0 {
		return resource.OutputDetails{}, models.ErrNotFound
	}
	return resource.OutputDetails{Data: results.Data[0]}, err
}

func (repo *repository) AssociatePolicy(id string, input *policy.Input) (policy.OutputDetails, error) {
	result, err := repo.policySession.Execute(`
		MATCH(principal:Resource{id: $principalID})
		CREATE (principal)<-[:GRANTED_TO]-(policy:Policy{ name:$name, id:randomUUID() })
		RETURN {policy: policy, principals: COLLECT(principal)}`,
		map[string]interface{}{
			"principalID": id,
			"name":        input.Data.Attributes.Name,
		},
	)
	if len(result.Data) == 0 {
		return policy.OutputDetails{}, models.ErrDatabase
	}
	return policy.OutputDetails{Data: result.Data[0]}, err
}

// Update function Edits the contents of a node
func (repo *repository) Update(oldResource resource.Details, newResource *resource.Input) (resource.OutputDetails, error) {
	newParentID := extractParentID(newResource)
	statement := generateUpdateStatement(newParentID)
	results, err := repo.session.Execute(statement,
		map[string]interface{}{
			"child_id":      oldResource.ID,
			"name":          newResource.Data.Attributes.Name,
			"source_id":     newResource.Data.Attributes.SourceID,
			"new_parent_id": newParentID,
		})
	if len(results.Data) == 0 {
		return resource.OutputDetails{}, models.ErrNotFound
	}
	return resource.OutputDetails{Data: results.Data[0]}, err
}

// Delete function deletes a node from the graph
func (repo *repository) Delete(id string) error {
	_, err := repo.session.Execute(`
		MATCH (n:Resource { id: $id }) DETACH DELETE n`,
		map[string]interface{}{
			"id": id,
		})
	return err
}

func (repo *repository) GetAllAssociatedPolicies(id string) (policy.Output, error) {
	return repo.policySession.Execute(`
		MATCH (resource:Resource{id: $id}) 
		WITH resource
		MATCH(policy)-[:GRANTED_TO]->(resource)
		WITH policy
		OPTIONAL MATCH(policy)-[:PERMISSION]->(target:Resource)
		OPTIONAL MATCH(policy)-[:GRANTED_TO]->(principal:Resource)
		RETURN {policy:policy, principals:COLLECT(principal), targets:COLLECT(target)}`,
		map[string]interface{}{
			"id": id,
		})
}

// New is a factory method to generate repository instances
func New(session session.Session, policySession policies.Session) Repository {
	return &repository{
		session:       session,
		policySession: policySession,
	}
}

func extractParentID(newResource *resource.Input) string {
	var parentID string
	if newResource.Data.Relationships != nil {
		parentID = newResource.Data.Relationships.Parent.Data.ID
	}
	return parentID
}

func generateUpdateStatement(newParentID string) (statement string) {
	statement = `
		MATCH(child:Resource{id:$child_id})
		SET child.name=$name
		SET child.source_id=$source_id
		WITH child
		OPTIONAL MATCH (child)-[old_relationship:OWNED_BY]->(old_parent:Resource)`

	if newParentID == "" {
		statement += `RETURN {child: child, parent: old_parent}`
	} else {
		statement += `
		DETACH DELETE old_relationship
		WITH child

		OPTIONAL MATCH (new_parent:Resource{id:$new_parent_id})
		CREATE(child)-[:OWNED_BY]->(new_parent)
		RETURN {child: child, parent: new_parent}`
	}
	return
}
