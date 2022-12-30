CREATE TABLE slack_bot (
        id             		bigserial     PRIMARY KEY,
		name       			varchar(80),
		token       		varchar(80),
		channel_id 			varchar(80),
		cron_every 			bigint,
		last_cron_run       timestamp default NOW()
);

CREATE INDEX ON "slack_bot" ("id");


CREATE TABLE slack_bot_websites (
	id          bigserial     PRIMARY KEY,
	slack_id	bigint,
	website_id  bigint
);

CREATE INDEX ON "slack_bot_websites" ("id");
CREATE INDEX ON "slack_bot_websites" ("slack_id");