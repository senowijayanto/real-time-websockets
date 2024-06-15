# Real-time WebSocket Server with Golang and PostgreSQL

This project demonstrates how to implement a real-time WebSocket server using Golang, integrated with a PostgreSQL database. The server listens for database changes and pushes updates to connected WebSocket clients in real-time.

## Prerequisites

- Go 1.16+
- PostgreSQL 10+
- Basic knowledge of Golang and WebSockets
- Familiarity with SQL and PostgreSQL

## Setup

### 1. Clone the repository

```bash
git clone https://github.com/senowijayanto/real-time-websockets.git
cd real-time-websockets
```

### 2. Set up PostgreSQL
Create a new database and table, and set up the trigger for notifications:

```bash
CREATE DATABASE real_time_db;
\c real_time_db;

CREATE TABLE test_data (
    id SERIAL PRIMARY KEY,
    data TEXT NOT NULL
);

-- Trigger function to notify on changes
CREATE OR REPLACE FUNCTION notify_trigger() RETURNS trigger AS $$
BEGIN
    PERFORM pg_notify('events', NEW.data);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to call the function on insert
CREATE TRIGGER data_change
AFTER INSERT ON test_data
FOR EACH ROW
EXECUTE FUNCTION notify_trigger();
```

### 3. Configure the Go application
Update the PostgreSQL connection string in `main.go`:
```bash
connPool, err = pgx.Connect(context.Background(), "postgres://username:password@localhost:5432/real_time_db")
```
Replace `username` and `password` with your PostgreSQL credentials.

### 4. Install dependencies
```bash
go get github.com/gorilla/websocket
go get github.com/jackc/pgx/v4
```

### 5. Run the application
```
go run main.go
```

### 6. Open the HTML file
Open `index.html` in a web browser. The client will connect to the WebSocket server and display messages received from the server.

## Test the Connection
1. Ensure the WebSocket server is running.
2. Insert data into the `test_data` table in the PostgreSQL database to trigger notifications:
```
INSERT INTO test_data (data) VALUES ('Hello, WebSocket!');
```
3. Check the web browser to see if the message appears.

## Logging
The application logs important events and errors to the console. Ensure your terminal or log system captures these logs for debugging.
