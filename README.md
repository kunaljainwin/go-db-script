# 🚀 PostgreSQL Load Testing Tool (Go)

A lightweight load testing utility written in **Golang** to benchmark **PostgreSQL** performance using concurrent database operations like `INSERT`.

---

## 📌 Description

This tool is designed to:
- Stress test a PostgreSQL database.
- Simulate high concurrency using goroutines.
- Measure performance of INSERT operations (can be extended).
- Generate randomized data for realistic load testing.

---

## 🔧 Tech Stack

- **Language**: Go
- **Database**: PostgreSQL
- **Driver**: `lib/pq` (`database/sql`)
- **Concepts Used**: Goroutines, Channels, WaitGroups, Random Data

