# promova-test-task

Startup instructions:

### Project assembly
To build and run REST microservice, run `make server`. Pay attention that the instance of PostgreSQL database running on port 5432 with database called `promova_test_task` should be created.

### Local deployment
To have instances of PostgreSQL database and REST microservice as containers locally, run `make deploy`. This command will utilize `docker-compose.yaml` to raise instances of DB and REST service.
