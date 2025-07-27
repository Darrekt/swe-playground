# Redis

A `docker-compose` file is provided to start redis along with an interfacing golang container and a `redis-cli` container. The `go` files provided provide both a server and client to spin up and play with.

- [Redis docs: running with Docker](https://redis.io/docs/latest/operate/oss_and_stack/install/install-stack/docker/)

The protobuf compiler `protoc` is included in the `devcontainer` specification. From this directory level, you can compile `.proto` files with

```bash
protoc --go_out=. --go_opt=paths=source_relative redis/src/leaderboard/proto/submission.proto
```

## Rate Limiter

- [Redis docs: Building a rate limiter](https://redis.io/learn/howtos/ratelimiting)

The rate limiter server implements a few basic rate-limiting algorithms:

- Sliding Window
- Fixed Window
- Token Bucket
- Leaky Bucket

The client allows you to configure how many requests to send per second to this rate limiter via the command line arguments. A Grafana dashboard will show a live view of the rate-limiter as well as the result of each request.

## Query Cache

- [Redis docs: Building a query cache](https://redis.io/learn/howtos/solutions/microservices/caching)

The query cache provides a mechanism to store the result sets of SQL queries into Redis for quick lookup. The server will provide the following interfaces:

- Start in write-through or cache-aside mode
- Provide the results to a query, along with whether the query was in cache or not.
- Invalidate a given key-value pair in the cache

## Leaderboard

- [Redis Docs: Sorted Sets](https://redis.io/docs/latest/develop/data-types/sorted-sets/)

The example we will use here is a LeetCode contest. Players are submitting entries to different questions within the contest. The client program will generate entries with a score between `1` and `100` and accumulate them to a given user's cumulative score. On Grafana, we should be able to see the generated entries, along with a live view of the redis leaderboard.

Use `main.go` with `leaderboard-client` and `leaderboard-server` command line arguments in two appropriately configured containers. You can then `exec` into the `client` container to use the client program and send in entries to the leaderboard.

```bash
docker exec -it client go run main.go leaderboard-client
```

To view the leaderboard for a given question ID using `redis-cli`:

```bash
docker exec -it redis redis-cli
> ZREVRANGE $QUESTION_ID 0 -1 withscores
```

For sanity checking, you can see the entries received by the server using `docker logs server`.
