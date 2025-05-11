# Statistical Database with Differential Privacy

This project is a lightweight, in-memory SQL-like database engine built from
scratch with a focus on **privacy-preserving queries**. It supports core SQL
functionality (`CREATE TABLE`, `INSERT INTO`, `SELECT`, `GROUP BY`, and
aggregate functions) and adds **differential privacy**, **k‑anonymity**, and
**l‑diversity** to protect sensitive data. Developed as part of the ENGR3599
Special Topics in Computing – Introduction to Databases course at Olin College
of Engineering (Spring 2025) by Ahan Trivedi, Nividh Singh, and Ertug Umsur.

## Features

- **Command-Line Interface**  
  Issue SQL-like queries via the `input.sql` file.

- **Lexer & Parser**  
  Tokenizes and parses a simplified SQL grammar into an Abstract Syntax Tree
  (AST).

- **In-Memory Data Storage**  
  Tables and rows are represented using Go structs and slices/maps for fast,
  dependency-free operation.

- **SQL Support**

The following functions are supported, to a limited degree. An example is in
input.sql.

- `CREATE TABLE` with column definitions
  - You can create tables with types varchar, int and float (example on line 1
    of input.sql)
- `INSERT INTO` for adding rows
  - You can either insert into a table and then set all values for that row
    (example on line 23 of input.sql)
  - You can insert into a table, specify which values you're adding in
    parentheses after the table name, and then only add those values (example on
    line 173 of input.sql)
- `SELECT ... AS... GROUP BY ...` queries

  - Currently, you can't do where or having, but you can do group by. You can
    also add aliases when doing the select query. With every numerical column,
    you can do `COUNT`, `SUM`, `AVG`, `MIN`, `MAX`. Every column you select
    either needs to be one of these five statistics or needs to be in the
    groupby (example on line 174 of input.sql)

- **Differential Privacy**

  - Adds Laplace noise to numeric query results based on a configurable privacy
    budget (ε).
  - Controls cumulative privacy loss via an exponential decay rate across
    repeated `SELECT` queries.

- **k‑Anonymity & l‑Diversity Enforcement**

  - Removes rows whose quasi-identifier combinations occur fewer than _k_ times.
  - Drops sensitive-attribute columns unless they have at least _l_ distinct
    values.

- **Human-Readable Output**  
  Nicely formatted ASCII tables, respecting column visibility and aliases.

## Architecture

1. **Finle Input**  
   Users enter SQL-style queries in the `input.sql` file.

2. **Lexer** (`lexer.go`)  
   Breaks raw input into tokens (keywords, identifiers, literals, operators).

3. **Parser** (`parser.go`)  
   Builds an AST for each command (`CREATE`, `INSERT`, `SELECT`).

4. **Executor**

   - **DDL Commands**: Dispatches `CREATE` to `createTableFromAST` and `INSERT`
     to `insertIntoFromAST`.
   - **Query Commands**: `selectFromAST` retrieves rows, applies
     filters/grouping, computes aggregates, adds noise, and enforces privacy.

5. **Privacy Module** (`privacyFunctions.go`)

   - Calculates per-query ε from a global budget and decay rate.
   - Adds Laplace noise using `addNoise`.
   - Enforces k‑anonymity and l‑diversity on the noisy result set.

6. **Output** (`printTable`)  
   Prints results as ASCII tables showing only visible columns with applied
   aliases.

## Differential Privacy

Differential privacy prevents leakage of individual data points through
statistical queries. We implement:

- **Laplace Mechanism**  
  Adds noise drawn from a Laplace distribution `Lap(Δf/ε)` to numeric query
  results, where Δf is the sensitivity and ε is the privacy parameter.

- **Privacy Budget & Decay**  
  A total ε budget is allocated across queries with an exponential decay factor
  to manage cumulative privacy loss.

- **k‑Anonymity & l‑Diversity**  
  Complement DP by ensuring each returned record is indistinguishable among at
  least _k_ records and each sensitive attribute has at least _l_ distinct
  values.

## Getting Started

### Prerequisites

- Go (1.18+ recommended)

### Running the Project

Make sure you're in the correct directory for the project. Then, run the
following:

Clean up the dependencies

```bash
go mod tidy
```

Build the project

```bash
go build .
```

Run the project

```bash
go run .
```

### Editing the Code

There are a few edits you can make to the code.

1. First, you can edit the `input.sql` file with your SQL commands. These are
   limited to the functions listed in the features section above.
2. You can edit the privacy features in the file `main.go`. At the top of the
   code, lines 10-15, you can control parameters regarding the epsilon budget.
   Additionally, starting on line 80 you can customize how k-anonymity and
   l-diversity are enforced. Currently, they're enforcing both on all rows and
   columns, but you can change the k and l values and change which rows or
   columns they're applied to.
3. Finally on line 90, you can uncomment the next line to print the database.
   This is not a feature, but you can use this for debugging.

## Challenges

The most interesting problem we encountered was structuring the data to be both
flexible and type‑safe at runtime. At first, we had to design columns whose
types weren’t known until execution, which meant our schema needed to
accommodate any data type on the fly. As the project grew, we then had to figure
out how to store user‑defined aliases alongside those dynamic types. Perhaps the
trickiest part was implementing aggregates like average — we needed hidden sum
and count columns to compute the mean, yet ensure those helper columns never
showed up in the final output.

Beyond column design, we wrestled with representing queries as an abstract
syntax tree and then wiring that tree into our execution engine. To keep things
manageable, we broke the code into clear modules: one layer for database I/O,
another for AST traversal, and a suite of helper functions for grouping and
aggregation. Designing that modular structure was challenging, but it paid off
by making the system far easier to reason about, extend, and maintain.
