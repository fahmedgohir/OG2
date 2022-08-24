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
   -d '{"user":{"name":"user2"}}'
```
