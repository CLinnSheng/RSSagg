-- migration in simpler terms is like git control version but for database
-- commit changes to your dtabase schema

-- Up set of instructions to apply changes to the db schema 
-- +goose Up

CREATE TABLE users (
  id UUID PRIMARY KEY,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  name TEXT NOT NULL
);

-- revert back the changes made by the 'Up'
-- +goose Down
Drop TABLE users;
