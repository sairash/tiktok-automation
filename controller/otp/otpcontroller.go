package opt

// func VerifyOtp(c echo.Context) error {
// 	type token struct {
// 		Otp         int    `json:"otp" validate:"required,numeric,gt=999,lt=10000"`
// 		AccessToken string `json:"access_token" validate:"required,len=5"`
// 	}

// 	request := token{}
// 	if err := c.Bind(&request); err != nil {
// 		return helper.ErrorResponse(c, "Error Validation", nil)
// 	}

// 	err := helper.Validator(request)
// 	if err != nil {
// 		return helper.ErrorResponse(c, fmt.Sprintf("%e", err), nil)
// 	}

// 	otp := models.Otp{
// 		Otp:         request.Otp,
// 		AccessToken: request.AccessToken,
// 	}

// 	helper.Database.Db.Preload("User.Country").Select("id", "user_id", "created_at").First(&otp)

// 	if otp.Id == 0 {
// 		return helper.ErrorResponse(c, "Otp Error", nil)
// 	}

// 	fmt.Println(otp.User)

// 	t1 := time.Now().Add(time.Duration(otp.User.Country.TimezoneOffset*60) * time.Minute)

// 	t2 := t1.Sub(otp.CreatedAt.Add(time.Duration(otp.User.Country.TimezoneOffset*60) * time.Minute))

// 	if t2 > 5*time.Minute {
// 		return helper.ErrorResponse(c, "This otp is late", nil)
// 	}

// 	claims := &helper.UserJwtClaims{
// 		Token: otp.User.Token,
// 		Admin: false,
// 		RegisteredClaims: jwt.RegisteredClaims{
// 			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
// 		},
// 	}

// 	token_gen := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

// 	t, err := token_gen.SignedString([]byte(helper.EnvVariable("SERECT")))
// 	if err != nil {
// 		return helper.ErrorResponse(c, "Token Not Correct", nil)
// 	}

// 	helper.Database.Db.Delete(&models.Otp{}, otp.Id)

// 	return helper.SuccessResponse(c, "Otp Verified", echo.Map{"token": t})
// }
