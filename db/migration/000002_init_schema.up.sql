CREATE TABLE users (
        id              bigserial     PRIMARY KEY,
		name            varchar(80),
		login           varchar(80),
		password        varchar(80),
		role 			varchar(80),
		is_active       boolean,
		created_at		date NOT NULL,	
		updated_at      date NOT NULL default current_date,
		UNIQUE(login)
);

CREATE INDEX ON "users" ("login");