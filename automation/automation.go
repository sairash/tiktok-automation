package automation

import (
	"beepbop/helper"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/alphadose/haxmap"
	"golang.org/x/exp/slices"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

type TimeStruct struct {
	NextEventTime time.Time
	ReturnHashKey string
	ByUser        uint
}

type AutomationLog struct {
	TypeOfAutomation string
	Amount           int
	NeededAccountIds []uint
	UserId           uint
	Executing        bool
}

var TimeMap = haxmap.New[string, TimeStruct]()

var AutomationLogMap = haxmap.New[string, [][]AutomationLog]()

var UserAllPendingAutomationLog = haxmap.New[uint, []string]()

const time_wait_duration = 2 * time.Second

func TimeSeries() {
	for {
		time.Sleep(time_wait_duration)
		currentTime := time.Now()

		TimeMap.ForEach(func(key string, time_data TimeStruct) bool {
			if currentTime.After(time_data.NextEventTime) {
				go ExecuteFirstEventFromLog(time_data.ReturnHashKey)
				go RemoveFromUserPending(key, time_data.ByUser)
				TimeMap.Del(key)
			}
			return true
		})
	}
}

func RemoveFromUserPending(key string, UserId uint) {
	all_pending_time, ok := UserAllPendingAutomationLog.Get(UserId)
	if ok {
		new_val := []string{}
		for _, pending_time := range all_pending_time {
			if pending_time != key {
				new_val = append(new_val, pending_time)
			}
		}
		UserAllPendingAutomationLog.Set(UserId, new_val)
	} else {
		fmt.Println("No Such Data!")
	}
}

func AddToUserPending(key string, UserId uint) {
	all_pending_time, ok := UserAllPendingAutomationLog.Get(UserId)
	if ok {
		all_pending_time = append(all_pending_time, key)
	} else {
		all_pending_time = []string{
			key,
		}
	}
	UserAllPendingAutomationLog.Set(UserId, all_pending_time)
}

func ExecuteFirstEventFromLog(key string) {
	val, ok := AutomationLogMap.Get(key)
	if ok {
		if len(val) > 0 {
			if len(val[0]) > 0 {
				val[0][0].Executing = true
				AutomationLogMap.Set(key, val)
				val[0][0].Executing = false

				switch val[0][0].TypeOfAutomation {
				case "use_account":
					UseAccount(val[0][0].NeededAccountIds, key)
				case "create_account":
					CreateAccount(val[0][0].Amount, val[0][0].UserId, key)
				case "clear_account":
					if len(val[0][0].NeededAccountIds) > 0 {
						ClearAccount(val[0][0].NeededAccountIds, key)
					}
				case "refresh_account":
					if len(val[0][0].NeededAccountIds) > 0 {
						RefreshAccount(val[0][0].NeededAccountIds, key)
					}
				case "delete_posts":
					if len(val[0][0].NeededAccountIds) > 0 {
						delete_all_posts(val[0][0].NeededAccountIds, key)
					}
				case "post":
					if len(val[0][0].NeededAccountIds) > 0 {
						PostToTikTok(val[0][0].Amount, val[0][0].UserId, val[0][0].NeededAccountIds, key)
					}
				}

				val[0][0].NeededAccountIds = []uint{}

				switch len(val) {
				case 1:
					add_automation_to_backlog := []AutomationLog{
						val[0][0],
					}
					val = append(val, add_automation_to_backlog)
					val[0] = val[0][1:]
				case 2:
					add_automation_to_backlog := val[0][0]
					val[1] = append(val[1], add_automation_to_backlog)
					val[0] = val[0][1:]
				}

				if len(val[0]) == 0 {
					if len(val) > 1 {
						next_array_len := len(val[1])
						if next_array_len < 1 || val[1][next_array_len-1].TypeOfAutomation != "repeat" {
							AutomationLogMap.Del(key)
						} else {
							val = val[1:]
							InitiateAutomation(key, val, true)
						}
					} else {
						AutomationLogMap.Del(key)
					}
				} else {
					InitiateAutomation(key, val, true)
				}
			}
		}
	}
}

func AddNewlyCreatedUserIdTOOtherEvents(user_id []uint, automationKey string) {
	if automationlogs, automation_ok := AutomationLogMap.Get(automationKey); automation_ok {
		if len(automationlogs) > 0 {
			for auto_key, automations := range automationlogs[0] {
				if automations.TypeOfAutomation == "post" || automations.TypeOfAutomation == "clear_account" || automations.TypeOfAutomation == "refresh_account" || automations.TypeOfAutomation == "delete_posts" {
					automationlogs[0][auto_key].NeededAccountIds = append(automationlogs[0][auto_key].NeededAccountIds, user_id...)
				}
			}
			AutomationLogMap.Set(automationKey, automationlogs)
		}
	}
}

func RemoveElementFromArray(user_ids []uint, remove_ids []uint) []uint {
	return_user_id := []uint{}

	for _, user_id := range user_ids {
		remove := false
		for _, remove_id := range remove_ids {
			if remove_id == user_id {
				remove = true
			}
		}
		if !remove {
			return_user_id = append(return_user_id, user_id)
		}
	}

	return return_user_id
}

func RemoveUserIdTOOtherEvents(user_ids []uint, automationKey string) {
	if automationlogs, automation_ok := AutomationLogMap.Get(automationKey); automation_ok {
		if len(automationlogs) > 0 {
			for auto_key, automations := range automationlogs[0] {
				if automations.TypeOfAutomation == "post" || automations.TypeOfAutomation == "clear_account" || automations.TypeOfAutomation == "refresh_account" || automations.TypeOfAutomation == "delete_posts" {
					automationlogs[0][auto_key].NeededAccountIds = RemoveElementFromArray(automationlogs[0][auto_key].NeededAccountIds, user_ids)
				}
			}
			AutomationLogMap.Set(automationKey, automationlogs)
		}
	}
}

func SetToTimeMap(min int, UserId uint, return_key, key string) {
	TimeMap.Set(key, TimeStruct{
		NextEventTime: time.Now().Add(time.Duration(min) * time.Second),
		ByUser:        UserId,
		ReturnHashKey: return_key,
	})
	AddToUserPending(return_key, UserId)
}

func InitiateAutomation(identifier string, autolog [][]AutomationLog, set_to_AutomationLog bool) {
	if len(autolog[0]) > 0 {
		time_identifier := helper.RandomString(20)
		if autolog[0][0].TypeOfAutomation != "wait" {
			SetToTimeMap(2, autolog[0][0].UserId, identifier, time_identifier)
		} else {
			SetToTimeMap(autolog[0][0].Amount, autolog[0][0].UserId, identifier, time_identifier)
			autolog[0][0].Executing = true
		}
		if set_to_AutomationLog {
			AutomationLogMap.Set(identifier, autolog)
		}
	}
}

func GetUserAutomation(user_id uint) map[string][]AutomationLog {
	return_automation := map[string][]AutomationLog{}

	if keys, ok := UserAllPendingAutomationLog.Get(user_id); ok {
		for _, value := range keys {
			if automationlogs, automation_ok := AutomationLogMap.Get(value); automation_ok {
				automation_array := []AutomationLog{}
				for _, automations := range automationlogs {
					automation_array = append(automation_array, automations...)
				}
				automation_array[0].Executing = true
				return_automation[value] = automation_array
			}
		}
	}
	return return_automation

}

func GenerateAutomationLog(type_of_automations, values []string, user_id uint) error {
	accounts_available := false
	some_actions_done := false
	new_automation_log := [][]AutomationLog{{}}
	accounts_to_work_on := []uint{}
	for key, type_of_automation := range type_of_automations {
		switch type_of_automation {
		case "post":
			if !accounts_available {
				return fmt.Errorf("Create Accounts or Select accounts First To do these automations.")
			}
			if amount, err := strconv.Atoi(values[key]); err == nil {
				new_automation_log[0] = append(new_automation_log[0], AutomationLog{TypeOfAutomation: type_of_automation, Amount: amount, UserId: user_id, NeededAccountIds: accounts_to_work_on})
				some_actions_done = true
			} else {
				return err
			}
		case "clear_account":
			if !accounts_available {
				return fmt.Errorf("Create Accounts or Select accounts First To do these automations.")
			}
			if amount, err := strconv.Atoi(values[key]); err == nil {
				new_automation_log[0] = append(new_automation_log[0], AutomationLog{TypeOfAutomation: type_of_automation, Amount: amount, UserId: user_id, NeededAccountIds: accounts_to_work_on})
				some_actions_done = true
			} else {
				return err
			}
		case "refresh_account":
			if !accounts_available {
				return fmt.Errorf("Create Accounts or Select accounts First To do these automations.")
			}
			if amount, err := strconv.Atoi(values[key]); err == nil {
				new_automation_log[0] = append(new_automation_log[0], AutomationLog{TypeOfAutomation: type_of_automation, Amount: amount, UserId: user_id, NeededAccountIds: accounts_to_work_on})
				some_actions_done = true
			} else {
				return err
			}
		case "use_account":
			accounts_available = true
			some_actions_done = true
			user_ids := []uint{}

			for _, user_id := range strings.Split(values[key], ",") {
				user_id_int, err := strconv.Atoi(user_id)
				if err == nil {
					user_ids = append(user_ids, uint(user_id_int))
				}
			}
			new_automation_log[0] = append(new_automation_log[0], AutomationLog{TypeOfAutomation: type_of_automation, Amount: 0, UserId: user_id, NeededAccountIds: user_ids})
		case "create_account":
			accounts_available = true
			if amount, err := strconv.Atoi(values[key]); err == nil {
				new_automation_log[0] = append(new_automation_log[0], AutomationLog{TypeOfAutomation: type_of_automation, Amount: amount, UserId: user_id})
				some_actions_done = true
			} else {
				return err
			}
		default:
			if amount, err := strconv.Atoi(values[key]); err == nil {
				new_automation_log[0] = append(new_automation_log[0], AutomationLog{TypeOfAutomation: type_of_automation, Amount: amount, UserId: user_id})
			} else {
				return err
			}
		}
	}
	if !some_actions_done {
		return fmt.Errorf("No Meaningful Actions Were Done!")
	}
	go InitiateAutomation(helper.RandomString(20), new_automation_log, true)
	return nil
}

func RemoveAutomation(user_id uint, key string) error {
	if running_automation_keys, ok := UserAllPendingAutomationLog.Get(user_id); ok {
		if slices.Contains(running_automation_keys, key) {
			AutomationLogMap.Del(key)
			TimeMap.ForEach(func(key string, time_data TimeStruct) bool {
				if time_data.ReturnHashKey == key {
					TimeMap.Del(key)
				}
				return true
			})
			return nil
		}
	}

	return fmt.Errorf("")
}
