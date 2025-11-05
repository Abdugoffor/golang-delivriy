package seeder

func DBSeed() {
	ProductSeeder()
	UserSeeder()
	RoleSeeder()
	CompanySeeder()
	CompanyUserRoleSeeder()
	CategorySeeder()
}
