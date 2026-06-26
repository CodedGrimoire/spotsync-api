# Interview Notes

## Question 1: What is a Race Condition? Explain how you used GORM Transactions and Row Locks (FOR UPDATE) to solve the EV Spot Bottleneck in this assignment.

A race condition happens when two or more requests try to read and write shared data at the same time, and the final result depends on timing. In a parking reservation system, this can happen when only one EV charging spot is left, but two users try to reserve it at almost the same moment.

If the code only checks available capacity first and then creates the reservation later, both requests might read the same value. They both see one spot available, and both create reservations. That means the system overbooks the parking zone.

In SpotSync, I solved this in the reservation repository using a GORM database transaction. The important part is that the capacity check and reservation creation happen inside the same transaction. Inside that transaction, I lock the selected `parking_zones` row using `FOR UPDATE`:

```go
tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&zone, reqZoneID)
```

This means if one request is checking and reserving a spot for that zone, another request for the same zone has to wait. After the row is locked, I count active reservations for that zone. If active reservations are already equal to or greater than total capacity, the repository returns a zone-full error. If there is capacity, it creates the reservation in the same transaction.

So the transaction gives atomic behavior, and the row lock gives concurrency safety. Together, they prevent two users from reserving the final spot at the same time.

## Question 2: How does GORM handle database connections under the hood? Why is it important to configure connection pooling?

GORM is an ORM, but under the hood it uses Go's standard `database/sql` package for database connections. When we call `gorm.Open`, GORM creates a database handle. That handle is not just one permanent connection. It manages a pool of connections that can be reused by different requests.

In a web API like SpotSync, many HTTP requests can arrive at the same time. For example, users may log in, list parking zones, or create reservations. If every request opened a brand-new database connection and closed it immediately, the API would be slower and the database could run out of connections.

Connection pooling solves this by keeping a controlled number of database connections available for reuse. In this project, the database setup configures:

```go
SetMaxIdleConns(10)
SetMaxOpenConns(100)
SetConnMaxLifetime(time.Hour)
```

`SetMaxIdleConns` controls how many unused connections can stay ready. `SetMaxOpenConns` limits the total number of open connections, so the app does not overload NeonDB or PostgreSQL. `SetConnMaxLifetime` prevents connections from staying open forever, which helps avoid stale connections.

This is important for production because connection pooling affects performance, reliability, and database resource usage. A good pool lets the API handle traffic efficiently without opening too many connections.

## Question 3: How do Goroutines differ from OS threads, and how does the Go scheduler manage them efficiently?

Goroutines are lightweight units of work managed by the Go runtime. They are not the same as OS threads. An OS thread is created and managed by the operating system, and it is relatively expensive. A goroutine is much cheaper, so a Go program can run thousands or even millions of goroutines depending on the workload.

The Go runtime uses a scheduler to run many goroutines on a smaller number of OS threads. The scheduler decides which goroutine should run, pauses goroutines when they are waiting, and resumes others. For example, if one goroutine is waiting for a database query or network operation, another goroutine can use the CPU instead of wasting time.

In a Go web server like SpotSync, each incoming HTTP request can be handled concurrently. Echo and the Go HTTP server can process many requests at the same time using goroutines. This is useful for APIs because one user might be logging in while another user is viewing zones and another is creating a reservation.

The Go scheduler makes this efficient by multiplexing goroutines over OS threads. So we get concurrency without manually creating and managing threads. It keeps the code simpler while still allowing the server to handle many requests at once.
