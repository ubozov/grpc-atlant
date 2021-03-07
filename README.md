# grpc-atlant

## Build

### Deleloping
Envirolment variables avialable on [.env.dev](.env.dev) file

For run on developer mode execute below command:
```
make run
```
MongoDB and MongoDB-Express run on the docker containers.

### Deployment
Envirolment variables avialable on [.env](.env) file.
Ngninx configuration file avialable on `conf/nginx.conf`

Below command deploys and runs all environment on the local docker host:
```
make deploy
```

Set SCALE parameter (by default is 2) to scaling grpc server instance. For example:
```
make deploy SCALE=5
```

## Usage
### Pagination
| Command | Description | Default |
| --- | --- | --- |
| `token` | *columnName_value_RecordID* (for example `price_34.1_6042fac29124f581b23467e2`) | empty
| `limit` | sets the limit of records per page | returns all records
### Sorting
| Command | Description | Default |
| --- | --- | --- |
| `columnName` | **_id**, **name**, **price**, **counter** or **lastModified** | **_id**
| `orderType` | **ASC** or **DESC** | **ASC**





