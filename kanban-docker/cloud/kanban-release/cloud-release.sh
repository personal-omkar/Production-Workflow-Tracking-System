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

# # Define the data directory path for database
# DATA_DIR="./data"

# # Check if the data directory exists
# if [ ! -d "$DATA_DIR" ]; then
#     echo "Data directory does not exist. Creating..."
#     mkdir -p "$DATA_DIR"
# else
#     echo "Data directory already exists. Using existing database..."
# fi

# ./build.sh

containers=$(docker container ls --format '{{.Names}}' | grep 'rubber.')

#stop all containers
docker container stop $containers
docker container wait $containers
docker container rm $containers

docker load -i kanban-images.tar

docker-compose --env-file env -f cloud-kanban-compose.yml -p irpl-rubber-kanban build --no-cache 
docker-compose --env-file env -f cloud-kanban-compose.yml -p irpl-rubber-kanban up &

if pgrep -f "NatsClient"; then
    echo "Bash-Service is already running..."
  else
    echo "Bash-Service not found. Running ..."
    /RUBBER/scripts/service  NatsClient  0.0.0.0:4222 CMDLINE.BASH.COMMAND &
fi

