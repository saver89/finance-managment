CREATE TYPE "office_type" AS ENUM (
  'hq',
  'warehouse',
  'store'
);

CREATE TYPE "office_state" AS ENUM (
  'active',
  'blocked'
);

CREATE TYPE "user_state" AS ENUM (
  'active',
  'blocked'
);

CREATE TYPE "currency_state" AS ENUM (
  'active',
  'disabled'
);

CREATE TYPE "account_state" AS ENUM (
  'active',
  'disabled'
);

CREATE TYPE "office_currency_type" AS ENUM (
  'retail',
  'supply'
);

CREATE TYPE "transaction_type" AS ENUM (
  'income',
  'outcome',
  'transfer',
  'adjustment'
);

CREATE TABLE "office" (
  "id" bigserial PRIMARY KEY,
  "parent_id" bigint NOT NULL,
  "type" office_type NOT NULL,
  "name" varchar NOT NULL,
  "state" office_state NOT NULL DEFAULT 'active',
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "deleted_at" timestamptz
);

CREATE TABLE "user" (
  "id" bigserial PRIMARY KEY,
  "office_id" bigint NOT NULL,
  "username" varchar NOT NULL,
  "password_hash" varchar NOT NULL,
  "first_name" varchar,
  "last_name" varchar,
  "middle_name" varchar,
  "birthday" timestamp,
  "email" varchar,
  "phone" varchar,
  "created_by" bigint,
  "state" user_state NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "deleted_at" timestamptz
);

CREATE TABLE "currency" (
  "id" bigserial PRIMARY KEY,
  "name" varchar NOT NULL,
  "short_name" varchar NOT NULL,
  "state" currency_state NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "deleted_at" timestamptz
);

CREATE TABLE "account" (
  "id" bigserial PRIMARY KEY,
  "office_id" bigint NOT NULL,
  "name" varchar NOT NULL,
  "balance" numeric(32,6) NOT NULL DEFAULT 0,
  "currency_id" bigint NOT NULL,
  "created_by" bigint NOT NULL,
  "state" account_state NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "deleted_at" timestamptz
);

CREATE TABLE "office_currency" (
  "id" bigserial PRIMARY KEY,
  "office_id" bigint NOT NULL,
  "currency_id" bigint NOT NULL,
  "type" office_currency_type NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "deleted_at" timestamptz
);

CREATE TABLE "office_currency_rate" (
  "id" bigserial PRIMARY KEY,
  "office_id" bigint NOT NULL,
  "from_currency_id" bigint NOT NULL,
  "to_currency_id" bigint NOT NULL,
  "rate" numeric(32,6) NOT NULL DEFAULT 1,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "deleted_at" timestamptz
);

CREATE TABLE "transactions" (
  "id" bigserial PRIMARY KEY,
  "office_id" bigint NOT NULL,
  "type" transaction_type NOT NULL,
  "from_account_id" bigint NOT NULL,
  "to_account_id" bigint,
  "amount" numeric(32,6) NOT NULL,
  "currency_id" bigint NOT NULL,
  "created_by" bigint,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "deleted_at" timestamptz
);

ALTER TABLE "office" ADD FOREIGN KEY ("parent_id") REFERENCES "office" ("id");

ALTER TABLE "user" ADD FOREIGN KEY ("office_id") REFERENCES "office" ("id");

ALTER TABLE "user" ADD FOREIGN KEY ("created_by") REFERENCES "user" ("id");

ALTER TABLE "account" ADD FOREIGN KEY ("office_id") REFERENCES "office" ("id");

ALTER TABLE "account" ADD FOREIGN KEY ("currency_id") REFERENCES "currency" ("id");

ALTER TABLE "account" ADD FOREIGN KEY ("created_by") REFERENCES "user" ("id");

ALTER TABLE "office_currency" ADD FOREIGN KEY ("office_id") REFERENCES "office" ("id");

ALTER TABLE "office_currency" ADD FOREIGN KEY ("currency_id") REFERENCES "currency" ("id");

ALTER TABLE "office_currency_rate" ADD FOREIGN KEY ("office_id") REFERENCES "office" ("id");

ALTER TABLE "office_currency_rate" ADD FOREIGN KEY ("from_currency_id") REFERENCES "currency" ("id");

ALTER TABLE "office_currency_rate" ADD FOREIGN KEY ("to_currency_id") REFERENCES "currency" ("id");

ALTER TABLE "transactions" ADD FOREIGN KEY ("office_id") REFERENCES "office" ("id");

ALTER TABLE "transactions" ADD FOREIGN KEY ("from_account_id") REFERENCES "account" ("id");

ALTER TABLE "transactions" ADD FOREIGN KEY ("to_account_id") REFERENCES "account" ("id");

ALTER TABLE "transactions" ADD FOREIGN KEY ("created_by") REFERENCES "user" ("id");

ALTER TABLE "transactions" ADD FOREIGN KEY ("currency_id") REFERENCES "currency" ("id");
