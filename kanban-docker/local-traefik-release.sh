#!/bin/bash

# Function to display the usage of the script
function usage() {
   echo "Usage: $0 [SOURCE_DIR] [force] [--skip-version-inc]"
   echo "SOURCE_DIR         Optional. The source directory to use."
   echo "force              Optional. Pass 'force' to force the build."
   echo "--skip-version-inc Optional. Flag to skip version increment."
   exit 1
}

# Check if the first argument is -h or --help to display the usage and exit
if [[ "$1" == "-h" || "$1" == "--help" ]]; then
   usage
fi

sudo ./add_hosts_entries.sh

sudo ./init.sh

./build.sh

containers=$(docker container ls --format '{{.Names}}' | grep 'rubber-.')

echo $containers

#stop all containers
docker container stop $containers
docker container wait $containers
docker container rm $containers

docker image rm $(docker image ls | grep -E '^rubber\.' | awk '{print $3}')
docker image rm $(docker images --format "{{.ID}}"  -f "dangling=true")|tr -d '\n'

docker compose --env-file env-traefik -f local-traefik-compose.yml -p rubber-kanban build --no-cache 
docker compose --env-file env-traefik -f local-traefik-compose.yml -p rubber-kanban up &

if pgrep -f "NatsClient"; then
    echo "Bash-Service is already running..."
  else
    echo "Bash-Service not found. Running ..."
    /RUBBER/scripts/service "$(docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' irpl-rubber-kanban-nats):4322" CMDLINE.BASH.COMMAND &
fi

docker container prune -f
docker image rm $(docker images --format "{{.ID}}"  -f "dangling=true")|tr -d '\n'
