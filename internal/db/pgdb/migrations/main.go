package migrations

import "github.com/uptrace/bun/migrate"

var Migrations *migrate.Migrations

func init() {
	Migrations = migrate.NewMigrations()
	if err := Migrations.DiscoverCaller(); err != nil {
		panic(err)
	}
}
