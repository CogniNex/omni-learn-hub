CREATE TYPE "entity_type" AS ENUM ('teacher', 'student', 'parent');

CREATE TYPE "backend_type" AS ENUM('int','varchar','decimal','text','datetime','bool');

CREATE TABLE IF NOT EXISTS "users"(
                        "user_id" UUID PRIMARY KEY,
                        "phone_number" varchar,
                        "password_hash" varchar,
                        "password_salt" varchar,
                        "refresh_token" varchar,
                        "refresh_expires_in" timestamp
);

CREATE TABLE IF NOT EXISTS "otp_codes" (
                           "otp_id" SERIAL PRIMARY KEY,
                           "phone_number" varchar,
                           "code" VARCHAR(6) NOT NULL,
                           "generation_attempts" int NOT NULL,
                           "is_verified" BOOLEAN DEFAULT false,
                           "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                           "expires_at" TIMESTAMP
);

CREATE TABLE IF NOT EXISTS "otp_blacklist"(
                                          "otp_blacklist_id" SERIAL PRIMARY KEY,
                                          "phone_number" varchar,
                                          "next_unblock_date" TIMESTAMP
);

CREATE TABLE IF NOT EXISTS "user_profiles"(
                                "user_id" UUID PRIMARY KEY,
                                "first_name" varchar,
                                "entity_id" int,
                                "entity_type_id" int,
                                "last_name" varchar,
                                "date_of_birth" date,
                                "language_id" int,
                                "email" varchar,
                                "is_active" bool,
                                "created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
                                "updated_at" timestamp,
                                UNIQUE ("entity_id", "entity_type_id")

);

CREATE TABLE IF NOT EXISTS "languages"(
                            "language_id" SERIAL PRIMARY KEY,
                            "language_code" varchar,
                            "language_name" varchar
);

CREATE TABLE IF NOT EXISTS "roles"(
                         "role_id" SERIAL PRIMARY KEY,
                         "role_name" varchar(255)
);

CREATE TABLE IF NOT EXISTS "user_roles"(
                              "user_id" UUID,
                              "role_id" int
);

CREATE TABLE IF NOT EXISTS "entity_attributes"(
                                     "entity_attribute_id" SERIAL PRIMARY KEY,
                                     "entity_type_id" int,
                                     "entity_id" int,
                                     "attribute_id" int,
                                     "attr_num" int,
                                     "attr_varchar" varchar,
                                     "attr_bool" bool,
                                     "attr_datetime" timestamp,
                                     "attr_decimal" decimal
);

CREATE TABLE IF NOT EXISTS "entities"(
                            "entity_type_id" SERIAL PRIMARY KEY,
                            "entity_type" entity_type
);

CREATE TABLE IF NOT EXISTS "attributes"(
                              "attribute_id" SERIAL PRIMARY KEY,
                              "entity_type_id" int,
                              "attribute_name" varchar(255),
                              "backend_type" backend_type,
                              "is_visible" bool,
                              "is_required" bool
);


CREATE OR REPLACE FUNCTION increment_entity_id()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.entity_id IS NULL THEN
        -- Check if the sequence for the entity_type_id exists
        IF NOT EXISTS(SELECT 1 FROM pg_sequences WHERE sequencename = 'entity_id_seq_' || NEW.entity_type_id) THEN
            -- Create a new sequence for the entity_type_id if it doesn't exist
            EXECUTE 'CREATE SEQUENCE entity_id_seq_' || NEW.entity_type_id || ' START 1';
END IF;
        -- Set the entity_id using the corresponding sequence
EXECUTE 'SELECT nextval(''entity_id_seq_' || NEW.entity_type_id || ''')' INTO NEW.entity_id;
END IF;
RETURN NEW;
END;
$$ LANGUAGE plpgsql;


CREATE TRIGGER trigger_increment_entity_id
    BEFORE INSERT ON "user_profiles"
    FOR EACH ROW EXECUTE FUNCTION increment_entity_id();

ALTER TABLE IF EXISTS "user_profiles" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("user_id");

ALTER TABLE IF EXISTS "user_profiles" ADD FOREIGN KEY ("entity_type_id") REFERENCES "entities" ("entity_type_id");

ALTER TABLE IF EXISTS "user_profiles" ADD FOREIGN KEY ("language_id") REFERENCES "languages" ("language_id");

ALTER TABLE IF EXISTS "user_roles" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("user_id");

ALTER TABLE IF EXISTS "user_roles" ADD FOREIGN KEY ("role_id") REFERENCES "roles" ("role_id");

ALTER TABLE IF EXISTS "entity_attributes" ADD FOREIGN KEY ("attribute_id") REFERENCES "attributes" ("attribute_id");

ALTER TABLE IF EXISTS "entity_attributes" ADD FOREIGN KEY ("entity_id", "entity_type_id") REFERENCES "user_profiles" ("entity_id", "entity_type_id");

ALTER TABLE IF EXISTS "attributes" ADD FOREIGN KEY ("entity_type_id") REFERENCES "entities" ("entity_type_id");

