CREATE TYPE "entity_type" AS ENUM ('teacher', 'student', 'parent');

CREATE TYPE "backend_type" AS ENUM('int','varchar','decimal','text','datetime','bool');

CREATE TABLE IF NOT EXISTS "users"(
                        "user_id" UUID PRIMARY KEY,
                        "phone_number" varchar,
                        "otp_code" varchar,
                        "is_otp_verified" bool,
                        "password_hash" varchar,
                        "password_salt" varchar,
                        "refresh_token" varchar
);

CREATE TABLE IF NOT EXISTS "user_profiles"(
                                "user_id" UUID PRIMARY KEY,
                                "name" varchar,
                                "entity_id" int,
                                "entity_type_id" int,
                                "surname" varchar,
                                "date_of_birth" date,
                                "language_id" int,
                                "email" varchar,
                                "is_active" bool,
                                "created_at" timestamp with time zone,
                                "updated_at" timestamp with time zone,
                                UNIQUE ("entity_id", "entity_type_id")

);

CREATE TABLE IF NOT EXISTS "languages"(
                            "language_id" int PRIMARY KEY,
                            "language_code" varchar,
                            "language_name" varchar
);

CREATE TABLE IF NOT EXISTS "roles"(
                         "role_id" int PRIMARY KEY,
                         "role_name" varchar(255)
);

CREATE TABLE IF NOT EXISTS "user_roles"(
                              "user_id" UUID,
                              "role_id" int
);

CREATE TABLE IF NOT EXISTS "entity_properties"(
                                     "entity_attribute_id" int PRIMARY KEY,
                                     "entity_type_id" int,
                                     "entity_id" int,
                                     "property_id" int,
                                     "prop_num" int,
                                     "prop_varchar" varchar,
                                     "prop_bool" bool,
                                     "prop_datetime" timestamp with time zone,
                                     "prop_decimal" decimal
);

CREATE TABLE IF NOT EXISTS "entities"(
                            "entity_type_id" int PRIMARY KEY,
                            "entity_type" entity_type
);

CREATE TABLE IF NOT EXISTS "properties"(
                              "property_id" int PRIMARY KEY,
                              "entity_type_id" int,
                              "property_name" varchar(255),
                              "backend_type" backend_type,
                              "is_visible" bool,
                              "is_required" bool
);

ALTER TABLE IF EXISTS "user_profiles" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("user_id");

ALTER TABLE IF EXISTS "user_profiles" ADD FOREIGN KEY ("entity_type_id") REFERENCES "entities" ("entity_type_id");

ALTER TABLE IF EXISTS "user_profiles" ADD FOREIGN KEY ("language_id") REFERENCES "languages" ("language_id");

ALTER TABLE IF EXISTS "user_roles" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("user_id");

ALTER TABLE IF EXISTS "user_roles" ADD FOREIGN KEY ("role_id") REFERENCES "roles" ("role_id");

ALTER TABLE IF EXISTS "entity_properties" ADD FOREIGN KEY ("property_id") REFERENCES "properties" ("property_id");

ALTER TABLE IF EXISTS "entity_properties" ADD FOREIGN KEY ("entity_id", "entity_type_id") REFERENCES "user_profiles" ("entity_id", "entity_type_id");

ALTER TABLE IF EXISTS "properties" ADD FOREIGN KEY ("entity_type_id") REFERENCES "entities" ("entity_type_id");
