create table users (
    id serial primary key,
    username text not null unique,
    password text not null
);
create table posts (
    id serial primary key,
    title text not null,
    content text not null,
    author int not null references users(id),
    is_commentable bool not null,
    created_at timestamptz default current_timestamp
);
create table comments (
    id serial primary key,
    author int not null references users(id),
    message text not null check ( length(message) <= 2000 ),
    post int not null references posts(id),
    parent integer default -1,
    depth integer not null,
    path text not null default '',
    replies integer not null default 0,
    created_at timestamptz default current_timestamp
);

