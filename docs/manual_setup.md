#Manual Setup

##Environment Variables

- REDIS_URL, typically `localhost:6379`
- RETHINK_URL, typically `localhost:28015`
- COMPOSE_SERVICE_PASSWORD, whatever you want

##Services You Need To Launch

- Redis
- RethinkDB

**NOTE**: this work sample was originally created with Go1.6, RethinkDB 2.0.x and [gorethink 1.0.0](https://github.com/dancannon/gorethink/releases/tag/v1.0.0).  Since releasing the work sample both the driver and the database have new minor version releases.  In testing, we found connection issues with new versions of the driver and the database.  You may use the versions that you wish, but we wanted to keep you from pulling your hair out. Please note: we have vendored gorethink 1.0.0 with this project for ease of use. So if you want to use a different version you will need to update the Godeps.

To seed the databases with some test data, run the following:

***WARNING***

We do flush the Redis database, so if you have other things you are working on, either start redis on a different port or backup your data


```shell
go test -tags=seed
```

To run the app,

Note: You should be running Go1.6 or Go1.5 with Vendor Experiment enabled,
as this ensures you're using the correct dependency versions from the
./vendor directory

```shell
go build ./... && ./platform-ws-services
```

If successful, you should see the following message:

```shell
listening on 8000
```

and if you've seeded the data, you should be able to run:

```shell
curl -u x:$COMPOSE_SERVICE_PASSWORD 'http://localhost:8000/margo/deployments/987654321'
```

which should return:

```shell
{
  "111111": [
    {
      "account": "compose-test",
      "capsule_id": "111111",
      "capsule_name": "redis0",
      "deployment_id": "987654321",
      "name":"redis_role",
      "output": "",
      "status":0
    }
  ]
}
```

