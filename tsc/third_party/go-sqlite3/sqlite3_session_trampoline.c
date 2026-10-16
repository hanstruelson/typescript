// Copyright (C) OpenAI. Derived from the MIT licensed go-sqlite3 driver.

#include <stdint.h>
#include "sqlite3-binding.h"

extern int goSQLite3SessionFilter(uintptr_t, const char *);
extern int goSQLite3SessionConflict(uintptr_t, int);

static int session_filter(void *context, const char *table) {
  return goSQLite3SessionFilter(*(uintptr_t *)context, table);
}

static int session_conflict(void *context, int kind, sqlite3_changeset_iter *iter) {
  (void)iter;
  return goSQLite3SessionConflict(*(uintptr_t *)context, kind);
}

int sqlite3GoSessionApply(sqlite3 *db, int size, void *changeset, uintptr_t handle) {
  void *context = handle ? &handle : 0;
  return sqlite3changeset_apply(
      db, size, changeset,
      handle ? session_filter : 0,
      handle ? session_conflict : 0,
      context);
}

int sqlite3GoSetDefensive(sqlite3 *db, int enabled) {
  int current = 0;
  return sqlite3_db_config(db, SQLITE_DBCONFIG_DEFENSIVE, enabled, &current);
}

int sqlite3GoSetDQS(sqlite3 *db, int enabled) {
  int current = 0;
  int code = sqlite3_db_config(db, SQLITE_DBCONFIG_DQS_DDL, enabled, &current);
  if (code != SQLITE_OK) return code;
  return sqlite3_db_config(db, SQLITE_DBCONFIG_DQS_DML, enabled, &current);
}

int sqlite3GoEnableLoadExtension(sqlite3 *db, int enabled) {
  int current = 0;
  return sqlite3_db_config(db, SQLITE_DBCONFIG_ENABLE_LOAD_EXTENSION, enabled, &current);
}
