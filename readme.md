
# instruction

## 1. build binary
run `build_for_linux.sh`, this will compile `main.go` and output a binary file, and at the same time, cp `config.yaml` to vol directory for container

## 2. create network

run `create_network.sh`. This will create a new docker bridge for container to link to each other

## 3. build image
run `docker build -t hello:3`, this build the image to use

## 4. compose up

run `docker-compose up`

