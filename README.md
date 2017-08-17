# What it does:
Using docker-compose we can bring up 3 separate containers: a Stats Server, the Workers and a Redis cache.

# Before getting started:
Make sure you have the following installed:
- docker
- docker-compose
- go

# Getting started:
```
make getdeps
```

# Build the container:
```
make build
```

# Running it:
```
make up
```

## Note: build and run with one commandL
```
make
```

# Cleaning up:
```
make clean
```

# -----------
# For initial development (without build the image):
```
docker-compose -f docker-compose-dev.yml up
```

# Control + C:
To terminate the running docker-compose process
