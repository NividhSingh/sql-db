# Statistical Database with Differential Privacy

This project is a lightweight statistical database engine built from scratch with a focus on **privacy-preserving queries**. It supports core SQL-like functionality — including `CREATE TABLE`, `INSERT INTO`, and `SELECT` queries — and adds **differential privacy** to protect user data during query execution.  It was developed as part of the ENGR3599 Special Topics in Computing - Introduction to Databases course at Olin College of Engineering (Spring 2025) by Ahan Trivedi, Nividh Singh, Ertug Umsur.

---

## ✨ Features

- Command-line interface for issuing SQL-like queries
- Lexer and parser for simplified SQL grammar
- In-memory data storage using Go structs
- Support for:
  - `CREATE TABLE` with column definitions
  - `INSERT INTO` for adding rows
  - `SELECT ... WHERE ...` queries (with single condition)
- Differential Privacy applied to query results (e.g., Laplace noise added to counts)

---

## 🏗️ Architecture

1. **CLI Input**: Users enter SQL-style queries via the terminal.
2. **Lexer**: Breaks input into tokens (keywords, identifiers, operators, etc.).
3. **Parser**: Converts tokens into an Abstract Syntax Tree (AST).
4. **Executor**:
   - Handles `CREATE`, `INSERT`, and `SELECT` based on AST type.
   - For `SELECT`, retrieves matching rows and computes the result.
5. **Differential Privacy Module**:
   - Applies Laplace noise to SELECT query outputs (e.g., counts or aggregates).
   - Ensures statistical privacy guarantees even with repeated queries.

---

## 🔐 Differential Privacy

We apply differential privacy to protect sensitive information in query results. This is especially important in statistical databases, where simple aggregate queries can leak private information. Our MVP adds **Laplace noise** based on a configurable privacy parameter (ε) to:
- Count queries (e.g., `SELECT COUNT(*)`)
- Numeric aggregations (e.g., `AVG`, `SUM`) [optional/future]

---

## 🚀 Getting Started

### Prerequisites
- Go (1.18+ recommended)

### Running the Project
```bash
go run main.go
