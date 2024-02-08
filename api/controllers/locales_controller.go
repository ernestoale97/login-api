package controllers

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func GetLocales(c echo.Context) error {
	return c.JSON(
		http.StatusOK,
		&echo.Map{
			"fr_Fr": &echo.Map{
				"login_hello":          "Bonjour",
				"login_username_label": "Nom d'utilisateur",
				"login_password_label": "Mot de passe",
				"pushed_button":        "Vous avez appuyé sur le bouton @count fois",
				"connection_failure":   "Hubo un error al obtener las traducciones",
			},
			"en_US": &echo.Map{
				"login_hello":          "Hello",
				"login_username_label": "Username",
				"login_password_label": "Password",
				"pushed_button":        "You have pushed the button @count times",
				"connection_failure":   "There was an error while getting the locales",
			},
		},
	)
}

func GetConfig(c echo.Context) error {
	return c.JSON(
		http.StatusOK,
		[]string{
			"name",
			/*"lastname",
			"username",
			"email",
			"phone",
			"city",
			"street",
			"instagram",
			"postalcode",
			"country",
			"password",
			"password_confirmation",
			"question",*/
			"birth",
			/*"sex",
			"citizen",*/
			"agree",
		},
	)
}

func Register(c echo.Context) error {
	return c.JSON(
		http.StatusOK,
		&echo.Map{
			"code": "error",
			"data": &echo.Map{
				"message": "error al crear user",
			},
		},
	)
}

func AppLogin(c echo.Context) error {
	return c.JSON(
		http.StatusOK,
		&echo.Map{
			"code": "Ok",
			"data": &echo.Map{
				"access_token":      "3|vwafVpnBY1vBdxoMCtsPXHHOJnWbo7FXiwIyrQt9",
				"type_verification": "loged",
			},
		},
	)
}

func UserInfo(c echo.Context) error {
	return c.JSON(
		http.StatusOK,
		&echo.Map{
			"code": "Ok",
			"data": &echo.Map{
				"id":                     127299,
				"name":                   "Ernesto",
				"lastname":               "Quintero Suarez",
				"fullname":               "Ernesto Quintero Suarez",
				"pic":                    "https://appmobile.mmtech-solutions.com/storage/ernestoale97_1706911950.jpg",
				"deposit":                1160,
				"withdraw":               -1256,
				"profit":                 96,
				"profitPositive":         1,
				"balance":                849,
				"mainAccount":            1263542,
				"imgCurrencyMainAccount": "https://appmobile.mmtech-solutions.com/storage/cryptocurrency/65c283f9dce15_.jpg",
				"lastTransactions": []echo.Map{
					{
						"description": "Deposit from back office",
						"move":        "Credit",
						"amount":      1000,
						"date":        "29-07-2023 07:18:23",
						"upToDown":    false,
					},
					{
						"description": "Challenge Purchased: 64c4bd5d2699b",
						"move":        "Debit",
						"amount":      -49,
						"date":        "29-07-2023 07:18:53",
						"upToDown":    true,
					},
					{
						"description": "Challenge Purchased: 64c6b32f2fee1",
						"move":        "Debit",
						"amount":      -229,
						"date":        "30-07-2023 18:59:59",
						"upToDown":    true,
					},
					{
						"description": "Profit from MT account (2100249478)",
						"move":        "Reward",
						"amount":      140,
						"date":        "31-07-2023 11:03:58",
						"upToDown":    false,
					},
					{
						"description": "testing",
						"move":        "Debit",
						"amount":      -800,
						"date":        "01-08-2023 18:04:14",
						"upToDown":    true,
					},
				},
			},
		},
	)
}

