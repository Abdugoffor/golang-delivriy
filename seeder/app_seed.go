package seeder

import (
	"fmt"
	app_model "my-project/modul/app/model"

	"gorm.io/gorm"
)

// ---- Form type "enum" lar ----

const (
	AppFormTypeNumber   = "number"
	AppFormTypeText     = "text"
	AppFormTypeTextarea = "textarea"
	AppFormTypeEmail    = "email"
	AppFormTypeSelect   = "select"
	AppFormTypeCheckbox = "checkbox"
)

// ---- Ichki konfiguratsiya structlari ----

type optionDef struct {
	Name string
	Slug string
}

type formDef struct {
	Name    string
	Slug    string
	Type    string
	Options []optionDef // select / checkbox uchun
	Value   string      // AppValue jadvalidagi value
}

type pageDef struct {
	Name  string
	Slug  string
	Forms []formDef
}

type categoryDef struct {
	Name  string
	Slug  string
	Pages []pageDef
}

// ---- 4 ta kategoriya + har birida page + form konfiguratsiyasi ----

var categories = []categoryDef{
	// 1) Umumiy murojaatlar
	{
		Name: "Umumiy murojaatlar",
		Slug: "general-requests",
		Pages: []pageDef{
			{
				Name: "Ariza yuborish",
				Slug: "submit-request",
				Forms: []formDef{
					{
						Name:  "F.I.Sh",
						Slug:  "full_name",
						Type:  AppFormTypeText,
						Value: "Ali Valiyev",
					},
					{
						Name:  "Telefon raqam",
						Slug:  "phone",
						Type:  AppFormTypeText,
						Value: "+998901234567",
					},
					{
						Name:  "Email manzil",
						Slug:  "email",
						Type:  AppFormTypeEmail,
						Value: "ali@example.com",
					},
					{
						Name: "Murojaat turi",
						Slug: "request_type",
						Type: AppFormTypeSelect,
						Options: []optionDef{
							{Name: "Shikoyat", Slug: "complaint"},
							{Name: "Taklif", Slug: "suggestion"},
							{Name: "Savol", Slug: "question"},
						},
						Value: "complaint",
					},
					{
						Name:  "Qo'shimcha izoh",
						Slug:  "comment",
						Type:  AppFormTypeTextarea,
						Value: "Tizim juda sekin ishlayapti.",
					},
					{
						Name: "Shaxsiy ma'lumotlar qayta ishlanishiga rozilik",
						Slug: "privacy_accept",
						Type: AppFormTypeCheckbox,
						Options: []optionDef{
							{Name: "Roziman", Slug: "yes"},
						},
						Value: "yes",
					},
				},
			},
			{
				Name: "Shikoyat sahifasi",
				Slug: "complaint-page",
				Forms: []formDef{
					{
						Name:  "Shikoyat sarlavhasi",
						Slug:  "complaint_title",
						Type:  AppFormTypeText,
						Value: "Xizmat sifati past",
					},
					{
						Name:  "Shikoyat izohi",
						Slug:  "complaint_body",
						Type:  AppFormTypeTextarea,
						Value: "Operator juda qo‘pol munosabatda bo‘ldi.",
					},
					{
						Name: "Bog'lanish usuli",
						Slug: "contact_method",
						Type: AppFormTypeSelect,
						Options: []optionDef{
							{Name: "Telefon", Slug: "phone"},
							{Name: "Email", Slug: "email"},
							{Name: "SMS", Slug: "sms"},
						},
						Value: "phone",
					},
					{
						Name: "Javob olishga rozilik",
						Slug: "allow_response",
						Type: AppFormTypeCheckbox,
						Options: []optionDef{
							{Name: "Ha", Slug: "yes"},
							{Name: "Yo'q", Slug: "no"},
						},
						Value: "yes",
					},
				},
			},
			{
				Name: "Takliflar sahifasi",
				Slug: "suggestions-page",
				Forms: []formDef{
					{
						Name:  "Taklif matni",
						Slug:  "suggestion_text",
						Type:  AppFormTypeTextarea,
						Value: "Yangi mobil ilova ishlab chiqing.",
					},
					{
						Name:  "Baholash (1–10)",
						Slug:  "rating",
						Type:  AppFormTypeNumber,
						Value: "9",
					},
					{
						Name: "Anonim bo‘lsinmi?",
						Slug: "is_anonymous",
						Type: AppFormTypeCheckbox,
						Options: []optionDef{
							{Name: "Ha", Slug: "yes"},
						},
						Value: "yes",
					},
				},
			},
		},
	},

	// 2) Texnik qo'llab-quvvatlash
	{
		Name: "Texnik qo'llab-quvvatlash",
		Slug: "technical-support",
		Pages: []pageDef{
			{
				Name: "Login muammolari",
				Slug: "login-issues",
				Forms: []formDef{
					{
						Name:  "Foydalanuvchi nomi",
						Slug:  "username",
						Type:  AppFormTypeText,
						Value: "user123",
					},
					{
						Name:  "Muammo tavsifi",
						Slug:  "problem_desc",
						Type:  AppFormTypeTextarea,
						Value: "Parolni tiklash xabari kelmayapti.",
					},
					{
						Name:  "Aloqa email",
						Slug:  "support_email",
						Type:  AppFormTypeEmail,
						Value: "user123@mail.com",
					},
					{
						Name: "Skrinshot bor-yo‘qligi",
						Slug: "has_screenshot",
						Type: AppFormTypeCheckbox,
						Options: []optionDef{
							{Name: "Bor", Slug: "yes"},
							{Name: "Yo'q", Slug: "no"},
						},
						Value: "yes",
					},
				},
			},
			{
				Name: "Xatolik haqida xabar berish",
				Slug: "bug-report",
				Forms: []formDef{
					{
						Name:  "Xatolik sarlavhasi",
						Slug:  "bug_title",
						Type:  AppFormTypeText,
						Value: "500 Internal Server Error",
					},
					{
						Name:  "Xatolik tafsiloti",
						Slug:  "bug_detail",
						Type:  AppFormTypeTextarea,
						Value: "Profil sahifasiga kirganda 500 xatosi chiqyapti.",
					},
					{
						Name: "Muhimlik darajasi",
						Slug: "priority",
						Type: AppFormTypeSelect,
						Options: []optionDef{
							{Name: "Past", Slug: "low"},
							{Name: "O‘rta", Slug: "medium"},
							{Name: "Yuqori", Slug: "high"},
						},
						Value: "high",
					},
				},
			},
			{
				Name: "Boshqa texnik muammolar",
				Slug: "other-issues",
				Forms: []formDef{
					{
						Name:  "Sarlavha",
						Slug:  "title",
						Type:  AppFormTypeText,
						Value: "Kirish vaqtida sekinlik",
					},
					{
						Name:  "Batafsil izoh",
						Slug:  "detail",
						Type:  AppFormTypeTextarea,
						Value: "Sahifalar juda sekin yuklanmoqda.",
					},
					{
						Name: "Uskuna turi",
						Slug: "device_type",
						Type: AppFormTypeSelect,
						Options: []optionDef{
							{Name: "Kompyuter", Slug: "pc"},
							{Name: "Telefon", Slug: "mobile"},
							{Name: "Planshet", Slug: "tablet"},
						},
						Value: "pc",
					},
					{
						Name: "Test ishlab ko‘rishga rozilik",
						Slug: "allow_testing",
						Type: AppFormTypeCheckbox,
						Options: []optionDef{
							{Name: "Ha", Slug: "yes"},
						},
						Value: "yes",
					},
				},
			},
		},
	},

	// 3) Hisob va to'lovlar
	{
		Name: "Hisob va to'lovlar",
		Slug: "billing",
		Pages: []pageDef{
			{
				Name: "Hisob-faktura so‘rash",
				Slug: "invoice-request",
				Forms: []formDef{
					{
						Name:  "Tashkilot nomi",
						Slug:  "company_name",
						Type:  AppFormTypeText,
						Value: "Example LLC",
					},
					{
						Name:  "INN",
						Slug:  "inn",
						Type:  AppFormTypeText,
						Value: "123456789",
					},
					{
						Name:  "Email",
						Slug:  "billing_email",
						Type:  AppFormTypeEmail,
						Value: "billing@example.com",
					},
					{
						Name: "To‘lov turi",
						Slug: "payment_type",
						Type: AppFormTypeSelect,
						Options: []optionDef{
							{Name: "Bank o‘tkazmasi", Slug: "bank_transfer"},
							{Name: "Naqd", Slug: "cash"},
						},
						Value: "bank_transfer",
					},
				},
			},
			{
				Name: "To‘lov tarixi filtri",
				Slug: "payment-history-filter",
				Forms: []formDef{
					{
						Name:  "Boshlang‘ich sana",
						Slug:  "from_date",
						Type:  AppFormTypeText,
						Value: "2026-01-01",
					},
					{
						Name:  "Tugash sanasi",
						Slug:  "to_date",
						Type:  AppFormTypeText,
						Value: "2026-01-31",
					},
					{
						Name: "Faqat muvaffaqiyatli to‘lovlar",
						Slug: "only_success",
						Type: AppFormTypeCheckbox,
						Options: []optionDef{
							{Name: "Ha", Slug: "yes"},
						},
						Value: "yes",
					},
				},
			},
			{
				Name: "Karta biriktirish",
				Slug: "card-binding",
				Forms: []formDef{
					{
						Name:  "Karta egasi F.I.Sh",
						Slug:  "card_holder",
						Type:  AppFormTypeText,
						Value: "Ali Valiyev",
					},
					{
						Name:  "Karta raqami (mask)",
						Slug:  "card_number_mask",
						Type:  AppFormTypeText,
						Value: "8600 **** **** 1234",
					},
					{
						Name: "Avtomatik to‘lovga rozilik",
						Slug: "auto_payment",
						Type: AppFormTypeCheckbox,
						Options: []optionDef{
							{Name: "Ha", Slug: "yes"},
						},
						Value: "yes",
					},
				},
			},
		},
	},

	// 4) Profil sozlamalari
	{
		Name: "Profil sozlamalari",
		Slug: "profile-settings",
		Pages: []pageDef{
			{
				Name: "Asosiy ma'lumotlar",
				Slug: "basic-info",
				Forms: []formDef{
					{
						Name:  "Ism",
						Slug:  "first_name",
						Type:  AppFormTypeText,
						Value: "Ali",
					},
					{
						Name:  "Familiya",
						Slug:  "last_name",
						Type:  AppFormTypeText,
						Value: "Valiyev",
					},
					{
						Name:  "Yosh",
						Slug:  "age",
						Type:  AppFormTypeNumber,
						Value: "25",
					},
					{
						Name: "Jinsi",
						Slug: "gender",
						Type: AppFormTypeSelect,
						Options: []optionDef{
							{Name: "Erkak", Slug: "male"},
							{Name: "Ayol", Slug: "female"},
						},
						Value: "male",
					},
				},
			},
			{
				Name: "Xavfsizlik sozlamalari",
				Slug: "security",
				Forms: []formDef{
					{
						Name: "2FA yoqilganmi?",
						Slug: "two_fa_enabled",
						Type: AppFormTypeCheckbox,
						Options: []optionDef{
							{Name: "Ha", Slug: "yes"},
							{Name: "Yo'q", Slug: "no"},
						},
						Value: "yes",
					},
					{
						Name:  "Tiklash email",
						Slug:  "recovery_email",
						Type:  AppFormTypeEmail,
						Value: "recover@example.com",
					},
				},
			},
			{
				Name: "Bildirishnomalar",
				Slug: "notifications",
				Forms: []formDef{
					{
						Name: "Bildirishnoma kanallari",
						Slug: "notification_types",
						Type: AppFormTypeCheckbox,
						Options: []optionDef{
							{Name: "Email", Slug: "email"},
							{Name: "SMS", Slug: "sms"},
							{Name: "Push", Slug: "push"},
						},
						// bir nechta tanlov: "email,sms"
						Value: "email,sms",
					},
					{
						Name: "Marketing xabarlarga rozilik",
						Slug: "allow_marketing",
						Type: AppFormTypeCheckbox,
						Options: []optionDef{
							{Name: "Ha", Slug: "yes"},
							{Name: "Yo'q", Slug: "no"},
						},
						Value: "no",
					},
				},
			},
		},
	},
}

