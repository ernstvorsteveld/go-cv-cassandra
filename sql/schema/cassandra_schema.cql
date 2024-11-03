CREATE KEYSPACE if not exists cv WITH replication = {'class':'SimpleStrategy', 'replication_factor':1};

USE cv;
create table if not exists cv_experiences(id text primary key, experience text);
create table if not exists cv_tags(id text primary key, tag text);
