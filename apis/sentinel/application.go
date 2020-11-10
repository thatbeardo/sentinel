// @title Sentinel API
// @version 0.1.0

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @description ### A domain agnostic permission management and querying API. to manage
// @description
// @description At the most basic level the Sentinel API offers a way for consumers to manage access of application resources to other application resources through the use of contexts.
// @description
// @description A simple example demonstrating this flexibility of this relationship could be:
// @description
// @description User (U1) has a context (P1) granting read/write access to an account Resource (R)
// @description User (U2) has a context (P2) granting read access to the same account Resource (R)
// @description
// @description The application can then query the api asking if (U1) can preform write access on (R) which will return True
// @description The application can also query the api asking if (U2) can perform write access on (R) which will return False
// @description
// @description More advanced capabilities include:
// @description * Inheritance - Resource access granted to a parent permits the same access to all children
// @description * Invitations - This model allows users having `with_grant` permissions to invite other users to join in a self-service way
// @description * Context - Users with multiple contexts can have differing permission. The api permits the caller to specify context when performing permission checks.
// @description
// @description ## Resources
//@description ---
// @description Simply stated a resource represents the nouns in your application. Resources are polymorphic and might represent users,
// @description accounts, or any application specific entity requiring gated access.
// @description
// @description Resources are granted access to other resources through Contexts. It is possible for a resource to be the
// @description `principal` in once context while simultaneously a `target_resource` in another context.
// @description
// @description
// @description _See definitions below for further explanation of principals and target_resources._
// @description
// @description
// @description Resource may be owned by other resources. Resource inheritance will help reduce redundant permissions to similar resources
// @description if the principal has the same level of access to a class of resource. E.g. System Admins can read/write all accounts.
// @description
// @description **Precedence is determined by path length.** - In the event a resource has two contexts with conflicting permission
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
// @description The following Contexts exist - context (P1), context (P2)
// @description
// @description **AND**
// @description
// @description Parent Resource (PR) owns Child Resource (CR)
// @description
// @description **AND**
// @description
// @description Principal (P) has a context (P1) explicitly denying a permission, say "read", to (PR)
// @description
// @description **AND**
// @description
// @description Principal (P) has a context (P2) allowing the same permission, "read", to (CR)
// @description
// @description **Conclusion:**
// @description (P) can not read (PR)
// @description (P) can read (CR)
// @description
// @description In the event there are two contexts from a principal directly to a resource with conflicting permission, the permission
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
// @description The following Contexts exist - context (P1), context (P2)
// @description
// @description **AND**
// @description
// @description Principal (P) has a context (P1) explicitly denying a permission, say "read", to (R)
// @description
// @description **AND**
// @description
// @description Principal (P) has a context (P2) allowing the same permission, "read", to (R)
// @description
// @description **Conclusion:**
// @description (P) can not read (R)
// @description
// @description ## Contexts
// @description
// @description ---
// @description Contexts `allow` or `deny` permissions to principals on target resources. A single context can be granted to zero or
// @description more principals and affect zero or more target resources.
// @description
// @description A context contains zero or more permissions to target resources. A permission is always scoped to a single target resource within the context.
// @description However, it is possible to grant the same permission to multiple target resources, or grant multiple permissions to
// @description one target resource within the context of a single context.
// @description
// @description Each grant of a context to a principal contains a `with_grant` boolean allowing the aforementioned principal to
// @description propagate their permissions to additional principals if the `with_grant` option has been set to `true`.
// @description By default `with_grant=false`.
// @description
// @description ## Permissions
// @description ---
// @description A permission defines access to perform an action between a context and a Resource.
// @description Permissions take simply a `name` and whether the permission is `permitted`.
// @description The `name` is used to determine a principal\'s authorization to a target resource.
// @description
// @description Permission Rules:
// @description * A permission can have `permitted: "allow"` which grants all `principals` of the context access to the `target_resource` for that permission.
// @description * A permission can have `permitted: "deny"` which prevents all `principals` of the context from performing the permission action the `target_resource` for that permission.
// @description * A permission can be `revoked` which will remove the `allow` or `deny` constraint. Access may then be based on the inheritance rules described above.
// @description * There is a unique constraint on the `(context.id, Permission.name, Resource.id)` tuple.
// @description
// @description ### Definitions
// @description |Entity|Definition|
// @description |------|----------|
// @description |Resource|Any noun in your application. Resources (principals) access other Resources (target resource) through contexts.|
// @description |context|A named collection of permissions. Contexts are granted to resources (principals) and grant or deny access to affect resources (target resources) through permissions.|
// @description |Permission|An explicit allowance or refusal of a resource (principal) to perform an action on a resource (target resource).|
// @description |Principal|In the context of a context the principal refers to the resource which has been granted access to make use of the context.|
// @description |Target Resource|In the context of a context a target resource is any and all resources which the context allows or denies access through Permissions.|
// @description
// @description ### Getting Started
// @description Review the [Getting Started Guide](http://localhost:3000/documentation/GettingStarted.md) for minimal setup instructions
//
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
//
// @BasePath /
package main

