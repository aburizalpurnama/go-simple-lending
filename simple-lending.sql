CREATE TYPE "statuses" AS ENUM (
  'active',
  'paidoff'
);

CREATE TABLE "accounts" (
  "id" INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
  "name" varchar NOT NULL,
  "limit" numeric NOT NULL
);

CREATE TABLE "loans" (
  "id" INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
  "amount" numeric NOT NULL,
  "paid_amount" numeric NOT NULL DEFAULT 0,
  "date" datetime NOT NULL DEFAULT (now() at time zone 'utc'),
  "status" statuses NOT NULL DEFAULT 'active',
  "acount_id" int NOT NULL
);

CREATE TABLE "installments" (
  "id" INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
  "amount" numeric NOT NULL,
  "paid_amount" numeric NOT NULL DEFAULT 0,
  "due_date" datetime NOT NULL,
  "status" statuses NOT NULL DEFAULT 'active',
  "loan_id" int NOT NULL
);

CREATE TABLE "payments" (
  "id" INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
  "date" datetime NOT NULL DEFAULT (now() at time zone 'utc'),
  "amount" numeric NOT NULL,
  "acount_id" int NOT NULL
);

ALTER TABLE "loans" ADD FOREIGN KEY ("acount_id") REFERENCES "accounts" ("id");

ALTER TABLE "installments" ADD FOREIGN KEY ("loan_id") REFERENCES "loans" ("id");

ALTER TABLE "payments" ADD FOREIGN KEY ("acount_id") REFERENCES "accounts" ("id");
