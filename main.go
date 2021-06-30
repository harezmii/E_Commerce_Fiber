package main

import (
	"e_commerce_furniture_with_fiber/database"
	"e_commerce_furniture_with_fiber/entity"
	_ "github.com/GoAdminGroup/go-admin/adapter/gofiber"
	_ "github.com/GoAdminGroup/go-admin/modules/db/drivers/postgres"
	_ "github.com/GoAdminGroup/go-admin/modules/db/drivers/sqlite"
	_ "github.com/GoAdminGroup/themes/adminlte"
	"github.com/go-playground/locales/tr"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	tr_translations "github.com/go-playground/validator/v10/translations/tr"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/postgres"
	"github.com/gofiber/template/django"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func main() {
	db := database.Connection()
	println(db)
	// DJANGO TEMPLATE
	engine := django.New("./public/views", ".html")


	// FIBER APP
	app := fiber.New(fiber.Config{
		Views:   engine,
	})
//	r := adminFiber.Gofiber{}
//	eng := adminEngine.Default()
//
//	cfg := config.Config{
//		Env: config.EnvLocal,
//		Theme: "adminlte",
//		Databases: config.DatabaseList{
//			"default": {
//				Host:       "127.0.0.1",
//				Port:       "5432",
//				User:       "postgres",
//				Pwd: 		"suat",
//				Name:       "furniture",
//				MaxIdleCon: 50,
//				MaxOpenCon: 150,
//				Driver:     config.DriverPostgresql,
//			},
//		},
//		UrlPrefix: "admin", // The url prefix of the website.
//		IndexUrl:  "/",
//		// Store must be set and guaranteed to have write access, otherwise new administrator users cannot be added.
//		Store: config.Store{
//			Path:   "./uploads",
//			Prefix: "uploads",
//		},
//		Language: language.EN,
//	}
//
//	// Add configuration and plugins, use the Use method to mount to the web framework.
//	_ = eng.AddConfig(&cfg).AddGenerators(datamodel.Generators).AddDisplayFilterXssJsFilter().AddGenerator("user", datamodel.GetUserTable).Use(app)


	// STATIC FILE
	app.Static("/", "./public/static/")
	app.Use(cors.New(cors.Config{
		AllowHeaders: "Origin, Content-Type, Accept",
	}))
	// ROUTE
	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.Render("index", fiber.Map{"data": "Gel"})
	})


	app.Get("/contact", func(ctx *fiber.Ctx) error {
		return ctx.Render("contact", fiber.Map{})
	})

	app.Get("/faq",func(ctx *fiber.Ctx) error {
		var faqs []entity.Faq
		db.Find(&faqs)
		return ctx.Render("faq",fiber.Map{"data" : faqs})
	})


	// STORAGE
	storage := postgres.New(postgres.Config{
		Username: "postgres",
		Password: "suat",
		Host:            "127.0.0.1",
		Port:            5432,
		Database:        "furniture",
		Table:           "sessions_storage",
		Reset:           false,
		GCInterval:      1 * time.Second,
		SslMode:         "disable",
	})

	// SESSION
	store := session.New(session.Config{
		Storage: storage,
		Expiration: 5 * time.Minute,
	})


	app.Get("/login",func(ctx *fiber.Ctx) error {
		return ctx.Render("login",fiber.Map{})
	})
	app.Get("/logout",func(ctx *fiber.Ctx) error {
		session, sessionError := store.Get(ctx)
		if sessionError != nil {
			println("Session Error")
		}

		session.Destroy()

		return ctx.Redirect("/")
	})

	app.Get("/register",func(ctx *fiber.Ctx) error {
		return ctx.Render("register",fiber.Map{})
	})

	app.Post("/loginCheck",func(ctx *fiber.Ctx) error{
		var result entity.User
		email    := ctx.FormValue("email")
		password := ctx.FormValue("password")

		db.Where(&entity.User{UserEmail: email}).First(&result)

		if result.ID != 0 {
			compareError := bcrypt.CompareHashAndPassword([]byte(result.UserPassword),[]byte(password))

			if compareError != nil {
				println("Hata var gardaşlık")
			} else {
				session, sessionError := store.Get(ctx)
				if sessionError != nil {
					println("Session Error")
				}
				session.Set("authorized",true)

				if errSave := session.Save(); errSave != nil {
					panic(errSave)
				}
				return ctx.Redirect("/admin")
			}

		}
		return ctx.Redirect("/login")
	})

	app.Post("/registerCheck",func(ctx *fiber.Ctx) error{
		fullName    := ctx.FormValue("full-name")
		email       := ctx.FormValue("email")
		password    := ctx.FormValue("password")
		rePassword:= ctx.FormValue("re-password")

		println(email)
		validate := validator.New()
		uni := ut.New(tr.New())
		trans, _ := uni.GetTranslator("tr")
		_ = tr_translations.RegisterDefaultTranslations(validate, trans)

		validateError := validate.Struct(&entity.User{
			UserName: fullName,
			UserPassword: password,
			UserEmail: email,
		})
		//validationErrors := validateError.(validator.ValidationErrors)

		if validateError != nil {
			var stringList []string
			for _, err := range validateError.(validator.ValidationErrors) {
				stringList = append(stringList, err.Translate(trans))//Age Must be greater than 18
			}
			return  ctx.JSON(stringList)
		} else {
			if password !=  rePassword {
				return  ctx.JSON([]string{"Password ve Re-Password Alanı Eşleşmiyor."})
			} else {
				hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
				if err != nil {
					panic("Hata")
				}
				result := db.Create(&entity.User{
					UserName: fullName,
					UserPassword: string(hash),
					UserEmail: email,
				})
				print("Result" , result.RowsAffected)
				if result.RowsAffected == 1 {
					return ctx.Redirect("/login")
				} else {
					return  ctx.JSON([]string{"İşlem Tamamlanamadı. Girilen alanları kontrol ediniz."})

				}
			}
		}
		return ctx.Redirect("/register")
	})

	app.Get("/receive", func(ctx *fiber.Ctx) error {
		return ctx.JSON(&entity.User{UserName: "deneme",UserEmail: "suatcnby06@gmail.com"})
	})
	// ADMIN ROUTE
	admin := app.Group("/admin")
	admin.Get("/",func(ctx *fiber.Ctx) error {
		session, sessionError := store.Get(ctx)
		if sessionError != nil {
			println("Session Error")
		}
		if session.Get("authorized") == true {
			return ctx.Render("admin-index", fiber.Map{})
		} else {
			return ctx.Redirect("/login")
		}

	})
	admin.Get("/addCategory",func(ctx *fiber.Ctx) error {
		return ctx.Render("admin-category", fiber.Map{})
	})

	app.Post("/post",func(ctx *fiber.Ctx) error {
		ctx.Vary()
		return ctx.Send(ctx.Body())
	})
	// MATCH ANY REQUEST
	app.Use(func(ctx *fiber.Ctx) error {
		return ctx.Render("404", fiber.Map{})
	})


	// FIBER SERVE
	err := app.Listen(":3000")
	if err != nil {
		return
	}
}
func ValidateStruct(user entity.User) []*entity.ErrorResponse {
	var errors []*entity.ErrorResponse
	validate := validator.New()
	err := validate.Struct(user)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element entity.ErrorResponse
			element.FailedField = err.Param()
			element.Field = err.Field()
			element.Value = err.ActualTag()
			errors = append(errors, &element)
		}
	}
	return errors
}