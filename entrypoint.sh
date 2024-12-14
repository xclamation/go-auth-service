#!/bin/bash

DBSTRING="host=$DBHOST user=$DBUSER password=$DBPASSWORD dbname=$DBNAME sslmode=$DBSSL"

goose postgres "$DBSTRING" up

# Change to the /app directory
cd ..

# Start the application
exec "$@"