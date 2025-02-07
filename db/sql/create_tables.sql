create extension if not exists citext;

create table users (
    username citext primary key not null,
    password text not null
);

create table posts (
    id serial primary key not null,
    author citext not null references users(username),
    message text not null,
    is_commentable bool not null,
)