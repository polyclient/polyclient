#!/usr/bin/env bash

echo "⚠️ WARNING: You are about to remove all Docker-related data from your system, including running containers, images, volumes, and networks. This operation is irreversible and will reset Docker to its initial state. Are you sure you want to proceed? (y/n)"
read -r answer

if [ "$answer" != "${answer#[Yy]}" ]; then
    docker stop $(docker ps -aq)
    docker system prune -af
    docker volume prune -af
    docker network prune -f
    echo "☢️ Docker nuked"
else
    echo "Operation cancelled."
fi
