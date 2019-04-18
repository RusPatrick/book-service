CREATE USER docker_user WITH ENCRYPTED PASSWORD 'docker';
CREATE DATABASE docker OWNER docker_user;
\c docker;
CREATE SCHEMA IF NOT EXISTS books;

CREATE TABLE IF NOT EXISTS books.Books (
id serial primary key,
author varchar(100) not null,
title varchar(50) not null,
pages int not null,
publish_year int not null
);

CREATE TABLE IF NOT EXISTS books.users(
id serial primary key,
email varchar(100) not null unique,
password_hash text not null,
password_salt varchar(20) not null);


CREATE TABLE IF NOT EXISTS books.sessions(
id serial PRIMARY KEY,
user_id integer REFERENCES books.users(id) ON DELETE CASCADE,
session_id text not null,
exp int not null);

GRANT ALL PRIVILEGES ON SCHEMA books TO docker_user;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA books TO docker_user;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA books to docker_user;
