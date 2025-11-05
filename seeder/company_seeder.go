package seeder

import (
	"my-project/config"
	company_model "my-project/modul/company/model"
)

func CompanySeeder() {
	companies := []company_model.Company{
		{Name: "OpenAI LLC", ParentID: 0, UserID: 1, IsActive: true},
		{Name: "TechGroup", ParentID: 0, UserID: 2, IsActive: true},
		{Name: "MyCompany", ParentID: 0, UserID: 3, IsActive: true},
		{Name: "MyCompany 4", ParentID: 0, UserID: 4, IsActive: true},
		{Name: "MyCompany 5", ParentID: 0, UserID: 5, IsActive: true},
		{Name: "MyCompany 6", ParentID: 0, UserID: 6, IsActive: true},
	}

	for _, c := range companies {
		config.DB.FirstOrCreate(&c, company_model.Company{Name: c.Name})
	}
}
