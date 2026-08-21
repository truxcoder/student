CREATE TABLE students (
                          id TEXT PRIMARY KEY,
                          name TEXT NOT NULL,
                          gender TEXT NOT NULL,
                          birthday DATE NOT NULL,
                          score REAL,
                          classroom TEXT NOT NULL,
                          created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
                          updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
