create table users (
    id serial primary key,
    username text not null unique,
    password text not null
);

create index idx_users_id on users(id);
create index idx_users_username on users(username);

create table posts (
    id serial primary key,
    title text not null,
    content text not null,
    author int not null references users(id),
    is_commentable bool not null,
    created_at timestamptz default current_timestamp
);

create index idx_posts_id on posts(id);
create index idx_posts_author on posts(author);
create index idx_posts_created_at on posts(created_at);

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

create index idx_comments_id on comments(id);
create index idx_comments_author on comments(author);
create index idx_comments_post on comments(post);
create index idx_comments_created_at on comments(created_at);