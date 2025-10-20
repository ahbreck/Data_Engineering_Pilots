#!/bin/bash
set -e

# Wait until Postgres is healthy and ready
echo "Waiting for PostgreSQL to be ready..."
until pg_isready -h "${PGADMIN_SERVER_HOST}" -U "${PGADMIN_SERVER_USER}" -d "${PGADMIN_SERVER_DB}" >/dev/null 2>&1; do
  sleep 2
done
echo "PostgreSQL is ready!"

# Create servers.json for pgAdmin
cat > /pgadmin4/servers.json <<EOF
{
  "Servers": {
    "1": {
      "Name": "${PGADMIN_SERVER_NAME}",
      "Group": "Servers",
      "Host": "${PGADMIN_SERVER_HOST}",
      "Port": ${PGADMIN_SERVER_PORT},
      "MaintenanceDB": "${PGADMIN_SERVER_DB}",
      "Username": "${PGADMIN_SERVER_USER}",
      "Password": "${PGADMIN_SERVER_PASSWORD}",
      "SSLMode": "prefer"
    }
  }
}
EOF

echo "pgAdmin server configuration created."
