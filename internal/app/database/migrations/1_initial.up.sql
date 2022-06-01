CREATE TABLE IF NOT EXISTS "companies"(
    "id" SERIAL PRIMARY KEY,
    "name" TEXT NOT NULL,
    "code" TEXT NOT NULL,
    "country" TEXT NOT NULL,
    "website" TEXT NOT NULL,
    "phone" TEXT NOT NULL
);