CREATE TABLE source (
    id serial PRIMARY KEY,
    name varchar(20)
);

INSERT INTO
    source (name)
VALUES
    ("youtube");

/* creator is uniquely identified by */
CREATE TABLE creator (
    id varchar(40),
    username varchar(40),
    source integer REFERENCES source (id) PRIMARY KEY (id, source_id)
);

CREATE TABLE tag (
    id serial PRIMARY KEY,
    name varchar(20)
);

CREATE TABLE content (
    id varchar(40),
    creator varchar(40) NOT NULL,
    source integer,
    title varchar(50) NOT NULL,
    PRIMARY KEY (id, source_id),
    FOREIGN KEY (creator, source) REFERENCES creator
);

CREATE TABLE tagged_content (
    content_id varchar(40),
    source_id integer,
    tag_id varchar(40) REFERENCES tag (id),
);

/*we first search */