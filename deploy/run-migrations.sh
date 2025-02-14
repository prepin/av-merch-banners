#!/bin/sh

# Exit immediately if any command fails
set -e

# Print commands before executing them
set -x

# Run migrations
goose -dir=/app/migrations up

# Run seeds
goose -dir=/app/seeds -no-versioning up
