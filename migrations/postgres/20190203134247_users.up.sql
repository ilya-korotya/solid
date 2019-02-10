CREATE TABLE users (
    -- They say it is better to use a 'uuid_generate_v1mc()' instead of a 'gen_random_uuid()'
    -- But i don't find information about this
    id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    first_name VARCHAR(255) NOT NULL,
    second_name VARCHAR(255) NOT NULL,
    age SMALLINT NOT NULL
);
