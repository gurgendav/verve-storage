# Verve storage

## Description

In this repository you can find 2 implementations of a [storage service](Case%20Study%20Golang.pdf). The [first one](simple) is a simple in-memory storage, the [second one](advanced) is a storage based on Redis.

### [Simple implementation](simple)

The simple implementation is an all-in-one solution keeping all the data in memory.

### [Advanced implementation](advanced)

The advanced implementation is trying to solve scaling problem by using Redis as a storage and separating the service 
into 2 parts: the main API service and the worker reading the file and inserting the data into the Redis.

API service can be scaled horizontally by running multiple instances of it.



