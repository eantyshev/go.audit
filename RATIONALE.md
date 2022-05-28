## Specification and rationale

We need a scalable distributed logging service, and our basic requirements are:
* High insert rate, a big number of producers
* Shallow events schema
* Timings at producer cite are critical

Nice to have:
* No single point of failure
* Scalable in every dimension

### Data model and proposed API
Event:
* `id`
* `type`
* `consumer`
* `created_at`
* and a free-form `payload` (optional)

### Authentication
Assuming that our consumers are within some security perimeter, we choose the simplest approach: API key
Enforce setting the env variable AUDIT_API_KEY and expect the key in `X-Api-Key` header

### Storage
We need a flexible schema, multiple indexes for fast queries, and a high insert rate: MongoDB seems a good match.
We'll have a single database `dbevents` and a collection `events`
Must-have indexes: by created_at, consumer
Testing AIO environment will have a single MongoDB instance.

## Proposed API
POST /v1/event
```
{
    "type": "some event type",
    "consumer": "consumer1",
    "payload": {...}
}
```

GET /v1/find
request:
```
{
    // every field is optional
    "type": "event type 1",
    "consumer": "consumer 1",
    "created_from": "2022-05-22T00:00:00Z" // ISO8601 datetime
    "created_to": "2022-05-23T00:00:00Z" // ISO8601 datetime
}
```
response:
```
{
    "events": [
        {
            "id": "628e58fc08808bcd3b2af122",
            "type": "t2",
            "consumer": "c4",
            "created_at": "2022-05-25T16:27:40.004Z"
        },
        {
            "id": "628e592a08808bcd3b2af123",
            "type": "t2",
            "consumer": "c4",
            "created_at": "2022-05-25T16:28:26.23Z",
            "payload": {
                "t2_data": 12341243
            }
        }
    ]
}
```

## TODO
### MongoDB cluster
To get rid of a single POF for inserts we need at least 2 replicas: primary and secondary for writes and reads respectively.
Further, we may event introduce sharding by consumer ID hash

### Long-term storage
With MongoDB as a runtime storage, we should establish a synchronization procedure with a big-data storage like Hadoop, and introduce TTL of MongoDB events

### Consuming events via Message Queue
* Sometimes primary replica goes down, and the cluster doesn't serve write requests for ~10sec
* Even in normal operation producers don't want to wait for write confirmation.
Naturally, we need an MQ to store in-flight messages (events).
RabbitMQ is free, scalable, and configurable.
But to keep the prototype simple we implement an endpoint  `POST /v1/event` which synchronously saves a log record.