// ---- Asosiy seeder funksiyasi ----

// SeedAppData:
// - app_categories
// - app_pages
// - app_form
// - app_option
// - apps
// - app_values
// jadvaliga ma'lumot soladi.
func SeedAppData(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		for _, c := range categories {
			// 1) Category
			category := app_model.AppCategory{
				Name:     c.Name,
				Slug:     c.Slug,
				IsActive: true,
			}
			if err := tx.Create(&category).Error; err != nil {
				return fmt.Errorf("create category %s: %w", c.Slug, err)
			}

			for _, p := range c.Pages {
				// 2) Page
				page := app_model.AppPage{
					Name:          p.Name,
					Slug:          p.Slug,
					AppCategoryID: category.ID,
					IsActive:      true,
				}
				if err := tx.Create(&page).Error; err != nil {
					return fmt.Errorf("create page %s: %w", p.Slug, err)
				}

				for _, f := range p.Forms {
					// 4) Form
					form := app_model.AppForm{
						AppPageID: page.ID,
						Name:      f.Name,
						Slug:      f.Slug,
						Type:      f.Type,
						IsRequire: false,
						IsActive:  true,
					}
					if err := tx.Create(&form).Error; err != nil {
						return fmt.Errorf("create form %s: %w", f.Slug, err)
					}

					// 5) Options (faqat select / checkbox uchun)
					if f.Type == AppFormTypeSelect || f.Type == AppFormTypeCheckbox {
						for _, opt := range f.Options {
							option := app_model.AppOption{
								AppFormID: form.ID,
								Name:      opt.Name,
								Slug:      opt.Slug,
								IsActive:  true,
							}
							if err := tx.Create(&option).Error; err != nil {
								return fmt.Errorf("create option %s for form %s: %w", opt.Slug, f.Slug, err)
							}
						}
					}

					// 3) Har bir page uchun bitta App (ariza) yaratamiz
					app := app_model.App{
						UserID:        1,
						AppCategoryID: category.ID,
						IsActive:      true,
					}

					if err := tx.Create(&app).Error; err != nil {
						return fmt.Errorf("create app for page %s: %w", p.Slug, err)
					}

					// 6) AppValue: shu ariza (app) bo‘yicha form javobi
					appValue := app_model.AppValue{
						AppID:     app.ID,
						AppPageID: page.ID,
						AppFormID: form.ID,
						Value:     f.Value,
						IsActive:  true,
					}

					if err := tx.Create(&appValue).Error; err != nil {
						return fmt.Errorf("create app value for form %s: %w", f.Slug, err)
					}
				}
			}
		}

		return nil
	})
}

// Qulaylik uchun: default userID = 1 bilan chaqirish
func SeedAppDataDefault(db *gorm.DB) error {
	return SeedAppData(db)
}
