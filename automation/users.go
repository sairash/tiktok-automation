package automation

import (
	"beepbop/helper"
	"beepbop/models"
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var (
	currentDir, _ = os.Getwd()
	relativePath  = "./assets/posts/"

	fullPath = filepath.Join(currentDir, relativePath)

	devices_in_use SafeDevice

	min_ports = helper.MinMaxPorts("MIN_PORTS")
	max_ports = helper.MinMaxPorts("MAX_PORTS")
)

type SafeDevice struct {
	device []uint
	mu     sync.Mutex
}

type TiktokUploadResponse struct {
	PostId string `json:"post_id"`
}

func BlockDevice(device_id uint) {
	helper.Database.Db.Model(&models.Device{}).Where("id = ?", device_id).Update("blocked", 1)
}

func createNewDevice() models.Device {
	device := models.Device{}
	new_device := helper.GetDevices()

	helper.Database.Db.Where("d_id = ? AND i_id = ?", new_device.Data.DeviceId, new_device.Data.InstallId).Find(&device)

	fmt.Println(device)

	if device.Id != 0 {
		fmt.Println("Already Exists!")
		// return createNewDevice()
	}

	device.UserAgent = new_device.Data.UserAgent
	device.Cookie = new_device.Data.Cookie
	device.DId = new_device.Data.DeviceId
	device.IId = new_device.Data.InstallId
	device.DeviceToken = new_device.Data.DeviceToken
	device_info, err := json.Marshal(new_device.Data.DeviceInfo)

	if err != nil {
		fmt.Println(err)
		// return createNewDevice()
	}

	device.DeviceInfo = strings.Replace(string(device_info), "\"", "'", -1)

	helper.Database.Db.Create(&device)

	return device
}

func (ss *SafeDevice) addToDeviceInUse(user_id uint) {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	ss.device = append(ss.device, user_id)
	fmt.Println(ss.device)
}

func (ss *SafeDevice) removeFromDeviceInUse(user_id uint) {
	var result []uint
	ss.mu.Lock()
	defer ss.mu.Unlock()
	for _, v := range ss.device {
		if v != user_id {
			result = append(result, v)
		}
	}
	ss.device = result
	fmt.Println(ss.device)
}

type ErrorResponse struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}

func url_generator(path string) string {
	return fmt.Sprintf("http://127.0.0.1:%d/%s", rand.Intn(max_ports-min_ports+1)+min_ports, path)
}

func createAccountHelper(user_id uint, amount int) uint {
	device := models.Device{}
	if len(devices_in_use.device) > 0 {
		helper.Database.Db.Select("id", "device_info", "blocked").Where("id NOT IN (?) AND blocked = 0", devices_in_use.device).First(&device)
	} else {
		helper.Database.Db.Select("id", "device_info", "blocked").Where("blocked = 0").First(&device)
	}
	fmt.Println(device)
	if device.Id == 0 {
		device = createNewDevice()
	}

	go devices_in_use.addToDeviceInUse(device.Id)

	url := url_generator("register")
	fmt.Println(url)
	jsonBody := []byte(fmt.Sprintf(`{"device_data": "%s"}`, device.DeviceInfo))

	bodyReader := bytes.NewReader(jsonBody)

	req, err := http.NewRequest(http.MethodPost, url, bodyReader)
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		return 0
	}
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{
		Timeout: 30 * time.Minute,
	}

	res, err := client.Do(req)
	devices_in_use.removeFromDeviceInUse(device.Id)
	if err != nil {
		fmt.Printf("client: error making http request: %s\n", err)
		return 0
	}

	if res.StatusCode != http.StatusOK {
		if amount > 2 {
			return 0
		}
		var errorRes ErrorResponse
		decoder := json.NewDecoder(res.Body)
		if err := decoder.Decode(&errorRes); err != nil {
			fmt.Println("Error decoding JSON response:", err)
			return 0
		}

		if len(strings.Split(errorRes.Error, "|||||Blocked!|||||")) > 1 {
			BlockDevice(device.Id)
		}

		return createAccountHelper(user_id, amount+1)
	}

	var accountData models.Account
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&accountData); err != nil {
		fmt.Println("Error decoding JSON response:", err)
		return 0
	}

	accountData.UserId = user_id

	models.CreateAccount(&accountData, helper.Database.Db)

	return accountData.Id
}

func UseAccount(user_ids []uint, key string) {
	AddNewlyCreatedUserIdTOOtherEvents(user_ids, key)
}

func CreateAccount(amount int, user_id uint, automation_key string) {
	// Send request to Account Creation Service

	i := 0

	account_ids := []uint{}
	for i < amount {
		i++

		new_account := createAccountHelper(user_id, 0)
		if new_account == 0 {
			break
		}
		account_ids = append(account_ids, new_account)
	}
	AddNewlyCreatedUserIdTOOtherEvents(account_ids, automation_key)

}

