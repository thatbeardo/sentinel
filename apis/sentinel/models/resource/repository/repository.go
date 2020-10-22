package repository

import (
	"fmt"

	models "github.com/bithippie/guard-my-app/apis/sentinel/models"
	context "github.com/bithippie/guard-my-app/apis/sentinel/models/context/dto"
	contexts "github.com/bithippie/guard-my-app/apis/sentinel/models/context/session"
	resource "github.com/bithippie/guard-my-app/apis/sentinel/models/resource/dto"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/resource/session"
)

// Repository is used by the service to communicate with the underlying database
type Repository interface {
	Get(clientID, tenant string) (resource.Output, error)
	GetResourcesHub(clientID, tenant string) (resource.OutputDetails, error)
	GetByID(string) (resource.OutputDetails, error)

	Create(*resource.Input) (resource.OutputDetails, error)
	CreateMetaResource(*resource.Input) (resource.OutputDetails, error)

	AddContext(string, *context.Input) (context.OutputDetails, error)
	GetAllContexts(string) (context.Output, error)

	Update(resource.Details, *resource.Input) (resource.OutputDetails, error)
	Delete(string) error
}

type repository struct {
	session        session.Session
	contextSession contexts.Session
}

// Get retrieves all the resources present in the graph
func (repo *repository) Get(clientID, tenant string) (resource.Output, error) {
	return repo.session.Execute(`
		MATCH(client:Resource{source_id:$client_id})
		<-[:GRANTED_TO]-(:Context)-[:PERMISSION{name:"sentinel:read"}]->
		(tenant:Resource{source_id: $tenant})<-[:GRANTED_TO]-(:Context)-[:PERMISSION{name:"sentinel:read"}]->(hub:Resource)
		<-[:OWNED_BY*0..]-(child:Resource)
		OPTIONAL MATCH (child)-[:OWNED_BY]->(parent: Resource)
		OPTIONAL MATCH (context:Context)-[:GRANTED_TO]->(child)
		RETURN {child: child, parent: parent, context: COLLECT(context)}`,
		map[string]interface{}{
			"client_id": clientID,
			"tenant":    tenant,
		})
}

// GetByID function adds a resource node
func (repo *repository) GetByID(id string) (resource.OutputDetails, error) {
	results, err := repo.session.Execute(`
		MATCH(child:Resource{id: $id})
		OPTIONAL MATCH (child)-[:OWNED_BY]->(parent: Resource)
		OPTIONAL MATCH (context:Context)-[:GRANTED_TO]->(child)
		RETURN {child: child, parent: parent, context: COLLECT(context)}`,
		map[string]interface{}{
			"id": id,
		})
	if len(results.Data) == 0 {
		return resource.OutputDetails{}, models.ErrNotFound
	}
	return resource.OutputDetails{Data: results.Data[0]}, err
}

// GetResourcesHub returns the resource to which all other nodes hang off of
func (repo *repository) GetResourcesHub(clientID, tenant string) (resource.OutputDetails, error) {
	results, err := repo.session.Execute(`
		MATCH(client:Resource{source_id: $client_id})<-[:GRANTED_TO]-(:Context)
		-[:PERMISSION{name:"sentinel:read", permitted: "allow"}]->
		(tenant:Resource{source_id: $tenant})<-[:GRANTED_TO]-(:Context)
		-[:PERMISSION{name:"sentinel:read", permitted: "allow"}]->(hub:Resource)
		RETURN {child: hub}`,
		map[string]interface{}{
			"client_id": clientID,
			"tenant":    tenant,
		})
	if len(results.Data) == 0 {
		return resource.OutputDetails{}, models.ErrNotFound
	}
	return resource.OutputDetails{Data: results.Data[0]}, err
}

