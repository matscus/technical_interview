![logo](https://github.com/matscus/technical_interview/blob/main/images/logo.png)

# Technical interview
## A project created to test candidates for load testing positions.
![GitHub](https://img.shields.io/github/license/matscus/technical_interview?color=31E311)
![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/matscus/technical_interview)


### Composition
Simple api services, postgresdb, prometheus + exporters, grafana, rabbitmq.

### Work tested on operating systems 

* Mac Os Big Sur 11.2.2
* Arch Linux.
* Windows 10.

Before using, it is recommended to update the OS and Docker to the current versions.

### How to use

#### Get and run
```sh
git clone https://github.com/matscus/technical_interview 
cd technical_interview/{TASK DIR}
docker-compose up --build -d
```
#### Stop and remove data
``` sh
docker-compose down --rmi all -v
```

#### Alternatively you can use make commands
```sh
# Run unit test
make unittest

# Build binary
make engine

# Build docker image and run
make run

# Stop containers
make stop

# Stop containers, remove all image and volumes
make kill

# Remove binary file
make clean
```

### The tasks
* [The task №1](https://github.com/matscus/technical_interview/blob/main/first_task/README.md) ([RUS](https://github.com/matscus/technical_interview/blob/main/first_task/READMERUS.md))
