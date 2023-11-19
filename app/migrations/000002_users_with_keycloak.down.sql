ALTER TABLE users DROP COLUMN external_id;
ALTER TABLE users DROP COLUMN description;
ALTER TABLE users DROP CONSTRAINT username_unique;
ALTER TABLE users DROP CONSTRAINT email_unique;
