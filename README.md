[![License MIT](https://img.shields.io/badge/License-MIT-brightgreen.svg)](https://img.shields.io/badge/License-MIT-brightgreen.svg)

# go.audit

Trivial audit log service, providing API to store and retrieve audit events

## Local deployment
Execute `make build` to build docker images
Execute `make run` and the service will be brought up using `docker-compose.yml`
REST API will be listening at `localhost:8080`
## Usage

* Post the new audit event (using the default secret from ):
> curl -v -X POST -H "X-Api-Key: supersecret1"  -H "Content-Type: application/json" http://localhost:8080/v1/event -d '{"type": "t1", "consumer": "c1"}'
* Post another with custom payload
> curl -v -X POST -H "X-Api-Key: supersecret1"  -H "Content-Type: application/json" http://localhost:8080/v1/event -d '{"type": "t2", "consumer": "c1", "payload": {"t2_param": 123}}'
* List all events
> curl -v -X GET -H "X-Api-Key: supersecret1" -H "Content-Type: application/json" http://localhost:8080/v1/events -d '{}'
* List by type, consumer, created_from/created_to range (any of these is optional)
> curl -v -X GET -H "X-Api-Key: supersecret1" -H "Content-Type: application/json" http://localhost:8080/v1/events -d '{"type": "t2", "created_from": "2022-05-25T16:10:00Z", "consumer": "c1"}'
```
{
    "events": [
        {
            "id": "628f4fe793347802739d3a03",
            "type": "t2",
            "consumer": "c1",
            "created_at": "2022-05-26T10:01:11.923Z",
            "payload": {
                "t2_param": 123
            }
        }
    ]
}
```
__note__ that `id` and `created_at` fields appeared
## Unit tests

```make unittest```
Handlers are tested using in-memory storage

## Integration tests

Run `make test` for integration tests which run into separate container (defined in `docker-compose.test.yml`)
Testing environment will be deployed and the tests will be executed, see the exit code for result