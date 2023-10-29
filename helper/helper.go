package helper

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"image"
	"image/png"
	"os"

	"beepbop/models"
	"beepbop/seed"

	"github.com/gocolly/colly"
	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var BodyColors = []string{
	"B32821",
	"FFEA31",
	"D36E70",
	"BEBD7F",
	"D53032",
	"BDECB6",
	"1D334A",
	"CB3234",
	"F3A505",
	"F3DA0B",
	"45322E",
	"633A34",
	"00BB2D",
	"C51D34",
	"EFA94A",
	"E4A010",
	"7D8471",
	"31372B",
	"633A34",
	"587246",
	"D36E70",
	"B44C43",
	"D95030",
	"1B5583",
	"6D3F5B",
	"B5B8B1",
	"9DA1AA",
	"287233",
	"C7B446",
	"45322E",
	"8A9597",
	"E55137",
	"EC7C26",
	"2C5545",
	"F8F32B",
	"E6D690",
	"1D1E33",
	"924E7D",
	"FF7514",
	"7F7679",
	"9C9C9C",
	"75151E",
	"E6D690",
	"C1876B",
	"CB3234",
	"412227",
	"763C28",
	"6D3F5B",
	"7FB5B5",
	"063971",
	"FF7514",
	"4C9141",
	"231A24",
	"1D1E33",
	"CAC4B0",
	"89AC76",
	"646B63",
	"1D334A",
	"231A24",
	"343E40",
	"9E9764",
	"CBD0CC",
	"464531",
	"343B29",
}

type IpInfo struct {
	Ip       string `json:"ip"`
	City     string `json:"city"`
	Region   string `json:"region"`
	Country  string `json:"country"`
	Loc      string `json:"loc"`
	Timezone string `json:"timezone"`
	Org      string `json:"org"`
}

type DbInstance struct {
	Db *gorm.DB
}

type UserJwtClaims struct {
	Token string `json:"token"`
	Role  int    `json:"role"`
	jwt.RegisteredClaims
}

type RegisterDeviceResponse struct {
	Code int                `json:"code"`
	Msg  string             `json:"msg"`
	Data RegisterDeviceData `json:"data"`
}

type RegisterDeviceData struct {
	Cookie      string                 `json:"cookie"`
	DeviceInfo  map[string]interface{} `json:"device_info"`
	UserAgent   string                 `json:"user_agent"`
	DeviceId    string                 `json:"device_id"`
	InstallId   string                 `json:"install_id"`
	DeviceToken string                 `json:"device_token"`
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var (
	Validate = validator.New()
	Database DbInstance
)

func EnvVariable(key string) string {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()

	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}
	value, ok := viper.Get(key).(string)

	if !ok {
		log.Fatalf("Invalid type assertion")
	}

	return value
}

func RandomCharacterValueGen() (uint, uint, uint, uint, uint, uint, uint, models.CharacterSheet, models.CharacterSheet, models.CharacterSheet, models.CharacterSheet, models.CharacterSheet, models.CharacterSheet, models.CharacterSheet) {
	total := 0
	head_accessory := len(seed.HeadAccessory())
	head_accessory_len := rand.Intn(head_accessory)
	total += head_accessory

	body_hand_left := len(seed.BodyHandLeft())
	random_hand_left := rand.Intn(body_hand_left)
	body_hand_left_len := total + random_hand_left
	total += body_hand_left

	body_hand_right := len(seed.BodyHandRight())
	random_hand_right := rand.Intn(body_hand_right)
	body_hand_right_len := total + rand.Intn(random_hand_right)
	total += body_hand_right

	head_eye_left := len(seed.HeadEyeLeft())
	random_eye_left := rand.Intn(head_eye_left)
	head_eye_left_len := total + random_eye_left
	total += head_eye_left

	// head_eye_right := len(seed.HeadEyeRight())
	// random_eye_right := rand.Intn(head_eye_right)
	// head_eye_right_len := total + random_eye_right
	// total += head_eye_right

	random_eye_right := random_eye_left
	head_eye_right_len := total + random_eye_right
	total += head_eye_left

	head_mouth := len(seed.HeadMouth())
	random_head_mouth := rand.Intn(head_mouth)
	head_mouth_len := total + random_head_mouth
	total += head_mouth

	body_accessory := len(seed.BodyAccessory())
	random_body_accessory := rand.Intn(body_accessory)
	body_accessory_len := total + random_body_accessory
	total += body_accessory

	return uint(head_accessory_len), uint(body_hand_left_len), uint(body_hand_right_len), uint(head_eye_left_len), uint(head_eye_right_len), uint(head_mouth_len), uint(body_accessory_len), seed.HeadAccessory()[head_accessory_len], seed.BodyHandLeft()[random_hand_left], seed.BodyHandRight()[random_hand_right], seed.HeadEyeLeft()[random_eye_left], seed.HeadEyeRight()[random_eye_right], seed.HeadMouth()[random_eye_right], seed.BodyAccessory()[random_body_accessory]
}

