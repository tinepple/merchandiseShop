create table users(
    id int generated by default as identity primary key,
    username text unique not null,
    password text not null,
);

create table merchandise(
    id int generated by default as identity primary key,
    name text unique not null,
    price int not null
);

create table balances(
    user_id int references users(id) not null,
    balance int default 1000
);

create table transactions(
    id int generated by default as identity primary key,
    user_id_from int references users(id) not null,
    user_id_to int references users(id) not null,
    amount int not null,
    created_at timestamp without time zone default NOW()
--  добавила created at, про него в задании не было сказано
);

create table purchases(
    id int generated by default as identity primary key,
    user_id int references users(id) not null,
    merchandise_id int references merchandise(id) not null,
    created_at timestamp without time zone default NOW()
);

insert into merchandise
    (name, price)
values
    ('t-shirt', 80),
    ('cup', 20),
    ('book', 50),
    ('pen', 10),
    ('powerbank', 200),
    ('hoody', 300),
    ('umbrella', 200),
    ('socks',	10),
    ('wallet', 50),
    ('pink-hoody', 500);

-- psql -U tinepple -d trainee_task
