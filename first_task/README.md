
### The task №1

#### Prerequisites
* it is necessary that the following ports in the system are not occupied, they are used for the following applications
    * 9443 - REST API
    * 3000 - Grafana
    * 9090 - Prometheus
    * 5432 - Postgres
    * 9187 - Postgres_exporter
    * 8080 - Cadvisor
    * 9100 - Node_exporter
    * 9001 - use only Jmeter Listener

### Task goal

* It is necessary to load test the REST API of the Technical_interview service and provide the test results.
* When using jmeter with the specified plugin, the results will be collected automatically, if any, or do not delete the volumes data when the project is stopped using
``` sh
docker-compose down 
```
OR
``` sh
make stop
```
* for other load tools, it is enough to provide screenshots of charts. Displaying charts in Prometheus will be a plus


1. Get project.
2. Run project (credential from grafana admin\admin)
3. Сreate a script for [Apache Jmeter](https://apache-mirror.rbc.ru/pub/apache//jmeter/binaries/apache-jmeter-5.4.1.tgz). You need to use [jmeter-prometheus-listener-plugin](https://github.com/kolesnikovm/jmeter-prometheus-listener/releases/download/2.0.2/jmeter-prometheus-listener-plugin-2.0.2.jar) for pass metrics. You can also use MF LoadRunner, Gatling and others load tools. In this case, the provision of graphs of the load of these instruments will be required or use library prometheus for these tools and add target to the file ./reconfigure/prometheus/prometheus.yml.
    * Scenario for one user:
        * user creation
            * post request by "/api/v1/users/createrandomuser" with empty body.
        * cycle of operations (quantity is calculated based on requirements)
            * Get All users
                * get request by "/api/v1/users/getusers" and extract a random user from the response, for update in the next request.
            * Updating extracts user
                * post request by "/api/v1/users/updateuser" with json body, the structure is similar to the received user structure in the request "Get All users".
                * Restrictions: field ID not is update.
            * Check updating user
                * get reuest by "/api/v1/users/getuser" with params ID equal id to users updated in the previous request. 
        * Delete extract user
                * post request by "/api/v1/users/deleteuser" with json body of the user structure, including the ID field (other fields are ignored).
4. Non-functional requirements:
    * 250 vusers in active
    * Request Create User - 10tps 
    * Requests Get All Users, Update and Check - 300tps
    * Request Delete User - 10tps
    * SLA for all request - 1s