func ProfileSvgCreator(body_color, head_eye_left, head_eye_right, head_mouth, head_accessory string) string {
	svg_builder :=
		`<svg xmlns="http://www.w3.org/2000/svg" width="119.16mm" height="167.17mm" version="1.1" viewBox="0 0 119.16 167.17" xml:space="preserve">
    <g transform="translate(-30.17 -65.855)">
      <path d="m77.352 153.74c-5.7903 4.7314-15.818 21.79-15.128 48.054 1.5702 5.56 5.1863 5.3835 8.6047 5.9678 0 0 28.585 0.11246 41.499-0.26571 3.1075-0.091 7.1029-2.0232 7.2148-5.9103s-2.1083-37.325-22.352-53.713" fill="#ffffff" />
      <path d="m104.39 160.04c-10.567-4.458-29.843-1.6335-27.118 2.8292 0 0 32.496-0.8782 33.415 9.5527 0.38842 4.4076-30.426-9.0267-40.09-0.63863-0.82616 0.71712 40.192 6.6291 42.777 10.024-9.0244 10.214-28.951-2.2038-43.666-4.1851 6.3789 12.08 29.821 14.192 45.297 13.543-11.283 1.8951-31.708 4.7999-48.904-1.4237 1.539 11.12 9.7923 16.626 49.787 12.415" fill="none" stroke="#` +
			body_color +
			`" stroke-linecap="round" stroke-linejoin="round" stroke-width="9" />
      <g stroke="#000000" stroke-linecap="round" stroke-linejoin="round">
        <path d="m77.249 154.17c-6.1376 4.4131-15.818 21.79-15.128 48.054 1.5702 5.56 5.1863 5.3835 8.6047 5.9678 0 0 28.585 0.11246 41.499-0.26571 3.1075-0.091 7.1029-2.0232 7.2148-5.9103s-0.3657-34.546-17.001-49.685" fill="none" stroke-width="1.9" />
      </g>
    </g>
    <g transform="translate(-30.17 -65.855)">
      <path d="m74.355 105.23c-9.7531 4.7299-17.713 48.569 13.077 49.225 32.879 0.70115 25.693-56.078-4.0678-55.797" fill="#ffffff" />
      <g fill="none" stroke="#000000" stroke-linecap="round">
        <path d="m74.095 104.58c-9.7531 4.7299-17.124 49.157 13.665 49.814 32.879 0.70115 25.104-64.321-12.115-60.9" stroke-width="2.2" />
        <g stroke-width="1.9">
        ` +
			head_eye_left +
			head_eye_right +
			head_mouth +
			`
        </g>
      </g>
    </g>
    ` +
			head_accessory +
			`
  </svg>
  `
	return svg_builder
}

func JWTAuthUser(token interface{}, user interface{}) error {
	switch token.(type) {
	case string:
		if err := Database.Db.First(&user, "token =?", token).Error; err != nil {
			return err
		}
	case uint:
		if err := Database.Db.First(&user, "id =?", token).Error; err != nil {
			return err
		}
	}
	return nil
}

