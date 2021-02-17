export const contexts = `## Contexts

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
`