func ChallengesList(c echo.Context) error {
	return c.JSON(
		http.StatusOK,
		[]echo.Map{
			{
				"idaccount":      59219,
				"nameProduct":    "$ 5,000 Challenge Account",
				"status":         "Disable",
				"stage":          "Challenge",
				"login":          2100269820,
				"platform":       "MT4",
				"amount":         5000,
				"startDate":      "2023-08-19 19:41:22",
				"period":         "Monthly",
				"objetiveGoal":   "8",
				"objetiveResult": -0.07,
				"profit":         -3.61,
				"currentBalance": 4996.39,
				"currentEquity":  4996.39,
			},
			{
				"idaccount":      59367,
				"nameProduct":    "$5,000 Express Pro Account",
				"status":         "Disable",
				"stage":          "Funded",
				"login":          2100249478,
				"platform":       "MT4",
				"amount":         5000,
				"startDate":      "2023-10-20 23:26:45",
				"period":         "Monthly",
				"objetiveGoal":   0,
				"objetiveResult": 0,
				"profit":         0,
				"currentBalance": 5000,
				"currentEquity":  5000,
			},
			{
				"idaccount":      72731,
				"nameProduct":    "August Competition",
				"status":         "Disable",
				"stage":          "Challenge",
				"login":          2121930010,
				"platform":       "MT4",
				"amount":         10000,
				"startDate":      nil,
				"period":         "Monthly",
				"objetiveGoal":   "8",
				"objetiveResult": 0,
				"profit":         0,
				"currentBalance": 10000,
				"currentEquity":  10000,
			},
			{
				"idaccount":      72732,
				"nameProduct":    "$ 5,000 Challenge Account",
				"status":         "Disable",
				"stage":          "Challenge",
				"login":          2121930011,
				"platform":       "MT4",
				"amount":         5000,
				"startDate":      nil,
				"period":         "Monthly",
				"objetiveGoal":   "8",
				"objetiveResult": 0,
				"profit":         0,
				"currentBalance": 5000,
				"currentEquity":  5000,
			},
			{
				"idaccount":      72734,
				"nameProduct":    "$ 5,000 Challenge Account",
				"status":         "Disable",
				"stage":          "Challenge",
				"login":          2122000001,
				"platform":       "MT4",
				"amount":         5000,
				"startDate":      nil,
				"period":         "Monthly",
				"objetiveGoal":   "8",
				"objetiveResult": -80,
				"profit":         -4000,
				"currentBalance": 1000,
				"currentEquity":  1000,
			},
			{
				"idaccount":      88058,
				"nameProduct":    "$ 400,000 Challenge Account",
				"status":         "Enable",
				"stage":          "Challenge",
				"login":          2122011334,
				"platform":       "MT4",
				"amount":         400000,
				"startDate":      "2023-11-02 20:34:25",
				"period":         "Monthly",
				"objetiveGoal":   "8",
				"objetiveResult": -0.3,
				"profit":         -1207.59,
				"currentBalance": 398792.41,
				"currentEquity":  398792.41,
			},
			{
				"idaccount":      100915,
				"nameProduct":    "testing",
				"status":         "Enable",
				"stage":          "Challenge",
				"login":          2122020163,
				"platform":       "MT4",
				"amount":         10,
				"startDate":      nil,
				"period":         "Monthly",
				"objetiveGoal":   "10",
				"objetiveResult": 0,
				"profit":         0,
				"currentBalance": 10,
				"currentEquity":  10,
			},
			{
				"idaccount":      100916,
				"nameProduct":    "testing",
				"status":         "Enable",
				"stage":          "Challenge",
				"login":          2122020164,
				"platform":       "MT4",
				"amount":         10,
				"startDate":      nil,
				"period":         "Monthly",
				"objetiveGoal":   "10",
				"objetiveResult": 0,
				"profit":         0,
				"currentBalance": 10,
				"currentEquity":  10,
			},
		},
	)
}

func ChallengeMetrics(c echo.Context) error {
	return c.JSON(
		http.StatusOK,
		&echo.Map{
			"code": "Ok",
			"data": &echo.Map{
				"access_token":      "3|vwafVpnBY1vBdxoMCtsPXHHOJnWbo7FXiwIyrQt9",
				"type_verification": "loged",
			},
		},
	)
}

/*{
  "code": "Ok",
  "data": {
    "message": "An email has been sent to the address provided with an activation link, please check your email."
  }
}*/
