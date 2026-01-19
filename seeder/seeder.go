package seeder

import "my-project/config"

func DBSeed() {
	// ProductSeeder()
	UserSeeder()
	RoleSeeder()
	CompanySeeder()
	CompanyUserRoleSeeder()
	CategorySeeder()
	SeedAppData(config.DB)
}
