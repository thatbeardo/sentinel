export const creatingPermissions = `

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
`