func RandomNumber(n int) int {
	var num int
	for i := 0; i < n; i++ {
		num += rand.Intn(9)
		num = num * 10
	}
	num = num / 10

	return num
}

func RandomString(n int) string {
	result := make([]byte, n)

	for i := range result {
		result[i] = charset[rand.Intn(len(charset)-1)]
	}

	return string(result)
}

func Validator(request interface{}) error {
	err := Validate.Struct(request)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return err
		}

		for _, err := range err.(validator.ValidationErrors) {
			return fmt.Errorf(fmt.Sprintf("Required Field -> %s: %s; Got Value -> %v; Expected Value -> %s: %s", strings.ToLower(err.Field()), err.Type(), err.Value(), err.Tag(), err.Param()))
		}
	}
	return nil
}

func GetClientIP(r *http.Request) string {
	// First, try to get the IP from the X-Real-IP header
	ip := r.Header.Get("X-Real-IP")
	if ip != "" {
		return ip
	}

	// If X-Real-IP is not available, try the X-Forwarded-For header
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		ips := strings.Split(xff, ",")
		// The client's IP will be the first entry in the X-Forwarded-For list
		return strings.TrimSpace(ips[0])
	}

	// If both X-Real-IP and X-Forwarded-For are not available, fallback to RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	fmt.Println(ip)
	return ip
}

func JWT(c echo.Context) (UserJwtClaims, error) {
	cookie, err := c.Cookie("toke")
	if err != nil {
		return UserJwtClaims{}, fmt.Errorf("Missing token")
	}

	tokenString := cookie.Value

	fmt.Println(tokenString)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(EnvVariable("SERECT")), nil
	})

	if err != nil || !token.Valid {
		return UserJwtClaims{}, fmt.Errorf("Token is not valid")
	}

	claims := token.Claims.(jwt.MapClaims)

	fmt.Println(claims["token"].(string))

	return UserJwtClaims{
		Token: claims["token"].(string),
		Role:  int(claims["role"].(float64)),
	}, nil
}

func GetIPInfo(c echo.Context) (*IpInfo, error) {
	ip := GetClientIP(c.Request())
	url := fmt.Sprintf("http://ipinfo.io/%s/json", ip)
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var ipInfo IpInfo
	err = json.NewDecoder(response.Body).Decode(&ipInfo)
	if err != nil {
		return nil, err
	}

	return &ipInfo, nil
}

func SmsSender(to string, message string) error {
	apiUrl := "https://sairashgautam.com.np/"
	// apiUrl := "https://sms.aakashsms.com/sms/v3/send/"
	data := url.Values{}
	data.Set("auth_token", EnvVariable("SMS_AUTH_TOKEN"))
	data.Set("to", to)
	data.Set("text", message)
	payload := strings.NewReader(data.Encode())

	client := &http.Client{}
	req, err := http.NewRequest("POST", apiUrl, payload)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}
	bodyString := string(bodyBytes)
	fmt.Println(bodyString)
	return nil
}

func ErrorResponse(c echo.Context, message string, data echo.Map) error {
	return c.JSON(http.StatusBadRequest, echo.Map{"message": message, "data": data})
}

func ErrorResponse_er(c echo.Context, message error, data echo.Map) error {
	return c.JSON(http.StatusBadRequest, echo.Map{"message": message, "data": data})
}

func SuccessResponse(c echo.Context, message string, data echo.Map) error {
	return c.JSON(http.StatusOK, echo.Map{"message": message, "data": data})
}

func MessageCreator(val string, type_of_val string) string {
	switch type_of_val {
	case "otp":
		return fmt.Sprintf("%s is your BUUZZ security code. Do not share this with anyone.", val)
	default:
		return fmt.Sprintf("Buuzz Sent you message with value %s", val)
	}
}

func CreateAndSendOtp(user_id uint, user_phone string, db *gorm.DB) (string, error) {
	otp := RandomNumber(4)
	access_token := RandomString(5)
	models.CreateOtp(user_id, otp, access_token, db)
	// err := SmsSender(user_phone, MessageCreator(otp, "otp"))
	// return err
	return access_token, nil
}

