# Database Clients and Tools

Docker image with database clients and monitoring/debugging tools for database administration and troubleshooting.

## What's Included

- Database clients: PostgreSQL, MySQL, Redis, MongoDB
- Monitoring: `htop`, `netcat-openbsd`, `tcpdump`, `strace`
- Network utilities: `curl`, `wget`, `dnsutils`, `iputils-ping`
- JSON processing: `jq`

## Usage

```bash
# Build the image
make build

# Run interactively
make run

# Connect to a database (example)
make psql DB_HOST=localhost DB_NAME=mydb
```

## Comparison

See `Dockerfile.traditional` for the equivalent Dockerfile without apt-bundle. The traditional approach requires:
- Long RUN command with all database clients and tools listed inline
- Difficult to categorize which tools are for which purpose
- Hard to add or remove specific database clients
- No easy way to share tool lists across projects

Using apt-bundle allows you to organize dependencies in `Aptfile` with comments, making it clear what each tool provides and allowing easy sharing of database administration setups.

## Examples

```bash
# PostgreSQL
docker run -it db-tools psql -h dbhost -U user -d database

# MySQL
docker run -it db-tools mysql -h dbhost -u user -p database

# Redis
docker run -it db-tools redis-cli -h redishost
```

