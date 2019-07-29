#!/bin/sh

full_path=$(realpath $0)
project_path=$(dirname "$(dirname $full_path)")
env_path="$project_path/config/.env"
migration_path="$project_path/migration"

if [ $1 = "create" ]; then
    $( migrate create -ext sql -dir "${project_path}/migration" -digits 8 $2 )
    echo "created!"
    exit 1;
fi

while IFS= read -r line
do
    eval $line
done < "$env_path"

is_error=0

[ -z "$DATABASE_HOST" ] && { echo "Missing environment variable DATABASE_HOST" ; is_error=1; }
[ -z "$DATABASE_PORT" ] && { echo "Missing environment variable DATABASE_PORT" ; is_error=1; }
[ -z "$DATABASE_USER" ] && { echo "Missing environment variable DATABASE_USER" ; is_error=1; }
[ -z "$DATABASE_PASS" ] && { echo "Missing environment variable DATABASE_PASS" ; is_error=1; }
[ -z "$DATABASE_NAME" ] && { echo "Missing environment variable DATABASE_NAME" ; is_error=1; }
[ $is_error = 1 ] && { exit 1; }

conn="mysql://$DATABASE_USER:$DATABASE_PASS@tcp($DATABASE_HOST:$DATABASE_PORT)/$DATABASE_NAME"
if [ -n "$DATABASE_OPTION" ]; then
    conn=$conn?$DATABASE_OPTION
fi

migrate -path "$migration_path" -database "$conn" -verbose $1 $2