func RandomCharacterGen() (uint, uint, uint, uint, uint, uint, uint, string) {
	head_accessory, body_hand_left, body_hand_right, head_eye_left, head_eye_right, head_mouth, body_accessory, head_accessory_command, body_hand_left_command, body_hand_right_command, head_eye_left_command, head_eye_right_command, head_mouth_command, body_accessory_command := RandomCharacterValueGen()
	body_color := BodyColors[rand.Intn(len(BodyColors))]
	profile_image_command := body_color + "-" + head_accessory_command.Command + "-" + body_hand_left_command.Command + "-" + body_hand_right_command.Command + "-" + head_eye_left_command.Command + "-" + head_eye_right_command.Command + "-" + head_mouth_command.Command + "-" + body_accessory_command.Command

	if _, err := os.Stat("./assets/images/profile/" + profile_image_command); errors.Is(err, os.ErrNotExist) {
		image_svg := ProfileSvgCreator(body_color, head_eye_left_command.SvgPath, head_eye_right_command.SvgPath, head_mouth_command.SvgPath, head_accessory_command.SvgPath)
		MakeSvgToPng("./assets/images/profile/", profile_image_command, image_svg)
	} else if err != nil {
		log.Fatal(err)
	}
	return head_accessory, body_hand_left, body_hand_right, head_eye_left, head_eye_right, head_mouth, body_accessory, profile_image_command
}

func MakeSvgToPng(path, image_name, svgString string) {
	// Parse SVG
	icon, err := oksvg.ReadIconStream(strings.NewReader(svgString))
	if err != nil {
		fmt.Println("Error parsing SVG:", err)
		return
	}

	// Create an image for rendering with transparent background
	width := int(icon.ViewBox.W)
	height := int(icon.ViewBox.H)
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	scanner := rasterx.NewScannerGV(width, height, img, img.Bounds())
	raster := rasterx.NewDasher(width, height, scanner)
	icon.Draw(raster, 1.0)

	// Create and save the PNG image with transparent background
	pngFile, err := os.Create(path + image_name + ".png")
	if err != nil {
		fmt.Println("Error creating PNG file:", err)
		return
	}
	defer pngFile.Close()

	if err := png.Encode(pngFile, img); err != nil {
		fmt.Println("Error encoding PNG:", err)
		return
	}
}

func FolderExists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func MakeDir(path string) error {
	return os.Mkdir(path, 0755)

}

func PageDataCreator(c echo.Context, title, header, body, sub_body, link string, show_app_bar bool, role int) echo.Map {
	page_data := echo.Map{
		"title":          title,
		"info":           header,
		"body":           body,
		"secondary_body": sub_body,
		"go_to_link":     link,
		"show_app_bar":   show_app_bar,
	}
	page_header := echo.Map{}

	if show_app_bar {
		if role == 0 {
			user_claims, err := JWT(c)

			if err != nil {
				role = 0
			} else {
				role = user_claims.Role
			}
		}

		if role == 0 {
			page_header["buttons"] = map[string]string{
				"Login":  "/signin",
				"Signup": "/signup",
			}
		} else if role == 1 {
			page_header["buttons"] = map[string]string{
				"Admin Dashboard": "/admin",
				"Logout":          "/logout",
			}
		} else if role == 2 {
			page_header["buttons"] = map[string]string{
				"Dashboard": "/home",
				"Logout":    "/logout",
			}
		}
	}

	page_data["header"] = page_header

	return page_data
}

