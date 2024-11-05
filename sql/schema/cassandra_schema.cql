CREATE KEYSPACE if not exists cv WITH replication = {'class':'SimpleStrategy', 'replication_factor':1};

USE cv;
create table if not exists cv_experiences(id text primary key, name text, tags set<text>);
CREATE INDEX IF NOT EXISTS experiences_name ON cv_experiences (name);
create table if not exists cv_tags(id text primary key, name text);
