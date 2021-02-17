export const permission = `
## Permissions
A permission defines access to perform an action between a context and a Resource. Permissions take simply a \`name\` and whether the permission is \`permitted\`. The \`name\` is used to determine a principal's authorization to a target resource. 

Permission \`names\` can be any \`String\` that is meaningful to your application. It is recommended that you namespace your permissions with sufficient information to make it obvious what the permissions does.

Example Pattern:

\`\${org}:\${app}:\${resource}$:{action} => bithippie:sentinel:account:create\`

Permission Rules:

- A permission can have \`permitted: "allow"\` which grants all \`principals\` of the context access to the \`target_resource\` for that permission.
- A permission can have \`permitted: "deny"\` which prevents all \`principals\` of the context from performing the permission action the \`target_resource\` for that permission.

`