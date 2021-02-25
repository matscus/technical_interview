![logo](https://github.com/matscus/technical_interview/blob/main/image/logo.png)

# technical_interview
## A project created to test candidates for load testing positions.

### Composition
Simple rest api service(is a target of testing) , postgresdb, prometheus + exporters, grafana

### How to use
#### Get and run
```sh
git clone https://github.com/matscus/technical_interview 
cd technical_interview
docker-compose up --build -d
```
#### Stop and remove data
``` sh
docker-compose down --rmi all -v
```

### The task №1
* [The task №1](https://github.com/matscus/technical_interview/TASCONE.md)

