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

### integration test

Currently, integration tests are a bit tricky to use, please follow the steps below:

1/ `make stop` to stop the docker-compose and clean the images (in case you performed actions on the DB already)

2/ `make dev` to launch development dependencies only (The postgres DB uses the port 5432, the DB is ready when you can see `database system is ready to accept connections`)

3/ in a 2nd bash window: 
- `make migration` to insert mocked data in the DB
- `make run` to launch the API

4/ in a 3rd bash window, launch `make integration`