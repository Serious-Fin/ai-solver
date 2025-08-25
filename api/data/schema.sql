CREATE TABLE sqlite_sequence(name,seq);
CREATE TABLE IF NOT EXISTS "goTemplates" (
	"id"	INTEGER NOT NULL UNIQUE,
	"mainFunction"	TEXT NOT NULL,
	"testTemplate"	TEXT NOT NULL,
	"testHelpers"	TEXT,
	"problemFk"	INTEGER NOT NULL,
	PRIMARY KEY("id" AUTOINCREMENT)
);
CREATE TABLE IF NOT EXISTS "problems" (
	"id"	INTEGER NOT NULL UNIQUE,
	"title"	TEXT NOT NULL,
	"description"	TEXT NOT NULL,
	"testCases"	TEXT NOT NULL,
	"difficulty"	NUMERIC,
	PRIMARY KEY("id" AUTOINCREMENT)
);
CREATE TABLE IF NOT EXISTS "userCompletedProblems" (
	"userId"	INTEGER NOT NULL,
	"problemId"	INTEGER NOT NULL
);
CREATE TABLE IF NOT EXISTS "sessions" (
	"id"	TEXT NOT NULL UNIQUE,
	"userId"	INTEGER NOT NULL,
	"expiresAt"	TEXT,
	PRIMARY KEY("id")
);
CREATE TABLE IF NOT EXISTS "users" (
	"id"	TEXT NOT NULL UNIQUE,
	"email"	TEXT,
	"name"	TEXT,
	"profilePic"	TEXT
);