// Create function adds a node to the graph - typically invoked by customer API not guard-my-app
func (repo *repository) Create(input *resource.Input) (resource.OutputDetails, error) {
	results, err := repo.session.Execute(`
		MATCH (parent:Resource{id:$parent_id})
		WITH parent
		CREATE (child:Resource{name:$name, source_id:$source_id, context_id:$context_id, id:randomUUID()})-[:OWNED_BY]->(parent)
		RETURN {child:child, parent:parent}`,
		map[string]interface{}{
			"name":       input.Data.Attributes.Name,
			"source_id":  input.Data.Attributes.SourceID,
			"parent_id":  input.Data.Relationships.Parent.Data.ID,
			"context_id": input.Data.Attributes.ContextID,
		})
	if len(results.Data) == 0 {
		return resource.OutputDetails{}, models.ErrNotFound
	}
	return resource.OutputDetails{Data: results.Data[0]}, err
}

// CreateTenantResource is called only when guard my app is requesting to create a resource.
func (repo *repository) CreateMetaResource(input *resource.Input) (resource.OutputDetails, error) {
	results, err := repo.session.Execute(`
		CREATE (child:Resource{name:$name, source_id:$source_id, context_id:$context_id, id: randomUUID()})
		RETURN {child:child}`,
		map[string]interface{}{
			"name":       input.Data.Attributes.Name,
			"source_id":  input.Data.Attributes.SourceID,
			"context_id": input.Data.Attributes.ContextID,
		})
	if len(results.Data) == 0 {
		return resource.OutputDetails{}, models.ErrNotFound
	}
	return resource.OutputDetails{Data: results.Data[0]}, err
}

func (repo *repository) AddContext(id string, input *context.Input) (context.OutputDetails, error) {
	result, err := repo.contextSession.Execute(`
		MATCH(principal:Resource{id: $principalID})
		CREATE (principal)<-[:GRANTED_TO]-(context:Context{ name:$name, id:randomUUID() })
		RETURN {context:context, principals: COLLECT(principal)}`,
		map[string]interface{}{
			"principalID": id,
			"name":        input.Data.Attributes.Name,
		},
	)
	if len(result.Data) == 0 {
		return context.OutputDetails{}, models.ErrDatabase
	}
	return context.OutputDetails{Data: result.Data[0]}, err
}

// Update function Edits the contents of a node
func (repo *repository) Update(oldResource resource.Details, newResource *resource.Input) (resource.OutputDetails, error) {
	newParentID := extractParentID(newResource)
	statement := generateUpdateStatement(newParentID, newResource)
	results, err := repo.session.Execute(statement,
		map[string]interface{}{
			"child_id":      oldResource.ID,
			"name":          newResource.Data.Attributes.Name,
			"source_id":     newResource.Data.Attributes.SourceID,
			"context_id":    newResource.Data.Attributes.ContextID,
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

func (repo *repository) GetAllContexts(id string) (context.Output, error) {
	return repo.contextSession.Execute(`
		MATCH (resource:Resource{id: $id}) 
		WITH resource
		MATCH(context)-[:GRANTED_TO]->(resource)
		WITH context
		OPTIONAL MATCH(context)-[:PERMISSION]->(target:Resource)
		OPTIONAL MATCH(context)-[:GRANTED_TO]->(principal:Resource)
		RETURN {context:context, principals:COLLECT(principal), targets:COLLECT(target)}`,
		map[string]interface{}{
			"id": id,
		})
}

// New is a factory method to generate repository instances
func New(session session.Session, contextSession contexts.Session) Repository {
	return &repository{
		session:        session,
		contextSession: contextSession,
	}
}

func extractParentID(newResource *resource.Input) string {
	var parentID string
	if newResource.Data.Relationships != nil {
		parentID = newResource.Data.Relationships.Parent.Data.ID
	}
	return parentID
}

func generateUpdateStatement(newParentID string, input *resource.Input) (statement string) {
	var updatedFields string

	if input.Data.Attributes.Name != "" {
		updatedFields += `SET child.name=$name `
	}

	if input.Data.Attributes.ContextID != "" {
		updatedFields += `SET child.context_id=$context_id `
	}

	if input.Data.Attributes.SourceID != "" {
		updatedFields += `SET child.source_id=$source_id `
	}

	statement = fmt.Sprintf(`
		MATCH(child:Resource{id:$child_id})
		%s
		WITH child
		OPTIONAL MATCH (child)-[old_relationship:OWNED_BY]->(old_parent:Resource)`, updatedFields)

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
