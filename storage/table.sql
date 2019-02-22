create table user (
    user_id bigint not null primary key auto_increment,
    name varchar(64) not null,
    email varchar(128) not null,
    salt varchar(256),
    salted varchar(256),
    created bigint not null,
    updated bigint,
    last_login bigint
);

create table follow (
    user_id bigint not null,
    follower_id bigint not null,
    created timestamp not null default '1970-01-01 00:00:01'
    state int not null,
);