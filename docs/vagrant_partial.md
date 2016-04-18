#Partial Vagrant Development

The platform-ws-services project depends on RethinkDB and Redis.

We have provided a Vagrantfile to ease your development cycle and reduce
what dependencies you need to worry about.

Please note: This requires Vagrant 1.5+

##Deploying Databases

In order to provision a deployment issue a:

```shell
vagrant up redis rethinkdb
```

This will bring up:

* Redis VM running latest version of Redis at the address `10.0.32.3:6379`
* RethinkDB VM running version 2.0.5 of RethinkDB at address `10.0.32.2:28015`

## Development

You need the following environment variables set:

- REDIS_URL
- RETHINK_URL
- COMPOSE_SERVICE_PASSWORD

You can do this via `source setup-env.sh`


To seed the databases with some test data, run the following:

***WARNING***

We do flush the Redis database, so if you change the REDIS_URL environment
variable to point to your local redis, either start redis on a different port
or backup your data


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

##Of Note

This project was originally built with RethinkDB 2.0.5 and gorethink 1.0.0. The
RethinkDB VM is running RethinkDB 2.0.5 and we package the project with
gorethink 1.0.0 vendored.

The Vagrantfile is configured to launch the VMs to `10.0.32.2` and
`10.0.32.3`. If you are already on a `10.0.32.x` network, you will want to
edit the `private_network` configuration in the Vagrantfile, and update
the `setup-env.sh` script to reflect those changes
