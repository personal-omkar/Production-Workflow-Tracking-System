module irpl.com/kanban-dao

go 1.23.0

toolchain go1.24.3

require (
	github.com/go-ldap/ldap/v3 v3.4.10
	github.com/gorilla/mux v1.8.1
	github.com/lib/pq v1.10.9
	gorm.io/driver/postgres v1.5.9
	gorm.io/gorm v1.25.12
	irpl.com/kanban-commons v0.0.0
)

require (
	github.com/Azure/go-ntlmssp v0.0.0-20221128193559-754e69321358 // indirect
	github.com/go-asn1-ber/asn1-ber v1.5.7 // indirect
	github.com/golang-jwt/jwt/v4 v4.5.1 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgx/v5 v5.5.5 // indirect
	github.com/jackc/puddle/v2 v2.2.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/oklog/ulid/v2 v2.1.1 // indirect
	golang.org/x/crypto v0.38.0 // indirect
	golang.org/x/sync v0.14.0 // indirect
	golang.org/x/text v0.25.0 // indirect
)

replace irpl.com/kanban-commons => ../kanban-commons
