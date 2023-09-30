CREATE TABLE IF NOT EXISTS redirects (
id serial PRIMARY KEY,
urlpath VARCHAR ( 100 ) NOT NULL,
urltarget  VARCHAR ( 100 ) NOT NULL
);