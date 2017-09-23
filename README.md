# Mutably
![Build Status](http://teamcity.marcusposey.com/app/rest/builds/aggregated/strob:(buildType:(project:(id:Mutably)))/statusIcon.svg)

Mutably is a tool for learning natural languages. It currently offers partial
support for Dutch verb conjugation.

## Roadmap
Development is focused on verb conjugation with the hope of bringing this task
to voice platforms. To accomplish this, we need:

1. A dataset that forms the basis of conjugation
2. General methods to extract relational inflection data from the set
3. A common API for the front-end applications
4. An initial target platform for testing voice interfaces

## Project Structure
The project has three components for now:
* database
* anvil
* api
`database` uses Docker and PostgreSQL to define the production database environment.
You can refer to the .sql schema files there for details on how data is managed
by the project. The next componenet, `anvil`, is a tool for sifting through
datasets and using them to build the application's relational data model. With
cleaned data in the database, `api` creates a front-end REST layer that any
platform can use to interfact with the core logic.

## Development Pipeline
It is important that the main branch stays production ready. This goal is
accomplished by only introducing changes through PRs--which are only accepted
after passing all build checks. The build checks act as a measure of quality,
so it is important that they stay relevent in order to guage the quality of
current code. We can accomplish currency by ensuring all pull requests come
with tests for units introduced and the integration of those units with the
existing codebase.
