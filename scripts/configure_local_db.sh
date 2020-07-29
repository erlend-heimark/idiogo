#!/bin/bash

export DB_DRIVER=sqlserver
export DB_USERNAME="sa"
export DB_PASSWORD="sOdifn3ijnvsd8!"
export DB_SERVER="localhost"
export DB_PORT="1433"
export DB_DATABASE="dadjokes"

export DOCKER_CONTAINER_NAME="idiogo_mssql"

echo "Creating database $DB_DATABASE"
export CMD="docker exec $DOCKER_CONTAINER_NAME /opt/mssql-tools/bin/sqlcmd -S $DB_SERVER -U $DB_USERNAME
 -P $DB_PASSWORD -d master"

echo "Make sure we can connect to the docker image. Will give up after 10 times"
counter=0
while true
do
    echo "attempting to connect to db: $counter"
    $CMD > /dev/null 2>&1
    if [ $? == 0 ]
    then
        echo "success..."
        break
    fi

    ((counter++))
    if [ $counter == 10 ]
    then
        echo "failed to connect to docker image. Exiting"
        exit 1
    fi
    sleep 3
done

$CMD -Q "CREATE DATABASE [${DB_DATABASE}];"

export CMD="docker exec $DOCKER_CONTAINER_NAME /opt/mssql-tools/bin/sqlcmd -S $DB_SERVER -U $DB_USERNAME
 -P $DB_PASSWORD -d $DB_DATABASE"

echo "Copying config sql"
docker cp scripts/sqlsetup/create_schemas_and_tables.sql $DOCKER_CONTAINER_NAME:/create_schemas_and_tables.sql
echo "Running config sql"
$CMD -i create_schemas_and_tables.sql
