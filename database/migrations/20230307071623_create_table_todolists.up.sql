CREATE TABLE todolists
(
    id bigint NOT NULL AUTO_INCREMENT,
    title varchar (255) not null,
    status tinyint DEFAULT false,
    PRIMARY KEY (id)
);