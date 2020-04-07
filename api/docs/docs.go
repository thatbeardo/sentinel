// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag

package docs

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "API Support",
            "email": "hmavani7@gmail.com"
        },
        "license": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/resources": {
            "get": {
                "description": "Get all the nodes present in the graph",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Get all the resources",
                "responses": {
                    "200": {
                        "description": "ok",
                        "schema": {
                            "$ref": "#/definitions/resource.Response"
                        }
                    }
                }
            },
            "post": {
                "description": "Add a new resource to existing resources",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Create a new Resource",
                "responses": {
                    "202": {
                        "description": "ok",
                        "schema": {
                            "$ref": "#/definitions/resource.Response"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "resource.Response": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "array",
                    "items": {
                        "type": "Dto"
                    }
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "0.1.0",
	Host:        "localhost:8080",
	BasePath:    "/",
	Schemes:     []string{},
	Title:       "Sentinel API",
	Description: "### A domain agnostic permission management and querying API.\n\nAt the most basic level the Sentinel API offers a way for consumers to manage access of application resources to other application resources through the use of policies.\n\nA simple example demonstrating this flexibility of this relationship could be:\n\nUser (U1) has a policy (P1) granting read/write access to an account Resource (R)\nUser (U2) has a policy (P2) granting read access to the same account Resource (R)\n\nThe application can then query the api asking if (U1) can preform write access on (R) which will return True\nThe application can also query the api asking if (U2) can perform write access on (R) which will return False\n\nMore advanced capabilities include:\n* Inheritance - Resource access granted to a parent permits the same access to all children\n* Invitations - This model allows users having `with_grant` permissions to invite other users to join in a self-service way\n* Context - Users with multiple policies can have differing permission. The api permits the caller to specify context when performing permission checks.\n\n## ERD\n\n![Image](https://www.lucidchart.com/publicSegments/view/e66a76a3-8114-4b1c-9104-765f307c7869/image.png)\n\n## Resources\n---\nSimply stated a resource represents the nouns in your application. Resources are polymorphic and might represent users,\naccounts, or any application specific entity requiring gated access.\n\nResources are granted access to other resources through Policies. It is possible for a resource to be the\n`principal` in once policy while simultaneously a `target_resource` in another policy.\n\n\n_See definitions below for further explanation of principals and target_resources._\n\n\nResource may be owned by other resources. Resource inheritance will help reduce redundant permissions to similar resources\nif the principal has the same level of access to a class of resource. E.g. System Admins can read/write all accounts.\n\n**Precedence is determined by path length.** - In the event a resource has two policies with conflicting permission\nto a resource, the path with the shortest distance to the resource is preferred.\n\n#### Example:\n---\n**Given:**\n\nThe following resources exist - Parent Resource (PR), Child Resource (CR), Principal (P)\n\n**AND**\n\nThe following Policies exist - Policy (P1), Policy (P2)\n\n**AND**\n\nParent Resource (PR) owns Child Resource (CR)\n\n**AND**\n\nPrincipal (P) has a policy (P1) explicitly denying a permission, say \"read\", to (PR)\n\n**AND**\n\nPrincipal (P) has a policy (P2) allowing the same permission, \"read\", to (CR)\n\n**Conclusion:**\n(P) can not read (PR)\n(P) can read (CR)\n\nIn the event there are two policies from a principal directly to a resource with conflicting permission, the permission\ncheck will fail close and access will be denied.\n\n#### Example:\n---\n**Given:**\n\nThe following resources exist - Child Resource (R), Principal (P)\n\n**AND**\n\nThe following Policies exist - Policy (P1), Policy (P2)\n\n**AND**\n\nPrincipal (P) has a policy (P1) explicitly denying a permission, say \"read\", to (R)\n\n**AND**\n\nPrincipal (P) has a policy (P2) allowing the same permission, \"read\", to (R)\n\n**Conclusion:**\n(P) can not read (R)\n\n## Policies\n\n---\nPolicies `allow` or `deny` permissions to principals on target resources. A single policy can be granted to zero or\nmore principals and affect zero or more target resources.\n\nA policy contains zero or more permissions to target resources. A permission is always scoped to a single target resource within the policy.\nHowever, it is possible to grant the same permission to multiple target resources, or grant multiple permissions to\none target resource within the context of a single policy.\n\nEach grant of a policy to a principal contains a `with_grant` boolean allowing the aforementioned principal to\npropagate their permissions to additional principals if the `with_grant` option has been set to `true`.\nBy default `with_grant=false`.\n\n## Permissions\n---\nA permission defines access to perform an action between a Policy and a Resource.\nPermissions take simply a `name` and whether the permission is `permitted`.\nThe `name` is used to determine a principal\\'s authorization to a target resource.\n\nPermission Rules:\n* A permission can have `permitted: \"allow\"` which grants all `principals` of the policy access to the `target_resource` for that permission.\n* A permission can have `permitted: \"deny\"` which prevents all `principals` of the policy from performing the permission action the `target_resource` for that permission.\n* A permission can be `revoked` which will remove the `allow` or `deny` constraint. Access may then be based on the inheritance rules described above.\n* There is a unique constraint on the `(Policy.id, Permission.name, Resource.id)` tuple.\n\n### Definitions\n|Entity|Definition|\n|------|----------|\n|Resource|Any noun in your application. Resources (principals) access other Resources (target resource) through policies.|\n|Policy|A named collection of permissions. Policies are granted to resources (principals) and grant or deny access to affect resources (target resources) through permissions.|\n|Permission|An explicit allowance or refusal of a resource (principal) to perform an action on a resource (target resource).|\n|Principal|In the context of a Policy the principal refers to the resource which has been granted access to make use of the policy.|\n|Target Resource|In the context of a Policy a target resource is any and all resources which the policy allows or denies access through Permissions.|\n\n### Getting Started\nReview the [Getting Started Guide](http://localhost:3000/documentation/GettingStarted.md) for minimal setup instructions",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}
