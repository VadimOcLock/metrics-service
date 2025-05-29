CREATE TABLE metrics (
                         id VARCHAR(255) PRIMARY KEY,
                         type VARCHAR(10) NOT NULL CHECK (type IN ('gauge', 'counter')),
                         delta BIGINT,
                         value DOUBLE PRECISION
);