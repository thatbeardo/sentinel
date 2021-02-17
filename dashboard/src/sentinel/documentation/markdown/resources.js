export const resources = `
## Resources

Simply stated a resource represents the nouns in your application. Resources are polymorphic and might represent users, accounts, or any application specific entity requiring gated access.

A principal is defined as a resource that is requesting access to another resource. The latter is the target. To generalize, Principal requests access to a Target. In other words, your application will often ask Sentinel if a Principal *can access* a particular Target(s). 

Resources are granted access to other resources through Contexts. It is possible for a resource to be the \`principal\` in once context while simultaneously a \`target\` in another context.

Resource may be owned by other resources. Resource inheritance will help reduce redundant permissions to similar resources if the principal has the same level of access to a class of resource. E.g. System Admins can read/write all accounts.

When creating a new Resource, you don't have to mention if it is a Principal or a Target. The terms Principal and Targets only become meaningful when granting a context to a Principal or permissions to a target.
`