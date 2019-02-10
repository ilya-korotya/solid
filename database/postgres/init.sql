
CREATE DATABASE solid;
-- Set correct access for user because him inherits access form root(postgres) user 
CREATE USER lowcoder;
GRANT ALL ON DATABASE solid TO lowcoder;
