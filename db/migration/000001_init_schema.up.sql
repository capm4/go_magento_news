CREATE TABLE websites (
        id              bigserial     PRIMARY KEY AUTOINCREMENT,
		url             varchar(80),
		selector        varchar(80),
		attribute       varchar(80),
		last_url        varchar(80)
);

CREATE INDEX ON "websites" ("url");

CREATE TABLE config (
	id		bigserial     PRIMARY KEY,
	path 	varchar(80), 
	value 	TEXT
);

CREATE INDEX ON "config" ("path");