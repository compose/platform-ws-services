# -*- mode: ruby -*-
# vi: set ft=ruby :

# Vagrantfile API/syntax version. Don't touch unless you know what you're doing!
VAGRANTFILE_API_VERSION = "2"

Vagrant.configure(VAGRANTFILE_API_VERSION) do |config|

  config.vm.define "rethinkdb" do |rethinkdb|
    rethinkdb.vm.box="geerlingguy/ubuntu1404"
    rethinkdb.vm.box_version = "1.1.0"

    rethinkdb.vm.network "private_network", ip:"10.0.32.2"

    rethinkdb.vm.provision :shell do |sh|
        sh.inline = <<-EOF

          if [ ! -f /var/lock/provision.lock ]; then
            source /etc/lsb-release && echo "deb http://download.rethinkdb.com/apt $DISTRIB_CODENAME main" | sudo tee /etc/apt/sources.list.d/rethinkdb.list

            wget -qO- https://download.rethinkdb.com/apt/pubkey.gpg | sudo apt-key add -
            apt-get update --assume-yes
            apt-get install rethinkdb=2.0.5~0trusty --assume-yes

            cat > /etc/rethinkdb/instances.d/instance1.conf <<EOL
            runuser=vagrant
            rungroup=vagrant
            bind=all
EOL
            touch /var/lock/provision.lock

          fi  

          service rethinkdb restart
        EOF
      end
    end
  config.vm.define "redis" do |redis|
    redis.vm.box="geerlingguy/ubuntu1404"
    redis.vm.box_version = "1.1.0"

    redis.vm.network "private_network", ip:"10.0.32.3"

    redis.vm.provision :shell do |sh|
      sh.inline = <<-EOF
        if [ ! -f /var/lock/provision.lock ]; then
          apt-get update --assume-yes
          apt-get install redis-server --assume-yes

          sed -i 's/bind\ 127\.0\.0\.1/bind 0\.0\.0\.0/' /etc/redis/redis.conf
        fi
        service redis-server restart
      EOF
    end

  end
  config.vm.define "compose" do |compose|
    compose.vm.box="geerlingguy/ubuntu1404"
    compose.vm.box_version = "1.1.0"

    compose.vm.network "private_network", ip:"10.0.32.4"
    compose.vm.provision :shell do |sh|
      sh.inline = <<-EOF
        export GOPATH=/home/vagrant/go
        if [ ! -f /var/lock/provision.lock ]; then
          add-apt-repository ppa:ubuntu-lxc/lxd-stable --yes
          apt-get update --assume-yes
          apt-get install golang=2:1.6-1ubuntu2~ubuntu14.04.1~ppa1 --assume-yes

          echo "Setting Go Workspace"
          mkdir -p /home/vagrant/go/src
          mkdir -p /home/vagrant/go/pkg
          mkdir -p /home/vagrant/go/bin

          go get -u github.com/tools/godep

          echo "Linking platform-ws-services"
          mkdir -p $GOPATH/src/github.com/compose
          ln -s /vagrant $GOPATH/src/github.com/compose/platform-ws-services

          cat > /etc/init/platform-ws-services.conf <<EOL
# platform-ws-services
description "platform-ws-services daemon"
author "Compose.io"

# Listen and start after the vagrant-mounted event
start on net-device-up IFACE=eth1 and runlevel [2345]
stop on runlevel [!2345]

script
  while [ ! -f /vagrant/setup-env.sh ]
  do
    sleep 1
  done
  . /home/vagrant/go/src/github.com/compose/platform-ws-services/setup-env.sh
  /home/vagrant/go/bin/platform-ws-services --bind=0.0.0.0
end script
EOL
          initctl reload-configuration
          touch /var/lock/provision.lock
        fi


        echo "Restoring project environment"
        cd $GOPATH/src/github.com/compose/platform-ws-services
        $GOPATH/bin/godep restore
        go get -t .


        echo "Running Tests"
        go test -v

        echo "Building platform-ws-services"

        source ./setup-env.sh
        go build

        echo "Cleaning up any old deployments"
        service platform-ws-services stop

        echo "Seeding Databases"
        go test -tags=seed

        echo "Launching Compose Service"
        service platform-ws-services restart
      EOF
    end
  end
end
