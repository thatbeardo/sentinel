export const prerequisites = `
# Usage

## Prerequisites
First, you must retrieve an access token to make authenticated calls to Sentinel. You must have recieved a \`client_id\` and \`client_secret\`. If not or if you wish to rotate your secret in case it has been compromised, immediately contact \`contact@bithippie.com \` Make an http request as such to retrieve the access token: 

\`\`\`shell
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
`