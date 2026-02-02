#!/bin/bash
COMMAND=$1
ARG=$2

if [[ "$COMMAND" != "up" && "$COMMAND" != "down" && "$COMMAND" != "force" ]]; then
  echo "Usage: ./migrate.sh [up|down|force] [version/arg]"
  exit 1
fi

migrate --path sql/migrations \
--database postgres://admin:123@localhost:5432/learn_russian_db?sslmode=disable \
  $COMMAND $ARG