import (
	"os"

	authorizations "github.com/bithippie/guard-my-app/apis/sentinel/api/handlers/authorization"
	contexts "github.com/bithippie/guard-my-app/apis/sentinel/api/handlers/contexts"
	"github.com/bithippie/guard-my-app/apis/sentinel/api/handlers/grants"
	"github.com/bithippie/guard-my-app/apis/sentinel/api/handlers/healthcheck"
	"github.com/bithippie/guard-my-app/apis/sentinel/api/handlers/login"
	"github.com/bithippie/guard-my-app/apis/sentinel/api/handlers/permissions"
	"github.com/bithippie/guard-my-app/apis/sentinel/api/handlers/resources"
	"github.com/bithippie/guard-my-app/apis/sentinel/api/middleware"
	contextRepository "github.com/bithippie/guard-my-app/apis/sentinel/models/context/repository"
	contextService "github.com/bithippie/guard-my-app/apis/sentinel/models/context/service"
	contextSession "github.com/bithippie/guard-my-app/apis/sentinel/models/context/session"
	grantRepository "github.com/bithippie/guard-my-app/apis/sentinel/models/grant/repository"
	grantService "github.com/bithippie/guard-my-app/apis/sentinel/models/grant/service"
	grantSession "github.com/bithippie/guard-my-app/apis/sentinel/models/grant/session"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/neo4j"
	permissionRepository "github.com/bithippie/guard-my-app/apis/sentinel/models/permission/repository"
	permissionService "github.com/bithippie/guard-my-app/apis/sentinel/models/permission/service"
	permissionSession "github.com/bithippie/guard-my-app/apis/sentinel/models/permission/session"
	resourceRepository "github.com/bithippie/guard-my-app/apis/sentinel/models/resource/repository"
	resourceService "github.com/bithippie/guard-my-app/apis/sentinel/models/resource/service"
	resourceSession "github.com/bithippie/guard-my-app/apis/sentinel/models/resource/session"
	"github.com/bithippie/guard-my-app/apis/sentinel/statsd"

	authorizationRepository "github.com/bithippie/guard-my-app/apis/sentinel/models/authorization/repository"
	authorizationService "github.com/bithippie/guard-my-app/apis/sentinel/models/authorization/service"
	authorizationSession "github.com/bithippie/guard-my-app/apis/sentinel/models/authorization/session"
	"github.com/bithippie/guard-my-app/apis/sentinel/server"
	"github.com/gin-gonic/gin"
	// "github.com/newrelic/go-agent/_integrations/nrgin/v1"
	log "github.com/sirupsen/logrus"
)

func main() {
	statsdClient, err := statsd.New(os.Getenv("STATSD_HOST"), os.Getenv("STATSD_PORT"))
	driver := server.CreateDatabaseDriver()
	runner := neo4j.NewRunner(driver)

	contextsSession := contextSession.NewNeo4jSession(runner)
	contextRepository := contextRepository.New(contextsSession)
	contextService := contextService.New(contextRepository)

	resourceSession := resourceSession.NewNeo4jSession(runner)
	resourceRepository := resourceRepository.New(resourceSession, contextsSession)
	resourceService := resourceService.New(resourceRepository)

	permissionSession := permissionSession.NewNeo4jSession(runner)
	permissionRepository := permissionRepository.New(permissionSession)
	permissionService := permissionService.New(permissionRepository)

	grantSession := grantSession.NewNeo4jSession(runner)
	grantRepository := grantRepository.New(grantSession)
	grantService := grantService.New(grantRepository)

	authorizationSession := authorizationSession.NewNeo4jSession(runner)
	authorizationRepository := authorizationRepository.New(authorizationSession)
	authorizationService := authorizationService.New(authorizationRepository)

	engine := gin.Default()
	router := server.GenerateRouter(engine)

	if err == nil {
		log.Info("StatsD connection established. Tracking data")
		router.Use(middleware.Metrics(statsdClient))
	}

	healthcheck.Routes(router)
	login.Routes(router)
	// router.Use(nrgin.Middleware(server.InitNewRelicApp()))
	router.Use(middleware.VerifyClaimant)
	router.Use(middleware.VerifyToken)
	router.Use(middleware.VerifyTenant)

	resources.Routes(router, resourceService, authorizationService)
	permissions.Routes(router, permissionService, authorizationService)
	contexts.Routes(router, contextService, authorizationService)
	grants.Routes(router, grantService, authorizationService)
	authorizations.Routes(router, authorizationService)
  
	server.Orchestrate(engine, driver)
}
