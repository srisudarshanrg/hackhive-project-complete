package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"

	"github.com/srisudarshanrg/hackhive-project-competition/confidential"
	"github.com/srisudarshanrg/hackhive-project-competition/pkg/config"
	"github.com/srisudarshanrg/hackhive-project-competition/pkg/driver"
	"github.com/srisudarshanrg/hackhive-project-competition/pkg/models"
	"github.com/srisudarshanrg/hackhive-project-competition/pkg/render"
)

var Repo HandlerAccess
var db *sql.DB
var otp string
var loginStatus bool

type HandlerAccess struct {
	App *config.AppConfig
}

func SetAppConfigHandler(a *config.AppConfig) *HandlerAccess {
	return &HandlerAccess{
		App: a,
	}
}

func NewHandlers(r *HandlerAccess) {
	Repo = *r
}

func DatabaseAccess(database *sql.DB) {
	db = database
}

// Handlers for templates and posted forms

// Home is the handler for the home page
func (a *HandlerAccess) Home(w http.ResponseWriter, r *http.Request) {
	if !loginStatus {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	} else {
		loginStatusData := map[string]interface{}{}
		loginStatusData["loginStatus"] = "Logout"
		render.RenderTemplate(w, r, "home.page.tmpl", &models.TemplateData{
			Data: loginStatusData,
		})
	}
}

// CarboPrint is the handler for the carboprint page
func (a *HandlerAccess) CarboPrint(w http.ResponseWriter, r *http.Request) {
	if !loginStatus {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	} else {
		loginStatusData := map[string]interface{}{}
		loginStatusData["loginStatus"] = "Logout"
		render.RenderTemplate(w, r, "carboprint.page.tmpl", &models.TemplateData{
			Data: loginStatusData,
		})
	}
}

// RecycleLocator is the handler for the recycle locator page
func (a *HandlerAccess) RecycleLocator(w http.ResponseWriter, r *http.Request) {
	if !loginStatus {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	} else {
		loginStatusData := map[string]interface{}{}
		loginStatusData["loginStatus"] = "Logout"
		render.RenderTemplate(w, r, "recycle-locator.page.tmpl", &models.TemplateData{
			Data: loginStatusData,
		})
	}
}

// Login is the handler for login page
func (a *HandlerAccess) Login(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "login.page.tmpl", &models.TemplateData{
		CustomErrors: nil,
	})
}

// PostLogin is the handler for the login form
func (a *HandlerAccess) PostLogin(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	detail_entered := r.Form.Get("username_entered")
	password_entered := r.Form.Get("password_entered")

	searchDetailQuery := `select * from login_details where username=$1 or email=$2 or phone=$3`
	result, err := db.Exec(searchDetailQuery, detail_entered, detail_entered, detail_entered)

	if err != nil {
		log.Println(err)
	}

	rowsAffected, _ := result.RowsAffected()

	errorMap := map[string]string{}
	var errorString string

	if rowsAffected == 0 {
		errorString = "Username or Password incorrect"
		errorMap["notFound"] = errorString
		render.RenderTemplate(w, r, "login.page.tmpl", &models.TemplateData{
			CustomErrors: errorMap,
		})
	} else {
		var hashed_password string

		confirmPasswordQuery := `select password from login_details where username=$1 or email=$2 or phone=$3`
		row, err := db.Query(confirmPasswordQuery, detail_entered, detail_entered, detail_entered)
		if err != nil {
			log.Println(err)
		}

		defer row.Close()
		for row.Next() {
			err := row.Scan(&hashed_password)
			if err != nil {
				panic(err)
			}
		}

		check := GetPasswordFromHash(password_entered, hashed_password)

		if check {
			SetLoginStatus(true)
			http.Redirect(w, r, "/", http.StatusSeeOther)
		} else {
			errorString = "Username or Password incorrect"
			errorMap["notFound"] = errorString
			render.RenderTemplate(w, r, "login.page.tmpl", &models.TemplateData{
				CustomErrors: errorMap,
			})
		}
	}
}

// SignUp is the handler for sign up page
func (a *HandlerAccess) SignUp(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "sign-up.page.tmpl", &models.TemplateData{
		CustomErrors: nil,
	})
}

