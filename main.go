// @title Sentinel API
// @version 0.1.0

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @description ### A domain agnostic permission management and querying API.
// @description
// @description At the most basic level the Sentinel API offers a way for consumers to manage access of application resources to other application resources through the use of policies.
// @description
// @description A simple example demonstrating this flexibility of this relationship could be:
// @description
// @description User (U1) has a policy (P1) granting read/write access to an account Resource (R)
// @description User (U2) has a policy (P2) granting read access to the same account Resource (R)
// @description
// @description The application can then query the api asking if (U1) can preform write access on (R) which will return True
// @description The application can also query the api asking if (U2) can perform write access on (R) which will return False
// @description
// @description More advanced capabilities include:
// @description * Inheritance - Resource access granted to a parent permits the same access to all children
// @description * Invitations - This model allows users having `with_grant` permissions to invite other users to join in a self-service way
// @description * Context - Users with multiple policies can have differing permission. The api permits the caller to specify context when performing permission checks.
// @description
// @description ## ERD
// @description
// @description ![Image](https://www.lucidchart.com/publicSegments/view/e66a76a3-8114-4b1c-9104-765f307c7869/image.png)
// @description
// @description ## Resources
//@description ---
// @description Simply stated a resource represents the nouns in your application. Resources are polymorphic and might represent users,
// @description accounts, or any application specific entity requiring gated access.
// @description
// @description Resources are granted access to other resources through Policies. It is possible for a resource to be the
// @description `principal` in once policy while simultaneously a `target_resource` in another policy.
// @description
// @description
// @description _See definitions below for further explanation of principals and target_resources._
// @description
// @description
// @description Resource may be owned by other resources. Resource inheritance will help reduce redundant permissions to similar resources
// @description if the principal has the same level of access to a class of resource. E.g. System Admins can read/write all accounts.
// @description
// @description **Precedence is determined by path length.** - In the event a resource has two policies with conflicting permission
// @description to a resource, the path with the shortest distance to the resource is preferred.
// @description
// @description #### Example:
// @description ---
// @description **Given:**
// @description
// @description The following resources exist - Parent Resource (PR), Child Resource (CR), Principal (P)
// @description
// @description **AND**
// @description
// @description The following Policies exist - Policy (P1), Policy (P2)
// @description
// @description **AND**
// @description
// @description Parent Resource (PR) owns Child Resource (CR)
// @description
// @description **AND**
// @description
// @description Principal (P) has a policy (P1) explicitly denying a permission, say "read", to (PR)
// @description
// @description **AND**
// @description
// @description Principal (P) has a policy (P2) allowing the same permission, "read", to (CR)
// @description
// @description **Conclusion:**
// @description (P) can not read (PR)
// @description (P) can read (CR)
// @description
// @description In the event there are two policies from a principal directly to a resource with conflicting permission, the permission
// @description check will fail close and access will be denied.
// @description
// @description #### Example:
// @description ---
// @description **Given:**
// @description
// @description The following resources exist - Child Resource (R), Principal (P)
// @description
// @description **AND**
// @description
// @description The following Policies exist - Policy (P1), Policy (P2)
// @description
// @description **AND**
// @description
// @description Principal (P) has a policy (P1) explicitly denying a permission, say "read", to (R)
// @description
// @description **AND**
// @description
// @description Principal (P) has a policy (P2) allowing the same permission, "read", to (R)
// @description
// @description **Conclusion:**
// @description (P) can not read (R)
// @description
// @description ## Policies
// @description
// @description ---
// @description Policies `allow` or `deny` permissions to principals on target resources. A single policy can be granted to zero or
// @description more principals and affect zero or more target resources.
// @description
// @description A policy contains zero or more permissions to target resources. A permission is always scoped to a single target resource within the policy.
// @description However, it is possible to grant the same permission to multiple target resources, or grant multiple permissions to
// @description one target resource within the context of a single policy.
// @description
// @description Each grant of a policy to a principal contains a `with_grant` boolean allowing the aforementioned principal to
// @description propagate their permissions to additional principals if the `with_grant` option has been set to `true`.
// @description By default `with_grant=false`.
// @description
// @description ## Permissions
// @description ---
// @description A permission defines access to perform an action between a Policy and a Resource.
// @description Permissions take simply a `name` and whether the permission is `permitted`.
// @description The `name` is used to determine a principal\'s authorization to a target resource.
// @description
// @description Permission Rules:
// @description * A permission can have `permitted: "allow"` which grants all `principals` of the policy access to the `target_resource` for that permission.
// @description * A permission can have `permitted: "deny"` which prevents all `principals` of the policy from performing the permission action the `target_resource` for that permission.
// @description * A permission can be `revoked` which will remove the `allow` or `deny` constraint. Access may then be based on the inheritance rules described above.
// @description * There is a unique constraint on the `(Policy.id, Permission.name, Resource.id)` tuple.
// @description
// @description ### Definitions
// @description |Entity|Definition|
// @description |------|----------|
// @description |Resource|Any noun in your application. Resources (principals) access other Resources (target resource) through policies.|
// @description |Policy|A named collection of permissions. Policies are granted to resources (principals) and grant or deny access to affect resources (target resources) through permissions.|
// @description |Permission|An explicit allowance or refusal of a resource (principal) to perform an action on a resource (target resource).|
// @description |Principal|In the context of a Policy the principal refers to the resource which has been granted access to make use of the policy.|
// @description |Target Resource|In the context of a Policy a target resource is any and all resources which the policy allows or denies access through Permissions.|
// @description
// @description ### Getting Started
// @description Review the [Getting Started Guide](http://localhost:3000/documentation/GettingStarted.md) for minimal setup instructions
// @host localhost:8080
// @BasePath /
package main

import (
	"github.com/thatbeardo/go-sentinel/models/resource/repository"
	"github.com/thatbeardo/go-sentinel/models/resource/service"
	"github.com/thatbeardo/go-sentinel/models/resource/session"
	"github.com/thatbeardo/go-sentinel/server"
)

func main() {
	shutdown, neo4jsession := server.Initialize()

	rr := repository.New(session.NewNeo4jSession(neo4jsession))
	// resourceRepository := resource.NewNeo4jRepository(session)
	resourceService := service.NewService(rr)

	engine := server.SetupRouter(resourceService)
	server.Orchestrate(engine, shutdown)
}
