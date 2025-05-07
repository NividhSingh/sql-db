CREATE TABLE myTable (col1 VARCHAR (255) PRIMARY KEY, col2 INT);
INSERT INTO myTable (col1, col2) VALUES ('John', 42);
INSERT INTO myTable (col1, col2) VALUES ('Bob', 42);
INSERT INTO myTable (col1, col2) VALUES ('Bob', 22);
SELECT col1 AS c1, AVG(col2) FROM myTable GROUP BY col1;
SELECT blood_type, SUM(has_diabetes), sex AS male_or_female FROM MedicalRecords GROUP BY blood_type, sex;
