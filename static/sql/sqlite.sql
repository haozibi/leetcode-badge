CREATE TABLE "user_info" (
  "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  "user_slug" TEXT NOT NULL,
  "real_name" TEXT NOT NULL,
  "user_avatar" TEXT NOT NULL,
  "is_cn" integer NOT NULL,
  "updated_time" integer NOT NULL,
  "created_time" integer NOT NULL,
  CONSTRAINT "user_slug" UNIQUE ("user_slug" COLLATE BINARY ASC)
);

CREATE INDEX "user_slug"
ON "user_info" (
  "user_slug" COLLATE BINARY ASC
);

CREATE TABLE "history_record" (
  "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  "user_slug" TEXT NOT NULL,
  "is_cn" integer NOT NULL,
  "ranking" integer NOT NULL,
  "solved_num" INTEGER NOT NULL,
  "zero_time" integer NOT NULL,
  "created_time" integer NOT NULL,
  CONSTRAINT "user_slug+zero_time" UNIQUE ("user_slug" COLLATE BINARY ASC, "zero_time" COLLATE BINARY ASC)
);

CREATE INDEX "record_slug"
ON "history_record" (
  "user_slug" COLLATE BINARY ASC
);

