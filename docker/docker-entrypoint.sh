#!/usr/bin/env bash
set -e

if [ "$(id -u)" = '0' ]; then
  chown -R app:app $DATA_DIR
  exec gosu app "$0" "$@"
fi

exec start "$@"
