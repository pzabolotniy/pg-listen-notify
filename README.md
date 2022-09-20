# pg-listen-notify

[![golangci-lint](https://github.com/pzabolotniy/pg-listen-notify/actions/workflows/golangci-lint.yml/badge.svg?branch=main)](https://github.com/pzabolotniy/pg-listen-notify/actions/workflows/golangci-lint.yml)

PostgreSQL LISTEN/NOTIFY experiment

# Build && run

```bash
docker compose up -d --build --remove-orphans
```

## Check

```bash
docker compose ps
```

## Known issues

* If api or listener didn't start, just restart

```bash
docker compose up -d
```

# Goals of the experiment

* Research and develop PostgreSQL [notify/listen](https://www.postgresql.org/docs/current/sql-notify.html) functionality
* Implement REST to receive event-payload and notify _channel_
* Implement listener, which listens _channel_, selects event payload in exclusive mode
* Use [pgx](https://github.com/jackc/pgx) for db connection

# Results
1. There can be only one listener per channel
   * postgresql notifies every listener of a channel
   * it can be ok in some cases, but in my case i send eventID as a payload
      * all listeners catch notification (with eventID)
      * all listeners make SQL-SELECT to fetch event's payload. and this is not good
2. ***pgxpool.Pool** should be used instead of ***pgx.Conn**
3. I didn't find convenient collector for the [opentelemetry](https://opentelemetry.io/) and pgx
   * There is collector for the [sqlx](https://github.com/jmoiron/sqlx), but i don't use it =)
   * i used [github.com/exaring/otelpgx](github.com/exaring/otelpgx), but it doesn't trace placeholder values
   * need to find way to use **sqlx** for notify/listen and do not use **pgx** directly
4. PostgreSQL doesn't send notification to the listener when someone call NOTIFY.
   * listener should make any query (```SELECT 1```) to a database and then it will receive all the notification
   * you try it in two psql consoles.
   * pgx has some tricks and developer doesn't need to send "garbage" sql to a database.
