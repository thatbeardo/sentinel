const docs = `
## Getting Started
We are thrilled to have you onboard. This guide allows you to understand how to leverage Sentinel's capabilities and implement Authorization like a boss! 😎 

- [The Basics](#the-basics)
  - [Resources](#resources)
  - [Contexts](#contexts)
  - [Permissions](#permissions)
  - [Definitions](#definitions)
- [Usage](#usage)
  - [Prerequisites](#prerequisites)
  - [Use Case](#use-case)
  - [Creating Resources](#creating-resources)
  - [Creating Contexts](#creating-contexts)
  - [Creating Permissions](#creating-permissions)
  - [Authorization](#authorization)

# The Basics

## Resources

Simply stated a resource represents the nouns in your application. Resources are polymorphic and might represent users, accounts, or any application specific entity requiring gated access.

A principal is defined as a resource that is requesting access to another resource. The latter is the target. To generalize, Principal requests access to a Target. In other words, your application will often ask Sentinel if a Principal *can access* a particular Target(s). 

Resources are granted access to other resources through Contexts. It is possible for a resource to be the \`principal\` in once context while simultaneously a \`target\` in another context.

Resource may be owned by other resources. Resource inheritance will help reduce redundant permissions to similar resources if the principal has the same level of access to a class of resource. E.g. System Admins can read/write all accounts.

When creating a new Resource, you don't have to mention if it is a Principal or a Target. The terms Principal and Targets only become meaningful when granting a context to a Principal or permissions to a target.

## Contexts

Contexts \`allow\` or \`deny\` permissions to principals on target resources. A single context can be granted to zero or more principals and affect zero or more target resources.

A context contains zero or more permissions to target resources. A permission is always scoped to a single target resource within the context. However, it is possible to grant the same permission to multiple target resources, or grant multiple permissions to one target resource within the scope of a single context.

| Principal | Context | Permission [org:account:resource:action]  | Target Resource | Owned By  | Notes                                                                                                                |
|-----------|---------|-------------------------------------------|-----------------|-----------|----------------------------------------------------------------------------------------------------------------------|
| User1     | Default | bithippie:sentinel:account:read {allow}   | Resource1       | Resource0 | Allow User1 to read Resource1                                                                                        |
| User1     | Admin   | bithippie:sentinel:account:read {allow}   | Resource1       | Resource0 | Allow User1 to assume \`Admin\` context and read Resource1                                                             |
| User1     | Admin   | bithippie:sentinel:account:write {allow}  | Resource1       | Resource0 | Allow User1 to assume \`Admin\` context and write Resource1                                                            |
| User1     | Admin   | bithippie:sentinel:account:read {allow}   | Resource2       | Resource0 | Allow User1 to assume \`Admin\` context and read Resource2. This does not give User1 permissions to write to Resource2 |
| User2     | Default | bithippie:sentinel:account:read {allow}   | Resource2       | Resource0 | Allow User2 to read Resource2                                                                                        |
| User3     | Super   | bithippie:sentinel:account:read {allow}   | Resource0       | null      | Allow User3 to read all Resources owned by Resource0                                                                 |
| User3     | Super   | bithippie:sentinel:account:write {allow}  | Resource0       | null      | Allow User3 to write all Resources owned by Resource0                                                                |
| User3     | Super   | bithippie:sentinel:account:delete {allow} | Resource0       | null      | Allow User3 to delete all Resources owned by Resource0                                                               |

Each grant of a context to a principal contains a \`with_grant\` boolean allowing the aforementioned principal to propagate their permissions to additional principals if the \`with_grant\` option has been set to \`true\`. By default \`with_grant=false\`

## Permissions

A permission defines access to perform an action between a context and a Resource. Permissions take simply a \`name\` and whether the permission is \`permitted\`. The \`name\` is used to determine a principal's authorization to a target resource. 

Permission \`names\` can be any \`String\` that is meaningful to your application. It is recommended that you namespace your permissions with sufficient information to make it obvious what the permissions does.

Example Pattern:

\`\${org}:\${app}:\${resource}$:{action} => bithippie:sentinel:account:create\`

Permission Rules:

- A permission can have \`permitted: "allow"\` which grants all \`principals\` of the context access to the \`target_resource\` for that permission.
- A permission can have \`permitted: "deny"\` which prevents all \`principals\` of the context from performing the permission action the \`target_resource\` for that permission.

## Definitions

| Name        | Definition           |
| ------------- |:-------------:|
| Resource     | Any noun in your application. Resources (principals) access other Resources (target resource) through contexts. |
| Context      | A named collection of permissions. Contexts are granted to resources (principals) and grant or deny access to affect resources (target resources) through permissions.      |
| Permissions | An explicit allowance or refusal of a resource (principal) to perform an action on a resource (target resource).      |
| Principal | In the scope of a context the principal refers to the resource which has been granted access to make use of the context.      |
| Target | In the scope of a context a target resource is any and all resources which the context allows or denies access through Permissions.      |

&nbsp;
 
# Usage

## Prerequisites
First, you must retrieve an access token to make authenticated calls to Sentinel. You must have recieved a \`client_id\` and \`client_secret\`. If not or if you wish to rotate your secret in case it has been compromised, immediately contact \`contact@bithippie.com \` Make an http request as such to retrieve the access token: 

\`\`\`
curl --request POST \
  --url https://bithippie.auth0.com/oauth/token \
  --header 'content-type: application/json' \
  --data '{"client_id":"YOUR_CLIENT_ID_HERE","client_secret":"YOUR_CLIENT_SECRET_HERE","audience":"https://api.guardmy.app/","grant_type":"client_credentials"}'
\`\`\`

The response you will recieve will be of the following format: 

\`\`\`JSON
{
  "access_token": "a-really-long-string",
  "token_type": "Bearer"
}
\`\`\`

For all consequent reqeusts to Sentinel, you will embed this \`access_token\` in the header. We will cover the semantics of a call in following sections. For now, remember to cache this token for 24 hours and make as few requests as possible to generate the \`access_token\`. We are working on creating an endpoint that you can use to get the token without worrying about the caching. For this manual, we will also set the \`TENANT\` and \`ACCESS_TOKEN\` in our shell and use them throughout. 

\`\`\`shell
  export TENANT=YOUR_ENVIRONMENT
  export ACCESS_TOKEN=YOUR_ACCESS_TOKEN
\`\`\`

## Use Case

For this manual, we will setup all the Sentinel infrastructure keeping in mind this following use case. 

\`\`\`
Danny is a Floor Manager at an e-commerce company. Danny has 2 employees reporting to him directly - Rusty, and Linus. Danny and his team use a portal to determine salaries. Given the hierarchy, as a manager, Danny should be able to edit salries and modify it. However, as an employee reporting to Danny, I should not be able to edit salaries. Danny and their team use MoneyManager – a payroll system developed in house. MoneyManager is written in React and the entire team logs into the same portal. When MoneyManager loads up, it consults Sentinel to get information about the currently logged in user and displays suitable UI components based on the role of who is logged in. For this example, we will create a bunch of resources, contexts and permissions to make sure Sentinel knows that Danny has privileges to edit salaries but others don't. 
\`\`\`

## Creating Resources

As seen previously, \`resources\` are nouns in your applications. In the example above this means that each employee is a \`resource\` in the Sentinel realm. Further, since we need gated access to salaries, we make a salary resource too. We must first go ahead and create all 3 employees and the salary resource. Use the following curl commands to create Employees Danny, Rusty, and Linus. 

\`\`\`
curl -X POST "https://sentinelbeta.herokuapp.com/v1/resources/" -H "accept: application/json" -H "x-sentinel-tenant: $TENANT" -H "Authorization: bearer  $ACCESS_TOKEN" -H "Content-Type: application/json" -d "{ "data": { "attributes": { "name": "danny", "source_id": "danny@ecommerce.com" }, "type": "resource" }}"
\`\`\`

The result of creating a new resource looks something like this:

\`\`\`JSON
{
  "data": {
    "id": "528961a8-56e2-4cfa-8825-0ff83036a333",
    "type": "resource",
    "attributes": {
      "name": "danny",
      "source_id": "danny@ecommerce.com",
      "context_id": ""
    },
    "relationships": {
      "parent": {
        "data": {
          "type": "resource",
          "id": "3d17e5eb-1fb3-4a23-8bb0-e374d6184659"
        }
      }
    }
  }
}
\`\`\`
Take a note the \`id\` field. This is a UUID generated by Sentinel to track Danny. We will use this when creating contexts in the second step. 

\`\`\`
curl -X POST "https://sentinelbeta.herokuapp.com/v1/resources/" -H "accept: application/json" -H "x-sentinel-tenant: $TENANT" -H "Authorization: bearer  $ACCESS_TOKEN" -H "Content-Type: application/json" -d "{ "data": { "attributes": { "name": "rusty", "source_id": "rusty@ecommerce.com" }, "type": "resource" }}"
\`\`\`

\`\`\`
curl -X POST "https://sentinelbeta.herokuapp.com/v1/resources/" -H "accept: application/json" -H "x-sentinel-tenant: $TENANT" -H "Authorization: bearer  $ACCESS_TOKEN" -H "Content-Type: application/json" -d "{ "data": { "attributes": { "name": "linus", "source_id": "linus@ecommerce.com" }, "type": "resource" }}"
\`\`\`

Alternatively, you can use the swagger page found [here](https://sentinelbeta.herokuapp.com/swagger/index.html). If you decide to use the swagger page, remember to set the authorization header. You can do this by clicking on the \`Authorize\` button and entering \`bearer YOUR_ACCESS_TOKEN\` in the value field. 

Great! Now Sentinel tracks these 3 newly created employees for you. Our use case dictates that Danny should have \`read\` and \`write\` permissions to the salaries while Rusty and Linus only have \`read\` permission set. Before we can set these permissions, we need to create the salary resource. Let's do that. 

\`\`\`
curl -X POST "https://sentinelbeta.herokuapp.com/v1/resources/" -H "accept: application/json" -H "x-sentinel-tenant: $TENANT" -H "Authorization: bearer  $ACCESS_TOKEN" -H "Content-Type: application/json" -d "{ "data": { "attributes": { "name": "salary", "source_id": "salary-source" }, "type": "resource" }}"
\`\`\`

## Creating Contexts

A context may \`allow\` or \`deny\` permissions to Principals on Target resources. A \`context\` cannot exist without a \`resource\`. This is why the endpoint to create a \`context\` asks for an \`id\`. This \`id\` is the unique identifier for the resource to which you want to attach the context to. For our example, it makes sense and is intuitive to have one separate context for Danny who is the manager. And another context that can be shared by Rusty and Linus. This means we need to create two contexts and associate them with suitable resources. 

First let's create the "managerial" context. This will be attached to Danny. We need Danny's unique identifier before we can create the context. You can use the \`GET /v1/resources\` endpoint to fetch all resources and get Danny's id. Now, make the following call to create a new \`context\` named manager and associate it with Danny

\`\`\`
curl -X POST "https://sentinelbeta.herokuapp.com/v1/resources/{DANNY_ID}/contexts" -H "accept: application/json" -H "x-sentinel-tenant: $TENANT" -H "Authorization: bearer $ACCESS_TOKEN" -H "Content-Type: application/json" -d "{ "data": { "attributes": { "name": "manager" }, "type": "context" }}"
\`\`\`

Next, we create a context that will be shared by both Rusty and Linus. We can first make a similar call as seen for Danny for either one of them and then attach that newly created context to the other person. Let's begin with Linus. We first create a context in association with Linus just like above. 

\`\`\`
curl -X POST "https://sentinelbeta.herokuapp.com/v1/resources/{LINUS_ID}/contexts" -H "accept: application/json" -H "x-sentinel-tenant: $TENANT" -H "Authorization: bearer $ACCESS_TOKEN" -H "Content-Type: application/json" -d "{ "data": { "attributes": { "name": "employee" }, "type": "context" }}"
\`\`\`

Extract the context \`id\` from the response. We will use this \`id\` to grant the newly created context to our second employee - Rusty. Check out the [swagger page](https://sentinelbeta.herokuapp.com/swagger/index.html) and in particular look at the Grants endpoint. One such endpoint is the \`PUT /v1/grants/resources/{resource_id}/contexts/{context_id}\` endpoint. We will use this to \`GRANT\` an existing context to a Resource, in our case, Rusty. 

\`\`\`
curl -X PUT "https://sentinelbeta.herokuapp.com/v1/grants/resources/{RUSTY_ID}/contexts/{EMPLOYEE_CONTEXT_ID}" -H "accept: application/json" -H "x-sentinel-tenant: $TENANT" -H "Authorization: bearer $ACCESS_TOKEN" -H "Content-Type: application/json" -d "{ "data": { "attributes": { "with_grant": true }, "type": "grant" }}"
\`\`\`

After making this call, we now have two contexts. One, a Manager context, granted to Danny and the second, an Employee context, granted to both Linus and Rusty. We are all set to start setting permissions

## Creating Permissions
Permissions are set on contexts. And contexts, as seen so far, are granted to resources. Permissions answer the question "Does the principal have access to a target". Thus Permissions \`allow\` or \`deny\` to principals on a particular target. In our example it is very clear that all employees will become \`Principals\` while salary becomes the \`Target\`. In other words, an intuitive question to ask Sentinel is "What sort of permissions does Danny have on salaries?". Note how it becomes obvious that each question is specific to one Principal - we will see this again when we encounter the \`authorization\` endpoint. 

For now, we can begin with setting permissions to the contexts we have created. Let's begin with the Manager context. Our case study dictates that a manager must be able to read and edit salaries. Let's go ahead and implement this in our model. We want to \`allow\` \`read\` and \`allow\` \`write\` permissions on the \`salary\` resource. This means we are telling Sentinel that anyone who is granted the Manager context is "allowed" to "read" salary and "write" "salary". Recall how the Manager context is granted to Danny. This means after these calls we will have successfully implemented a model where Danny as a Manager can \`read\` and \`write\` salaries. You can name a permission as you wish - a good habit is to name them using an idiomatic napspaces as described in the earlier sections. 
Assuming your company is called "enginerd" and you are developing permissions for your client "rolce", we chose the following names \`enginerd:rolce:read\` and \`enginerd:rolce:write\`. Here's a call to \`allow\` \`enginerd:rolce:read\` permissions on salary for the Manager context

\`\`\`
curl -X PUT "https://sentinelbeta.herokuapp.com/v1/permissions/{MANAGER_CONTEXT_ID}/resources/{SALARY_RESOURCE_ID}" -H "accept: application/json" -H "x-sentinel-tenant: $TENANT" -H "Authorization: bearer $ACCESS_TOKEN" -H "Content-Type: application/json" -d "{ "data": { "attributes": { "name": "enginerd:rolce:read", "permitted": "allow" }, "type": "permission" }}"
\`\`\`

Similarly, let us \`allow\` \`enginerd:rolce:write\` permissions on salary resource for the Manager context. 

\`\`\`
curl -X PUT "https://sentinelbeta.herokuapp.com/v1/permissions/{MANAGER_CONTEXT_ID}/resources/{SALARY_RESOURCE_ID}" -H "accept: application/json" -H "x-sentinel-tenant: $TENANT" -H "Authorization: bearer $ACCESS_TOKEN" -H "Content-Type: application/json" -d "{ "data": { "attributes": { "name": "enginerd:rolce:write", "permitted": "allow" }, "type": "permission" }}"
\`\`\`

With these permissions set, Danny now has \`enginerd:rolce:read\` and \`enginerd:rolce:write\` access to salary resource. Let's make the same call on Employee context. This time we only want \`allow\` \`enginerd:rolce:read\` permissions. The call would look like this: 

\`\`\`
curl -X PUT "https://sentinelbeta.herokuapp.com/v1/permissions/{EMPLOYEE_CONTEXT_ID}/resources/{SALARY_RESOURCE_ID}" -H "accept: application/json" -H "x-sentinel-tenant: $TENANT" -H "Authorization: bearer $ACCESS_TOKEN" -H "Content-Type: application/json" -d "{ "data": { "attributes": { "name": "enginerd:rolce:read", "permitted": "allow" }, "type": "permission" }}"
\`\`\`

## Authorization

With that last call, we have our model set. Sentinel now has some data to answer authorization questions. To sum it up, we have 4 resources - Danny, Rusty, Linus, and Salary. We have two contexts. First, a Manager context, granted to Danny and second, an employee context granted to both Rusty and Linus. The Manager context allows the principal - in this case Danny - to "read" and "write" on the Target Resource - in this case Salary. Similarly, The Employee context allows the principal - in this case Linus and Rusty - to only "read" on the Target resource viz Salary. Let's take a look at how can we ask Sentinel questions about authorization. 

We begin with a basic question. Imagine in our use case, a front-end application is loading up and it has realized that Danny is logging into the system. An interesting question to ask would be "List all permissions that are granted to Danny to any Target resources." The following curl command does exactly this

\`\`\`
curl -X GET "https://sentinelbeta.herokuapp.com/v1/principal/{DANNY_ID}/authorization?context_id={MANAGER_CONTEXT_ID}&depth=0&include_denied=false" -H "accept: application/json" -H "x-sentinel-tenant: $TENANT" -H "Authorization: bearer $ACCESS_TOKEN"
\`\`\`

Notice how we have to mention the \`context_id\` with which Danny is trying to log in. Of course, this becomes cumbersome. To avoid embedding the \`context_id\` this way in every call, we can instruct Sentinel to store the Manager context as Dannys default context. This can be achieved by the following call. 

\`\`\`
curl -X PATCH "https://sentinelbeta.herokuapp.com/v1/resources/{DANNY_ID}" -H "accept: application/json" -H "x-sentinel-tenant: $TENANT" -H "Authorization: bearer $ACCESS_TOKEN" -H "Content-Type: application/json" -d "{ "data": { "attributes": { "context_id": "57adee56-9fcf-4194-be43-da05bdcf329e", "source_id": "danny@ecommerce.com" }, "type": "resource" }}"
\`\`\`

By making this \`PATCH\` call we have set the default context id for Danny. It is a way of telling Sentinel that if no \`context_id\` is provided in the \`authorization\` call (the call before the last one) then use the Manager context ID. Let's try the previous authorization call again but this time, without providing the context id. 

\`\`\`
curl -X GET "https://sentinelbeta.herokuapp.com/v1/principal/{DANNY_ID}/authorization?depth=0&include_denied=false" -H "accept: application/json" -H "x-sentinel-tenant: test-one" -H "Authorization: bearer $ACCESS_TOKEN"
\`\`\`

The result in this case and in the previous authorization call will be 

\`\`\`json
{
  "data": [
    {
      "id": "5f874a26-a4ec-49a7-ab87-b48ef5fb4197",
      "type": "resource",
      "attributes": {
        "name": "salary",
        "source_id": "salary-source",
        "context_id": ""
      },
      "relationships": {
        "permissions": {
          "data": [
            {
              "name": "enginerd:rolce:write",
              "permitted": "allow"
            },
            {
              "name": "enginerd:rolce:read",
              "permitted": "allow"
            }
          ]
        }
      }
    }
  ]
}
\`\`\`

This response is strictly about the Principal we provided in the request. In this case, it is about Danny. Sentinel is telling us that Danny has \`enginerd:rolce:write\` and \`enginerd:rolce:read\` permissions for Target resource \`salary\`.

Let's make the same call for Linus. Recall how Linus does not have the default \`context_id\` set. So we will have to provide that in the authorization request. The request to figure out Linus' permissions is as follows

\`\`\`
curl -X GET "https://sentinelbeta.herokuapp.com/v1/principal/{LINUS_ID}/authorization?context_id={EMPLOYEE_CONTEXT_ID}&depth=0&include_denied=false" -H "accept: application/json" -H "x-sentinel-tenant: $TENANT" -H "Authorization: bearer $ACCESS_TOKEN"
\`\`\`

and the response like before will look like this

\`\`\`json
{
  "data": [
    {
      "id": "5f874a26-a4ec-49a7-ab87-b48ef5fb4197",
      "type": "resource",
      "attributes": {
        "name": "salary",
        "source_id": "salary-source",
        "context_id": ""
      },
      "relationships": {
        "permissions": {
          "data": [
            {
              "name": "enginerd:rolce:read",
              "permitted": "allow"
            }
          ]
        }
      }
    }
  ]
}
\`\`\`

This response is strictly about the Principal we provided in the request. In this case, it is about Linus. Sentinel is telling us that Linus has only \`enginerd:rolce:read\`  permissions for Target resource \`salary\`. With this, we have a basic idea about storing permissions with Sentinel. Next up, let's look at some advanced use cases.`;

export default docs;
