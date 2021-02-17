export const creatingContexts = `
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
`