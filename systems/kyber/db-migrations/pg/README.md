# PG DB-Migrations Directory

This directory contains all PG migrations for this system. We use `golang-migrate` for managing migrations.

## Naming Convention

```
{version}_{description}.{up|down}.sql
```

Example:

- `000001_create_users.up.sql`
- `000001_create_users.down.sql`
