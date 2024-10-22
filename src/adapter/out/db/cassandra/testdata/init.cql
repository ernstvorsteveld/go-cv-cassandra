CREATE KEYSPACE testcv WITH replication = {'class':'SimpleStrategy', 'replication_factor':1};

CREATE TABLE testcv.experiences(id text primary key, name text, tags set<text>);
CREATE INDEX IF NOT EXISTS experiences_name ON testcv.experiences (name);

CREATE TABLE testcv.tags(id text primary key, name text);
