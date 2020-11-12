package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/rudyraharjo/emurojaah/admin/group"
	"github.com/rudyraharjo/emurojaah/models"
	"github.com/spf13/viper"
)

type groupAdmService struct {
	groupRepo group.Repository
}

// NewAdmGroupService func
func NewAdmGroupService(repo group.Repository) group.Service {
	return &groupAdmService{repo}
}

func (s *groupAdmService) DeleteDuplicateGroupMembers() (int, error) {
	c, err := s.groupRepo.DeleteDuplicateGroupMembers()
	if c == 0 || err != nil {
		return 0, err
	}

	if c == 2 {
		return 2, nil
	}
	return 1, nil
}

func (s *groupAdmService) GetListGroups() []models.GroupWithStatus {
	data, err := s.groupRepo.GetListGroups()
	if err != nil {
		return nil
	}
	return data
}

func (s *groupAdmService) GetListGroupMember(groupID int) []models.ResponseGroupMember {
	data, err := s.groupRepo.GetListGroupMember(groupID)
	if err != nil {
		return nil
	}
	return data
}

func (s *groupAdmService) SendNotifBelomBaca(groupID int) (int, error) {

	Members, err := s.groupRepo.GetUserBelomBaca(groupID)

	if len(Members) > 0 {

		urlFcm := viper.GetString(`fcm_key.url`)
		serverKey := viper.GetString(`fcm_key.server_key`)

		notification := models.NotificationSendFcm{
			Title:    "Emurojaah Pengingat ..",
			Priority: "high",
			Body:     "Assalamualaikum..! Jangan lupa untuk melanjutkan bacaannya ya..",
		}

		for _, member := range Members {

			PostJSON := models.RequestSingleSendFcm{
				To:           member.TokenFirebase,
				Notification: notification,
			}

			payloadByte, _ := json.Marshal(PostJSON)

			var payload = bytes.NewReader(payloadByte)
			req, errReq := http.NewRequest("POST", urlFcm, payload)
			req.Header.Set("Authorization", "key="+serverKey)
			req.Header.Set("Content-Type", "application/json")

			if errReq != nil {
				//fmt.Println(errReq)
			}

			client := new(http.Client)
			response, errRes := client.Do(req)

			if errRes != nil {
				panic(errRes)
			}

			defer response.Body.Close()

			if response != nil {

				body, _ := ioutil.ReadAll(response.Body)

				History := models.History{
					Description: string(body),
					Process:     "Response Notif User ID " + strconv.Itoa(member.UserID) + "",
					CreatedAt:   time.Now(),
				}

				err := s.groupRepo.InsertHistory(History)
				if err != nil {
					fmt.Println(err)
				}

			}

		}
	}

	if err != nil {
		return 1, err
	}

	return 0, nil

}

func (s *groupAdmService) GenerateAddGroup(Type string) []models.Group {
	getgroups, err := s.groupRepo.GetLastNoIndexGroup(Type)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	now := time.Now()
	groupName := ""
	limitUser := 0
	switch strings.Title(Type) {
	case "Juz":
		limitUser = 30
		groupName = "Group Juz"
	case "Ayat":
		limitUser = 6236
		groupName = "Group Ayat"
	default:
		limitUser = 604
		groupName = "Group Halaman"
	}

	newGroup := models.Group{
		NoGroupIndex:  getgroups[0].NoGroupIndex + 1,
		GroupName:     groupName,
		MaxMember:     limitUser,
		CurrentMember: 0,
		GroupType:     Type,
		CreatedAt:     now,
	}

	groupNew, err := s.groupRepo.CreateGroup(newGroup)

	if err != nil {
		return nil
	}

	return groupNew
}
