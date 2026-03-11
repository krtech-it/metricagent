CREATE TABLE metrics (
    id VARCHAR(50) PRIMARY KEY,
    m_type VARCHAR(50) NOT NULL,
    delta bigint default null,
    value double precision default null
)