func UserSidebar(url string) []echo.Map {
	var side = []echo.Map{
		{
			"text":    "Users",
			"is_text": true,
		},
		{
			"text":      "Disposable",
			"is_text":   false,
			"url":       "/home",
			"svg":       template.HTML(`<path stroke-linecap="round" stroke-linejoin="round" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />`),
			"is_active": false,
		},
		{
			"text":      "Contained",
			"is_text":   false,
			"url":       "/home/contained",
			"svg":       template.HTML(`<path stroke-linecap="round" stroke-linejoin="round" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z" />`),
			"is_active": false,
		},
		{
			"text":    "Posts",
			"is_text": true,
		},
		{
			"text":      "All Posts",
			"is_text":   false,
			"url":       "/home/posts",
			"svg":       template.HTML(`<path stroke-linecap="round" stroke-linejoin="round" d="M7.5 7.5h-.75A2.25 2.25 0 004.5 9.75v7.5a2.25 2.25 0 002.25 2.25h7.5a2.25 2.25 0 002.25-2.25v-7.5a2.25 2.25 0 00-2.25-2.25h-.75m0-3l-3-3m0 0l-3 3m3-3v11.25m6-2.25h.75a2.25 2.25 0 012.25 2.25v7.5a2.25 2.25 0 01-2.25 2.25h-7.5a2.25 2.25 0 01-2.25-2.25v-.75" />`),
			"is_active": false,
		},
		{
			"text":    "Automations",
			"is_text": true,
		},
		{
			"text":      "Automations",
			"is_text":   false,
			"url":       "/home/automations",
			"svg":       template.HTML(`<path stroke-linecap="round" stroke-linejoin="round" d="M16.023 9.348h4.992v-.001M2.985 19.644v-4.992m0 0h4.992m-4.993 0l3.181 3.183a8.25 8.25 0 0013.803-3.7M4.031 9.865a8.25 8.25 0 0113.803-3.7l3.181 3.182m0-4.991v4.99" />`),
			"is_active": false,
		},
	}

	for key, element := range side {
		if element["is_text"] != true {
			if url == element["url"] {
				side[key]["is_active"] = true
			}
		}
	}

	return side
}

func IsAccBanned(screen_name string) (map[string][]int, bool) {
	url := "https://www.tiktok.com/@" + screen_name
	c := colly.NewCollector()
	c.SetProxy("http://brd-customer-hl_1a3ab2f7-zone-ispproxy:unz9l48239vv@brd.superproxy.io:22225")

	mep := make(map[string][]int)
	c.OnHTML("strong", func(e *colly.HTMLElement) {
		value, err := strconv.Atoi(e.Text)
		if err == nil {
			mep[e.Attr("data-e2e")] = append(mep[e.Attr("data-e2e")], value)
		}
	})

	c.Visit(url)

	fmt.Println(mep)

	if len(mep) < 1 {
		return mep, true
	} else {
		return mep, false
	}
}

func GetVideoData(screen_name string, post_id string) map[string][]int {
	url := "https://www.tiktok.com/@" + screen_name + "/video/" + post_id
	c := colly.NewCollector()
	c.SetProxy("http://brd-customer-hl_1a3ab2f7-zone-ispproxy:unz9l48239vv@brd.superproxy.io:22225")

	mep := make(map[string][]int)
	c.OnHTML("strong", func(e *colly.HTMLElement) {
		value, err := strconv.Atoi(e.Text)
		if err == nil {
			mep[e.Attr("data-e2e")] = append(mep[e.Attr("data-e2e")], value)
		}
	})

	c.Visit(url)

	return mep
}

func GetDevices() RegisterDeviceResponse {
	url := "https://tiktok-video-no-watermark2.p.rapidapi.com/service/registerDevice?aid=1233&version=290304&os=9&idc=useast5"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("X-RapidAPI-Key", EnvVariable("RAPID_API_TOKEN_KEY"))
	req.Header.Add("X-RapidAPI-Host", "tiktok-video-no-watermark2.p.rapidapi.com")

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		fmt.Println(err)
	}

	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		fmt.Println("Error Fetching data!")
	}

	body, _ := io.ReadAll(res.Body)

	var response RegisterDeviceResponse
	err = json.Unmarshal(body, &response)

	if err != nil {
		fmt.Println(err)
	}

	response.Data.DeviceInfo["use_store_region_cookie"] = "1"

	return response

}

func MinMaxPorts(key string) int {
	port, err := strconv.Atoi(EnvVariable(key))
	if err != nil {
		return 5000
	}

	return port
}