func upload(session, tiktok_user_id, type_of_post, post_id, path, desc, music string, amount int) string {
	device := models.Device{}
	if len(devices_in_use.device) > 0 {
		helper.Database.Db.Select("id", "device_info", "blocked").Where("id NOT IN (?) AND blocked = 0", devices_in_use.device).First(&device)
	} else {
		helper.Database.Db.Select("id", "device_info", "blocked").Where("blocked = 0").First(&device)
	}
	fmt.Println(device)
	if device.Id == 0 {
		device = createNewDevice()
	}

	go devices_in_use.addToDeviceInUse(device.Id)
	path = strings.Join(strings.Split(path, "\\"), "/")
	url := url_generator("upload")
	jsonBody := []byte(fmt.Sprintf(`{"session":"%s","tiktok_user_id":"%s","type_of_post":"%s","post_id":"%s","path":"%s","desc":"%s","music":"%s", "device_data": "%s"}`, session, tiktok_user_id, type_of_post, post_id, path, desc, music, device.DeviceInfo))

	bodyReader := bytes.NewReader(jsonBody)

	req, err := http.NewRequest(http.MethodPost, url, bodyReader)
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		return ""
	}
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{
		Timeout: 30 * time.Minute,
	}

	res, err := client.Do(req)
	devices_in_use.removeFromDeviceInUse(device.Id)
	if err != nil {
		fmt.Printf("client: error making http request: %s\n", err)
		if amount >= 2 {
			return ""
		}
		return upload(session, type_of_post, post_id, path, desc, music, tiktok_user_id, amount+1)
	}

	if res.StatusCode != http.StatusOK {
		if amount >= 2 {
			return ""
		}

		var errorRes ErrorResponse
		decoder := json.NewDecoder(res.Body)
		if err := decoder.Decode(&errorRes); err != nil {
			fmt.Println("Error decoding JSON response:", err)
			return ""
		}

		if len(strings.Split(errorRes.Error, "|||||Blocked!|||||")) > 1 {
			BlockDevice(device.Id)
		}
		return upload(session, type_of_post, post_id, path, desc, music, tiktok_user_id, amount+1)
	}

	var accountData TiktokUploadResponse
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&accountData); err != nil {
		fmt.Println("Error decoding JSON response:", err)
		return upload(session, type_of_post, post_id, path, desc, music, tiktok_user_id, amount+1)
	}

	return accountData.PostId
}

func PostToTikTok(post_id int, user_id uint, account_ids []uint, automation_key string) {
	accounts := []models.Account{}

	helper.Database.Db.Select("id", "tik_user_id", "session", "cleared", "is_banned").Where("cleared = ? AND is_banned = ?", 0, 0).Find(&accounts, account_ids)

	post := models.Post{}

	helper.Database.Db.Where("id = ?", post_id).First(&post)

	for _, account := range accounts {
		helper.Database.Db.Create(&models.AccountPost{
			UserId:    user_id,
			PostID:    post.Id,
			AccountId: account.Id,
			TikId:     upload(account.Session, account.TikUserId, post.Type, post.Path, fullPath, post.Desc, post.Music, 0),
		})
	}

}

func CheckBanned(tik_user_id string, account_id uint, posts []models.AccountPost) bool {
	mep, banned := helper.IsAccBanned(tik_user_id)

	fmt.Println(mep)

	if banned {
		helper.Database.Db.Model(&models.Account{}).Where("tik_user_id = ?", tik_user_id).Update("is_banned", true)
		return true
	}

	sum := 0
	post_len := len(posts) - 1
	for key, num := range mep["video-views"] {
		if key <= post_len {
			fmt.Println("Error Here")
			helper.Database.Db.Model(&models.AccountPost{}).Where("id = ?", posts[key].Id).Update("total_views", num)
			fmt.Println("No Error Here")
		}
		sum += num
	}

	fmt.Println("Error Here 1")
	helper.Database.Db.Model(&models.Account{}).Where("id = ?", account_id).Update("total_likes", mep["likes-count"][0]).Update("total_views", sum).Update("followers", mep["followers-count"][0])
	fmt.Println("No Error Here 1")

	return false
}

func RefreshAccountValue(tik_user_id, post_id string, account_id uint, posts []models.AccountPost) {
	if !CheckBanned(tik_user_id, account_id, posts) {
		mep := helper.GetVideoData(tik_user_id, post_id)
		fmt.Println(mep)
		fmt.Println("Error Here 2")
		if mep["like-count"] != nil {
			helper.Database.Db.Model(&models.AccountPost{}).Where("tik_id = ?", post_id).Update("total_likes", mep["like-count"][0])
		}
	}
}

func ClearAccount(account_ids []uint, key string) {
	helper.Database.Db.Model(&models.Account{}).Where("id IN ?", account_ids).Update("cleared", 1)

	RemoveUserIdTOOtherEvents(account_ids, key)
}

func RefreshAccount(accounts_uids []uint) {
	accounts := []models.Account{}
	helper.Database.Db.Select("id", "tik_user_id").Where("id IN ?", accounts_uids).Find(&accounts)
	fmt.Println(accounts)
	for _, account := range accounts {
		posts := []models.AccountPost{}
		helper.Database.Db.Order("id desc").Select("id", "tik_id").Where("account_id = ?", account.Id).Find(&posts)

		for _, post := range posts {
			RefreshAccountValue(account.TikUserId, post.TikId, account.Id, posts)
		}

	}
}
