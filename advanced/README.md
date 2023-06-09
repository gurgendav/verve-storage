# Verve storage redis implementation

## Services

### API
API service is an HTTP server that exposes 1 endpoint:

`/promotions/{id}` - returns promotion by id

### Worker (file-parser)
Worker service is a simple application that reads the csv file every 30 minutes and inserts the data into Redis.

## How to run

```bash
docker-compose up
```
API will be available at port `2020`
