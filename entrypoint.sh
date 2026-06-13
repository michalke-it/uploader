#!/bin/sh
set -e

# Allow the caller to supply the UID/GID that the server should run as, so
# files written to the mounted /app/uploads volume are owned correctly on the
# host. Defaults to 1000:1000 if not provided.
PUID="${PUID:-1000}"
PGID="${PGID:-1000}"

# (Re)create the group with the requested GID.
if ! getent group appgroup >/dev/null 2>&1; then
    groupadd -g "$PGID" appgroup
else
    groupmod -o -g "$PGID" appgroup
fi

# (Re)create the user with the requested UID.
if ! getent passwd appuser >/dev/null 2>&1; then
    useradd -o -u "$PUID" -g "$PGID" -d /app -s /sbin/nologin appuser
else
    usermod -o -u "$PUID" -g "$PGID" appuser
fi

# Make sure the upload target and log file are writable by that user.
mkdir -p /app/uploads
chown -R "$PUID:$PGID" /app/uploads
chmod u+rwx /app/uploads
touch /app/uploads.log
chown "$PUID:$PGID" /app/uploads.log

echo "Running as UID=$PUID GID=$PGID"

# Drop privileges and exec the server.
exec gosu "$PUID:$PGID" "$@"
