# API
## Documentation
The REST API is documented with Swagger. Files for working with Swagger UI and
Swagger Editor can be found in `./swagger`. They should remain in sync with
each other.
## Testing
TeamCity uses `docker-compose.test.yaml` to set up an environment in which to
run `go test`. This method keeps the API in sync with the production database
model defined in the (project) root database directory.
