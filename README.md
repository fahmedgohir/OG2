# og2

## Run with:
```
docker compose up
```

## Rebuild with:
```
docker compose build og2
```

## Requests:

### Create a session:
```
curl -X POST localhost:8081/user \
   -H 'Content-Type: application/json' \
   -d '{"user":{"name":"user1"}}'
```

### Query a session:
```
curl -X GET localhost:8081/dashboard \
   -H 'Content-Type: application/json' \
   -d '{"user":{"name":"user1"}}'
```

## Output:
```
{"user":{"name":"user1"},"resources":{"iron":410,"copper":123,"gold":82},"factories":{"iron_factory":{"level":1,"resource":"iron"},"copper_factory":{"level":1,"resource":"copper"},"gold_Factory":{"level":1,"resource":"gold"}},"last_updated":1661355401}
```
