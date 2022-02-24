# netra
> _Etymology. From Sanskrit नेत्र (netra, “eyes”)._

This is Yet Another Issue Tracking Service!

## Prerequisites
- docker
- docker-compose
- make
## Get Started
- `make build && make run`
- Access API at localhost:3000
- Access CockroachDB at localhost:8080

## Basic Design
- Server with Route handling and middleware is done using Go Chi
  - Simpler nested routing with provision of context building middlewares
  - Predefined middlewares like RequestID, Logger, Recoverer
- Handlers are built around a struct which talks to persistence layer structs
  - Abstracts out the DB calls for re-usability e.g., Context building allows fetching of an Issue from DB for Get, Delete (and Update if I had created one)
  - Allows dependency injection at instantiation level
- Database adapter is upper.io
  - Simpler to use and handy mechanisms of running ORM-like queries, SQL builders or plain SQL
  - Integrates with almost every major DB
- ~~Choice of DB~~ This entire section is so inaccurate as CockroachDB is not a typical SQL database. It's a robust Google Spanner-sequel that works like a geodistributed quorum based distributed database with SQL APIs. Internally, each node uses a NoSQL database engine. This should scale horizontally too, however, for this problem we will be underutilizing its complex join/transactional capabilities. CockroachDB is usually preferred for an even distribution of READ/WRITE traffic
  - ~~CockroachDB, a SQL database, was chosen for this dev project to showcase API design. In the usual case, each org with their set of users will be accessing a private set of issues. In this case, each org can be issued a SQL database~~
  - ~~Based on deployment strategy for the service, for example, if issue is a global concept for all users (like GitHub issues), it won't be scalable to deploy in a SQL database. A distributed NoSQL database like DynamoDB might be better~~
## TODO/Improvements
- Update API: I didn't have time to make it for any good use
- Docs: I didn't have time due to work commitments so I blazed past the code
- Better Logging: I wanted to use a log framework like zap or something to make debugging easier
- Better Errors: I scarcely used the errors package for wrapping up nests up to the user, this could be improved
- Traces: I wanted to setup traces to investigate deeper into API performance of Search
- Search: This was minimal, most time was spent on dockerizing the app and making sure core API works. I attempted a simple Bag of Words-based Full Text Search using idea from https://www.cockroachlabs.com/blog/full-text-indexing-search/ but it took too much time
- Auth: There are Basic, JWT and OAuth2-based middleware available to setup Auth for the app, however, implementing any of them while actively developing on core functionality was slowing down development. For example, I kept getting prompted for Username and Password when using Basic auth
- Simplified Model: Having separate models for persistence and service layer made some things clean but at the same time converting one struct to another was a repeated chore. I want to investigate the best standards around this
- Testing Server & Router: I wanted to write integration tests for server and routing but ran out of time
- CLI: Using Cobra for defining DB config to distinguish Docker instance vs local run would have been nice
