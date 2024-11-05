CREATE KEYSPACE testcv WITH replication = {'class':'SimpleStrategy', 'replication_factor':1};

CREATE TABLE testcv.cv_experiences(id text primary key, name text, tags set<text>);
CREATE INDEX IF NOT EXISTS cv_experiences_name ON testcv.cv_experiences (name);

CREATE TABLE testcv.cv_tags(id text primary key, name text);
