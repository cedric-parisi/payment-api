# payment-api

payment-api provides CRUD operation on a payment resource.

go-swagger documentation can be found [here](https://app.swaggerhub.com/apis-docs/cedric-parisi/payment-api/1.0.0#/payments/getPayment)

# installation

```
git clone https://github.com/cedric-parisi/payment-api.git
```
Once downloaded, go to the root directory and run
```
make install
``` 
It will install all the dependencies.


# local development

To help on the development, a docker-compose with a postgres database and a jaeger client are launched using: 
```
make dev
```

To stop and remove the docker images, launch:
```
make stop
```

## database migration

Running: 
```
make migration
```

will:
- Create or update schema according to models definition. 
- Insert mock data from `mock.json` (must be in the root directory)

## run the API

Use:
```
make run
```

## test

To get an HTML representation of the code coverage, use:
```
make cover
```

To launch all unit test, use:
```
make test
```

`mockery` is used to generate mocks over the interfaces, to update/create them, launch:
```
make update-mocks
```

Integration test can be launched using: 
```
make integration
```
<b> 
the current version of integration test requires to have a clean local db.
All these operations are required because the integration test are working with hardcoded id, meaning that after a launch some payments won't be present anymore. 

v2 will fetch a payment from the db before performing an delete/update action
<b>