// PostSignUp is the handler for the sign up form
func (a *HandlerAccess) PostSignUp(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println("Could not parse form")
	}

	username := r.Form.Get("username")
	password := r.Form.Get("password")
	email := r.Form.Get("email")
	phone := r.Form.Get("phone")

	searchUniqueQueryUsername := `select username from login_details where username = $1 or email = $2 or phone = $3`
	result, err := db.Exec(searchUniqueQueryUsername, username, email, phone)

	if err != nil {
		log.Println(err)
	}

	rowsAffected, _ := result.RowsAffected()

	if rowsAffected == 0 {
		hashed_password, err := HashPassword(password)

		if err != nil {
			log.Fatal(err)
		}

		addRowQuery := `insert into login_details (username, password, email, phone) values($1, $2, $3, $4)`
		_, err = db.Exec(addRowQuery, username, hashed_password, email, phone)
		if err != nil {
			log.Println(err)
		}

		driver.DisplayRows(db)

		http.Redirect(w, r, "/login", http.StatusSeeOther)
	} else {
		errorMap := map[string]string{}

		errorText := "This user already exists. Choose another one."
		errorMap["uniqueDetail"] = errorText

		render.RenderTemplate(w, r, "sign-up.page.tmpl", &models.TemplateData{
			CustomErrors: errorMap,
		})
	}
}

// ForgotPassword is the handler for the forgot password page
func (a *HandlerAccess) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "forgot-password.page.tmpl", &models.TemplateData{})
}

func (a *HandlerAccess) PostForgotPassword(w http.ResponseWriter, r *http.Request) {
	otp := rand.Intn(999999)
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	email := r.Form.Get("email")

	if email == "" {
		errorMap := map[string]string{}

		errorMsg := "Enter valid email address to get OTP"
		errorMap["emailRequired"] = errorMsg

		render.RenderTemplate(w, r, "forgot-password.page.tmpl", &models.TemplateData{
			CustomErrors: errorMap,
		})

		return
	}

	checkEmailExistQuery := `select * from login_details where email=$1`
	result, err := db.Exec(checkEmailExistQuery, email)
	if err != nil {
		log.Println(err)
	}

	affected, _ := result.RowsAffected()

	if affected == 0 {
		errorMap := map[string]string{}

		errorString := "This email address does not have an account here. Would you like to sign up?"
		errorMap["emailNotFound"] = errorString

		render.RenderTemplate(w, r, "forgot-password.page.tmpl", &models.TemplateData{
			CustomErrors: errorMap,
		})
	}

	// Send email
	from := "srisudarshanrg@gmail.com"
	password := confidential.EmailPassword()
	to := []string{email}

	message := []byte(fmt.Sprintf("OTP for %s is %d", email, otp))

	err = SendEmail(from, to, message, password)

	if err != nil {
		log.Println(err)
	}

	var otp_convert = strconv.Itoa(otp)
	SetOtp(otp_convert)

	http.Redirect(w, r, "/otp-confirm", http.StatusSeeOther)
}

// ConfirmOTP is the handler for confirm otp page
func (a *HandlerAccess) ConfirmOTP(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "confirm-otp.page.tmpl", &models.TemplateData{})
}

// PostConfirmOTP handles the response by the form in the confirm-otp.page.tmpl page
func (a *HandlerAccess) PostConfirmOTP(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	otpRecieved := r.Form.Get("confirm_otp")

	// get correct otp. otp variable is accessible throughout the file
	otp_correct := otp

	if otpRecieved == otp_correct {
		http.Redirect(w, r, "/reset-password", http.StatusSeeOther)
	} else {
		otpError := map[string]string{}

		otpErrorString := "The OTP you have entered is incorrect. Please try again"
		otpError["otpError"] = otpErrorString

		render.RenderTemplate(w, r, "confirm-otp.page.tmpl", &models.TemplateData{
			CustomErrors: otpError,
		})
	}
}

// ResetPassword is the handler for the reset password page
func (a *HandlerAccess) ResetPassword(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "reset-password.page.tmpl", &models.TemplateData{})
}

