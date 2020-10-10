package sqlite

import (
	"github.com/pkg/errors"
)

var (
	userTableSQL = `CREATE TABLE "user_info" (
  "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  "user_slug" TEXT NOT NULL,
  "real_name" TEXT NOT NULL,
  "user_avatar" TEXT NOT NULL,
  "is_cn" integer NOT NULL,
  "updated_time" integer NOT NULL,
  "created_time" integer NOT NULL
);

CREATE UNIQUE INDEX "user_info_slug"
ON "user_info" (
  "user_slug" COLLATE BINARY ASC,
  "is_cn" COLLATE BINARY ASC
);
`

	recordTableSQL = `CREATE TABLE "history_record" (
  "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  "user_slug" TEXT NOT NULL,
  "is_cn" integer NOT NULL,
  "ranking" integer NOT NULL,
  "solved_num" INTEGER NOT NULL,
  "zero_time" integer NOT NULL,
  "created_time" integer NOT NULL
);

CREATE UNIQUE INDEX "record_slug"
ON "history_record" (
  "user_slug" COLLATE BINARY ASC,
  "is_cn" COLLATE BINARY ASC,
  "zero_time" COLLATE BINARY ASC
);
`
)

func (l *lite) createTable() (err error) {

	_, err = l.db.Exec(userTableSQL)
	if err != nil {
		return errors.Wrapf(err, "create user_info table error")
	}

	_, err = l.db.Exec(recordTableSQL)
	if err != nil {
		return errors.Wrapf(err, "create history_record table error")
	}
	return nil
}
