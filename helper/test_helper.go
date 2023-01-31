package helper

import (
	"strconv"
	"time"

	"github.com/bayudha2/go-test-0/models"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const TableUserCreationQuery = `
	CREATE TABLE IF NOT EXISTS "public"."users" (
		"id" varchar(36) UNIQUE NOT NULL,
		"fullname" varchar(30) NOT NULL,
		"username" varchar(50) UNIQUE NOT NULL,
		"password" varchar(255) NOT NULL,
		"email" varchar(255) NOT NULL,
		"created_at" timestamptz NOT NULL DEFAULT now(),
		PRIMARY KEY ("id")
);`

const TableSessionCreationQuery = `
	CREATE TABLE IF NOT EXISTS "public"."sessions" (
    "id" varchar(36) UNIQUE NOT NULL,
    "username" varchar(50) NOT NULL,
    "refresh_token" varchar NOT NULL,
    "expires_at" timestamptz NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT "sessions_username_fkey" FOREIGN KEY ("username") REFERENCES "public"."users"("username"),
    PRIMARY KEY ("id")
);`

const TableProductCreationQuery = `
	CREATE TABLE IF NOT EXISTS "public"."products" (
		"id" varchar(36) UNIQUE NOT NULL,
		"name" varchar(30) NOT NULL,
		"price" numeric(10,2) NOT NULL DEFAULT 0.0,
		"created_at" timestamptz NOT NULL DEFAULT NOW(),
		"updated_at" timestamptz NOT NULL DEFAULT NOW(),
		PRIMARY KEY ("id")
	)
`

const TableProductIndexingQuery = `
	CREATE INDEX IF NOT EXISTS name_idx ON "public"."products"("name")
`

const TablePostCreationQuery = `
	CREATE TABLE IF NOT EXISTS "public"."posts" (
		"id" varchar(36) UNIQUE NOT NULL,
		"user_id" varchar(36) NOT NULL,
		"description" text NOT NULL,
		"created_at" timestamptz NOT NULL DEFAULT NOW(),
		"updated_at" timestamptz NOT NULL DEFAULT NOW(),
		CONSTRAINT "posts_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "public"."users"("id"),
		PRIMARY KEY ("id")
	);
`

func AddUsers(count int) {
	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
		hashPassword, _ := bcrypt.GenerateFromPassword([]byte("inipassword"+strconv.Itoa(i)), bcrypt.DefaultCost)
		models.DB.Exec("INSERT INTO users(id, fullname, username, password, email, created_at) VALUES($1, $2, $3, $4, $5, $6)",
			"iniuserid"+strconv.Itoa(i),
			"inifullname"+strconv.Itoa(i),
			"iniusername"+strconv.Itoa(i),
			string(hashPassword),
			"iniemail@"+strconv.Itoa(i)+".com",
			time.Now())
	}
}

func AddPost(count int, userid string) {
	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
		models.DB.Exec(`INSERT INTO posts(id, user_id, description, created_at, updated_at) 
		VALUES($1, $2, $3, $4, $5) 
		RETURNING id, user_id, description, created_at, updated_at`,
			"inipostid"+strconv.Itoa(i),
			userid,
			"ini post user ini yang ke - "+strconv.Itoa(i),
			time.Now(),
			time.Now(),
		)
	}
}

func AddProducts(count int) string {
	if count < 1 {
		count = 1
	}

	var productID string

	for i := 0; i < count; i++ {
		productID = uuid.New().String()
		models.DB.Exec("INSERT INTO products(id, name, price, created_at, updated_at) VALUES($1, $2, $3, $4, $5)",
			productID,
			"iniproduk"+strconv.Itoa(i),
			i*10.0,
			time.Now(),
			time.Now(),
		)
	}

	return productID
}

func AddSession(refresh string, expires int) {
	models.DB.Exec("INSERT INTO sessions(id, username, refresh_token, expires_at, created_at) VALUES($1, $2, $3, $4, $5) RETURNING username, refresh_token, expires_at, created_at",
		uuid.New().String(),
		"iniusername0",
		refresh,
		expires,
		time.Now())
}
