// Code generated by SQLBoiler 4.16.1 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package entity

import "testing"

func TestUpsert(t *testing.T) {
	t.Run("AccountRoles", testAccountRolesUpsert)

	t.Run("Accounts", testAccountsUpsert)

	t.Run("Roles", testRolesUpsert)

	t.Run("SchemaMigrations", testSchemaMigrationsUpsert)

	t.Run("TransferJobs", testTransferJobsUpsert)
}
