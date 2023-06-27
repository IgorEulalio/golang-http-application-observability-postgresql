## Architecture

![image](https://github.com/vmfarms/golang-sample-http-app/assets/41654187/3780ba0b-3216-4228-b0ab-bfe5f5e6e1df)

Inside this architecture, we can see the following components:
- Repository Service: Our main service is instrumented with OpenTelemetry and configured to send metrics and traces to the OpenTelemetry Collector. It also integrates with another service called Configuration Service, which. is also instrumented so we can follow the trace entirely.
- Configuration Service: Another service that returns the **configurationId**, which the repository service needs when creating new repositories in the database.
- OpenTelemetry Collector: This collector receives data from the services, processes it, and then sends it to configured monitoring backends.
- Monitoring backends: Just showing that we can have multiple other technologies being used as backends: Zipkin, Graphite, XRay....
- Jaeger: Tool used for distributed tracing.
- Prometheus: Tool used for metrics aggregation. Please note that the application uses a vendor-agnostic protocol, which means that OpenTelemetry offers a common library that can be used in multiple languages and does not have any lock with the metrics backend itself. 
## Local Setup

### Requirements
- [Golang](https://go.dev/doc/install) installed and GOPATH set version 1.16+
- [Docker](https://docs.docker.com/engine/install/ubuntu/)
- [jq](https://howtoinstall.co/en/jq)

### Running application
After you've installed everything, you can spin up the infrastructure and the configuration microservice, which will be used to simulate an external microservice.

First, clone this repository inside your pre-configured GOPATH, access it, and then run the following commands:
```
export PROJECT_HOME=$(pwd)
cd local-setup/
docker-compose up -d
```

Wait for all the containers to be ready, then run your application:
```
cd $PROJECT_HOME
go run ./cmd/repositories-service/main.go
```

### Testing endpoints
With the application running properly, you can test the endpoints available:
1. Create repository:
```
curl --request POST \                                                                                          at  16:18:55
  --url http://localhost:8080/repositories \
  --header 'Content-Type: application/json' \
  --data '{
        "name": "my_super_cool_repository",
        "owner": "igor",
        "configurationId": "10"
}' | jq '.ID'
```

2. Get All Repositories:
```
curl --request GET \
  --url http://localhost:8080/repositories
```

3. Get Repository by ID. Replace with the repositoryId returned by create repository:
```
curl --request GET \
  --url http://localhost:8080/repositories/{reposotiry_id}
```

4. Delete repository. Replace with the repositoryId returned by create repository:
```
curl --request DELETE \
  --url http://localhost:8080/repositories/{reposotiry_id}
```




