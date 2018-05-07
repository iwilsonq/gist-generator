What is a good way to go about testing my /graphql endpoint for all of my different queries and mutations?

```javascript
export default (
	{ id }: DBChannel,
	{ first, after }: PaginationOptions,
	{ loaders }: GraphQLContext
) => {
	const cursor = decode(after)

	const lastDigits = cursor.match(/-(\d+)$/)
	const lastUserIndex = lastDigits && lastDigits.length > 0 && parseInt(lastDigits[1], 10)

	return getMembersInChannel(id, { first, after: lastUserIndex })
		.then(users => loaders.user.loadMany(users))
		.then(result => ({
			pageInfo: {
				hasNextPage: result && result.length >= first
			},
			edges: result.filter(Boolean).map((user, index) => ({
				cursor: encode(`${user.id}-${lastUserIndex + index + 1}`),
				node: user
			}))
		}))
}
```

Personal access tokens function like ordinary OAuth access tokens. They can be used instead of a password for Git over HTTPS, or can be used to authenticate to the API over Basic Authentication.

```javascript
function sayHello() {
	return function(name) {
		console.log('Hello ', name)
	}
}
```

Personal access tokens function like ordinary OAuth access tokens. They can be used instead of a password for Git over HTTPS, or can be used to authenticate to the API over Basic Authentication.
