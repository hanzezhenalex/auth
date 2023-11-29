#!/bin/bash

set -exo pipefail

DOCKER_MYSQL_NAME='mysql-for-auth'

MYSQL_DEFAULT_DATABASE='alex-auth'
MYSQL_ROOT_USER="root"
MYSQL_ROOT_PASSWORD="alex"
MYSQL_BUSINESS_USER="alex"
MYSQL_BUSINESS_PASSWORD="alex"

docker run -d -p 3306:3306 --name "${DOCKER_MYSQL_NAME}" \
    -e MYSQL_ROOT_PASSWORD="${MYSQL_ROOT_PASSWORD}" \
    -e MYSQL_DATABASE="${MYSQL_DEFAULT_DATABASE}" \
    -e MYSQL_USER="${MYSQL_BUSINESS_USER}" \
    -e MYSQL_PASSWORD="${MYSQL_BUSINESS_PASSWORD}" \
    mysql/mysql-server:latest

sleep 30 # waiting for mysql

DOCKER_MYSQL_CONTAINER_ID=$(docker ps | grep "${DOCKER_MYSQL_NAME}" | awk '{print $1}')
if [ -z "${DOCKER_MYSQL_CONTAINER_ID}" ]; then
  echo "mysql is not working"
  exit 1
fi

echo "mysql container id: ${DOCKER_MYSQL_CONTAINER_ID}"

# prepare db
echo "create database ${MYSQL_DEFAULT_DATABASE}"
docker exec "${DOCKER_MYSQL_CONTAINER_ID}" mysql -u"${MYSQL_ROOT_USER}" -p"${MYSQL_ROOT_PASSWORD}" \
 -e "CREATE DATABASE IF NOT EXISTS ${MYSQL_DEFAULT_DATABASE};"

echo "grant full privilege for ${MYSQL_USER} on wechat"
docker exec "${DOCKER_MYSQL_CONTAINER_ID}" mysql -u"${MYSQL_ROOT_USER}" -p"${MYSQL_ROOT_PASSWORD}" \
  -e "GRANT ALL PRIVILEGES ON "\`${MYSQL_DEFAULT_DATABASE}\`".* TO '${MYSQL_BUSINESS_USER}'@'%';"