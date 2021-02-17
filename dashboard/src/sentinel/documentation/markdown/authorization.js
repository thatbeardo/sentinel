export const authorization = `

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

This response is strictly about the Principal we provided in the request. In this case, it is about Linus. Sentinel is telling us that Linus has only \`enginerd:rolce:read\`  permissions for Target resource \`salary\`. With this, we have a basic idea about storing permissions with Sentinel. Next up, let's look at some advanced use cases.`