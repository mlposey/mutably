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
* database - a PostgreSQL database for storing word data
* anvil    - a tool exploring and importing archives
* api      - a RESTful API that provides unified access to core service logic

## Starting The Service
Before starting the service, you should install the following: 
1. Docker
2. Docker Compose
3. wget
4. lbzip2

Then carry out the following two steps: 
1. Run `get-archive.sh` from the `archive` folder. This downloads a Wiktionary
archive for parsing. The initial download is ~700M and the decompressed version
is ~6G. The service only needs the decompressed file.
2. Run `docker-compose up` in the root project directory, passing in the required
environment variables. See [the docker-compose file](./docker-compose.yaml) for
a list of required variables.

## Development Pipeline
It is important that the main branch stays production ready. This goal is
accomplished by only introducing changes through PRs--which are only accepted
after passing all build checks. The build checks act as a measure of quality,
so it is important that they stay relevent in order to guage the quality of
current code. We can accomplish currency by ensuring all pull requests come
with tests for units introduced and the integration of those units with the
existing codebase.
