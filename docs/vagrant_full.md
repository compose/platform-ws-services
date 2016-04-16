#Full Vagrant Development

The platform-ws-services project depends on RethinkDB and Redis.

We have provided a Vagrantfile to ease your development cycle and reduce
what dependencies you need to worry about.

Please note: This requires Vagrant 1.5+

##Deploying

In order to provision a deployment issue a:

```shell
vagrant up
```

This will bring up:

* Redis VM running latest version of Redis
* RethinkDB VM running version 2.0.5 of RethinkDB
* A VM which downloads dependencies, runs tests, builds, seeds, and runs the
  platform-ws-services project

##Using

If everything deployed correctly you should be able to:

```shell
source setup-env.sh
curl -u x:$COMPOSE_SERVICE_PASSWORD $COMPOSE_SERVICE_URL/deployments/987654321
```

and get

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


where `$COMPOSE_SERVICE_PASSWORD` and `$COMPOSE_SERVICE_URL` are defined in
`setup-env.sh`

##Developing

You can work within the VM if you want via `vagrant ssh compose`. But you
probably prefer to write code on your own machine with your favorite
editors. To do this, just edit the code then perform a `vagrant provision
compose` and the compose VM will test, rebuild and relaunch the new
platform-ws-services.

##Of Note

This project was originally built with RethinkDB 2.0.5 and gorethink 1.0.0. The
RethinkDB VM is running RethinkDB 2.0.5 and we package the project with
gorethink 1.0.0 vendored.

The Vagrantfile is configured to launch the VMs to `10.0.32.2` `10.0.32.3` and
`10.0.32.4`. If you are already on a `10.0.32.x` network, you will want to edit
the `private_network` configuration in the Vagrantfile, and update the
`setup-env.sh` script to reflect those changes
