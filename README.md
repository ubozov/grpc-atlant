# grpc-atlant

Used tools:
- golang
- grpc
- protobuf
- mongodb
- nginx
- mongo-express
- docker
- docker-compose

## Build

### Deleloping
Envirolment variables are avialable on [.env.dev](.env.dev) file

The command below launches the `mongodb`, `mongo-express` on docker container and grpc-server at `localhost:10000`:
```
make run
```
The `mongo-epxress` is avialable at `localhost:8081`.

### Deployment
Envirolment variables are avialable on [.env](.env) file.
The `nginx` configuration file available on `conf/nginx.conf`

Below command deploys and runs all environment on the local docker host:
```
make deploy
```

Set SCALE parameter (by default is 2) to scaling grpc-server instance. For example:
```
make deploy SCALE=5
```
The `mongo-express` is avialable at `localhost:8081`.

The `nginx` available on the `localhost:1000` and upstreams all traffic between deployed grpc-servers instances.

## Usage
### Pagination
| Parameter | Description | Default |
| --- | --- | --- |
| `token` | *columnName_value_RecordID* (for example `price_34.1_6042fac29124f581b23467e2`) | empty
| `limit` | sets the limit of records per page | returns all records
### Sorting
| Parameter | Description | Default |
| --- | --- | --- |
| `columnName` | **_id**, **name**, **price**, **counter** or **lastModified** | **_id**
| `orderType` | **ASC** or **DESC** | **ASC**

