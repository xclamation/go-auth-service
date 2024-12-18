#!/bin/bash

DBSTRING="host=$DBHOST user=$DBUSER password=$DBPASSWORD dbname=$DBNAME sslmode=$DBSSL"

goose postgres "$DBSTRING" up

# Start the application
exec ./main "$@"