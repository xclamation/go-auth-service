#!/bin/bash

# Wait for the database to be ready
# until pg_isready -h $DBHOST -p 5432 -U $DBUSER; do
#   echo "Waiting for database to be ready..."
#   sleep 2
# done

DBSTRING="host=$DBHOST user=$DBUSER password=$DBPASSWORD dbname=$DBNAME sslmode=$DBSSL"

#goose postgres "$DBSTRING" up
goose -dir ./migrations postgres "$DBSTRING" up

# Start the application
exec ./main "$@"