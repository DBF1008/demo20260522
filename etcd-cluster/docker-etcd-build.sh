#! /bin/bash

export COMPOSE_HTTP_TIMEOUT=500
export DOCKER_CLIENT_TIMEOUT=500
export COMPOSE_PARALLEL_LIMIT=1024
echo y | docker network prune

docker network create etcd
docker-compose -f docker-compose-etcd.yml up -d

etcdctl --endpoints=127.0.0.1:12379,127.0.0.1:22379,127.0.0.1:32379 member list