func (a *HandlerAccess) PostResetPassword(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	email := r.Form.Get("email")
	newPassword := r.Form.Get("reset_password")
	confirmPassword := r.Form.Get("confirm_password")

	if newPassword == confirmPassword {
		newHashPassword, err := HashPassword(newPassword)

		if err != nil {
			log.Println(err)
		}

		resetPasswordQuery := `update login_details set password=$1 where email=$2`
		_, err = db.Exec(resetPasswordQuery, newHashPassword, email)
		if err != nil {
			log.Println(err)
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		errorMapPassword := map[string]string{}

		errorPassword := "Error changing password"
		errorMapPassword["errorPassword"] = errorPassword

		render.RenderTemplate(w, r, "reset-password.page.tmpl", &models.TemplateData{
			CustomErrors: errorMapPassword,
		})
	}
}

func (a *HandlerAccess) ResourceStatus(w http.ResponseWriter, r *http.Request) {
	if !loginStatus {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	} else {
		loginStatusData := map[string]interface{}{}
		loginStatusData["loginStatus"] = "Logout"
		render.RenderTemplate(w, r, "resource-status.page.tmpl", &models.TemplateData{
			CustomErrors: nil,
			Data:         loginStatusData,
		})
	}
}

func (a *HandlerAccess) PostResourceStatus(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	country_entered := r.Form.Get("country")
	country_entered = strings.ToLower(country_entered)

	getCountryQuery := "select * from resource_status where lower(country)=$1"
	result, err := db.Exec(getCountryQuery, country_entered)
	if err != nil {
		log.Println(err)
	}

	affected, _ := result.RowsAffected()

	if affected == 0 {
		errorMap := map[string]string{}

		errorString := "No such country exists. Please enter an existing country."
		errorMap["noCountry"] = errorString

		render.RenderTemplate(w, r, "resource-status.page.tmpl", &models.TemplateData{
			CustomErrors: errorMap,
		})

		return
	}

	var country, oil, electricity, coal, natural_gas, biofuel string
	var created_at, updated_at interface{}
	var id int

	row, err := db.Query(getCountryQuery, country_entered)
	if err != nil {
		log.Println(err)
	}
	defer row.Close()

	for row.Next() {
		err = row.Scan(&id, &country, &oil, &electricity, &coal, &natural_gas, &biofuel, &created_at, &updated_at)
		if err != nil {
			log.Println(err)
		}
	}

	type CountryDetailOverall struct {
		Country     string
		Oil         string
		Electricity string
		Coal        string
		NaturalGas  string
		Biofuel     string
	}

	specificCountryOverall := CountryDetailOverall{
		Country:     country,
		Oil:         oil,
		Electricity: electricity,
		Coal:        coal,
		NaturalGas:  natural_gas,
		Biofuel:     biofuel,
	}

	data := map[string]interface{}{}
	data["countryRowOverall"] = specificCountryOverall

	render.RenderTemplate(w, r, "resource-status.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

func (a *HandlerAccess) CoalProduction(w http.ResponseWriter, r *http.Request) {
	if !loginStatus {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	} else {
		loginStatusData := map[string]interface{}{}
		loginStatusData["loginStatus"] = "Logout"
		render.RenderTemplate(w, r, "coal-production.page.tmpl", &models.TemplateData{
			CustomErrors: nil,
			Data:         loginStatusData,
		})
	}
}

func (a *HandlerAccess) PostCoalProduction(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	countryResponse := r.Form.Get("countryCoal")
	yearResponse := r.Form.Get("yearCoal")

	checkAvailableQuery := `select * from coal_production where country=$1 and year=$2`
	result, err := db.Exec(checkAvailableQuery, countryResponse, yearResponse)
	if err != nil {
		log.Println(err)
	}

	affected, _ := result.RowsAffected()

	if affected == 0 {
		errorMap := map[string]string{}

		errorString := "Data not available for this country or year"
		errorMap["coalCountryNotFound"] = errorString

		render.RenderTemplate(w, r, "coal-production.page.tmpl", &models.TemplateData{
			CustomErrors: errorMap,
		})

		return
	}

	row, _ := db.Query(checkAvailableQuery, countryResponse, yearResponse)

	var country, code, year, coal_production string
	var created_at, updated_at interface{}
	var id int

	for row.Next() {
		err = row.Scan(&id, &country, &code, &year, &coal_production, &created_at, &updated_at)
		if err != nil {
			log.Println(err)
		}
	}

	type CoalCountry struct {
		Country        string
		Code           string
		Year           string
		CoalProduction string
	}

	specificCountryCoal := CoalCountry{
		Country:        country,
		Code:           code,
		Year:           year,
		CoalProduction: coal_production,
	}

	data := map[string]interface{}{}
	data["countryCoalDetail"] = specificCountryCoal

	render.RenderTemplate(w, r, "coal-production.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

func (a *HandlerAccess) OilProduction(w http.ResponseWriter, r *http.Request) {
	if !loginStatus {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	} else {
		loginStatusData := map[string]interface{}{}
		loginStatusData["loginStatus"] = "Logout"
		render.RenderTemplate(w, r, "oil-production.page.tmpl", &models.TemplateData{
			CustomErrors: nil,
			Data:         loginStatusData,
		})
	}
}

func (a *HandlerAccess) PostOilProduction(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	countryResponse := r.Form.Get("countryCoal")
	yearResponse := r.Form.Get("yearCoal")

	checkAvailableQuery := `select * from coal_production where country=$1 and year=$2`
	result, err := db.Exec(checkAvailableQuery, countryResponse, yearResponse)
	if err != nil {
		log.Println(err)
	}

	affected, _ := result.RowsAffected()

	if affected == 0 {
		errorMap := map[string]string{}

		errorString := "Data not available for this country or year"
		errorMap["coalCountryNotFound"] = errorString

		render.RenderTemplate(w, r, "coal-production.page.tmpl", &models.TemplateData{
			CustomErrors: errorMap,
		})

		return
	}

	row, _ := db.Query(checkAvailableQuery, countryResponse, yearResponse)

	var country, code, year, coal_production string
	var created_at, updated_at interface{}
	var id int

	for row.Next() {
		err = row.Scan(&id, &country, &code, &year, &coal_production, &created_at, &updated_at)
		if err != nil {
			log.Println(err)
		}
	}

	type CoalCountry struct {
		Country        string
		Code           string
		Year           string
		CoalProduction string
	}

	specificCountryCoal := CoalCountry{
		Country:        country,
		Code:           code,
		Year:           year,
		CoalProduction: coal_production,
	}

	data := map[string]interface{}{}
	data["countryCoalDetail"] = specificCountryCoal

	render.RenderTemplate(w, r, "coal-production.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

func (a *HandlerAccess) GasProduction(w http.ResponseWriter, r *http.Request) {
	if !loginStatus {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	} else {
		loginStatusData := map[string]interface{}{}
		loginStatusData["loginStatus"] = "Logout"
		render.RenderTemplate(w, r, "gas-production.page.tmpl", &models.TemplateData{
			CustomErrors: nil,
			Data:         loginStatusData,
		})
	}
}

func (a *HandlerAccess) PostGasProduction(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	countryResponse := r.Form.Get("countryCoal")
	yearResponse := r.Form.Get("yearCoal")

	checkAvailableQuery := `select * from coal_production where country=$1 and year=$2`
	result, err := db.Exec(checkAvailableQuery, countryResponse, yearResponse)
	if err != nil {
		log.Println(err)
	}

	affected, _ := result.RowsAffected()

	if affected == 0 {
		errorMap := map[string]string{}

		errorString := "Data not available for this country or year"
		errorMap["coalCountryNotFound"] = errorString

		render.RenderTemplate(w, r, "coal-production.page.tmpl", &models.TemplateData{
			CustomErrors: errorMap,
		})

		return
	}

	row, _ := db.Query(checkAvailableQuery, countryResponse, yearResponse)

	var country, code, year, coal_production string
	var created_at, updated_at interface{}
	var id int

	for row.Next() {
		err = row.Scan(&id, &country, &code, &year, &coal_production, &created_at, &updated_at)
		if err != nil {
			log.Println(err)
		}
	}

	type CoalCountry struct {
		Country        string
		Code           string
		Year           string
		CoalProduction string
	}

	specificCountryCoal := CoalCountry{
		Country:        country,
		Code:           code,
		Year:           year,
		CoalProduction: coal_production,
	}

	data := map[string]interface{}{}
	data["countryCoalDetail"] = specificCountryCoal

	render.RenderTemplate(w, r, "coal-production.page.tmpl", &models.TemplateData{
		Data: data,
	})
}
