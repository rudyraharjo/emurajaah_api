package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/rudyraharjo/emurojaah/group"
	"github.com/rudyraharjo/emurojaah/models"
	"github.com/spf13/viper"
)

type groupService struct {
	groupRepo group.Repository
}

// NewGroupService Service
func NewGroupService(repo group.Repository) group.Service {
	return &groupService{repo}
}

func (s *groupService) JoinGroupByEmail(groupID int, email string, grouptype string) (*models.GroupMember, []models.Group, bool, int) {

	dataNewMember, code, errIns := s.groupRepo.AddGroupMemberByEmail(groupID, email, grouptype)
	if errIns != nil {
		fmt.Println(errIns)
		return nil, nil, false, 1
	}

	Group, _ := s.groupRepo.GetNoGroupIndex(groupID)

	if code == 1 {
		return nil, nil, false, 1
	}

	if code == 2 {
		return nil, nil, false, 2
	}

	if code == 3 {
		return nil, nil, false, 3
	}

	return dataNewMember, Group, true, code

}

func (s *groupService) JoinBulkGroup(reqBulk models.RequestJoinBulkGroup) ([]models.Group, bool, int) {

	code, err := s.groupRepo.AddGroupMemberBulk(reqBulk)
	if err != nil {
		return nil, false, 1
	}

	if code == 1 {
		return nil, false, 1
	}

	if code == 2 {
		return nil, false, 2
	}

	Group, _ := s.groupRepo.GetNoGroupIndex(2)
	return Group, true, 0
}

func (s *groupService) JoinGroupNew(userId int, category string) (*models.GroupMember, []models.Group, bool, int) {
	dataNewMember, code, err := s.groupRepo.AddGroupMemberNew(userId, strings.Title(category))
	if err != nil {
		fmt.Println(err)
		return nil, nil, false, 1
	}

	if code == 1 {
		return nil, nil, false, 1
	}

	if code == 2 {
		return nil, nil, false, 2
	}

	Group, _ := s.groupRepo.GetNoGroupIndex(dataNewMember.GroupID)

	return dataNewMember, Group, true, 0
}

func (s *groupService) JoinGroup(userId int, category string) (*models.GroupMember, []models.Group, bool) {

	var nMembers *models.GroupMember
	now := time.Now()
	reservedIndex := 0
	categoryType := strings.Title(category)

	// GetAvailableGroupByType With limitUser
	grp := s.GetAvailableGroupByType(categoryType)
	//ListGroupMembers, err := s.groupRepo.GetListGroupMembersByIDGroup(grp.Id)

	for {

		genNumberIndex := generateNumber(categoryType)

		arrGroupMembers, err := s.groupRepo.GetListGroupMembersByIDGroupAndCurrenIndex(grp.Id, genNumberIndex)

		if err != nil {
			fmt.Println("error on check ListGroupMembersByIDGroup")
			return nil, nil, false
		}

		if arrGroupMembers == 0 {
			reservedIndex = genNumberIndex
			break
		}

	}

	newMember := models.GroupMember{
		CurrentIndex: reservedIndex,
		UserID:       userId,
		GroupID:      grp.Id,
		GroupType:    categoryType,
		IsDone:       0,
		CreatedAt:    now,
	}

	dataNewMember, errIns := s.groupRepo.AddGroupMember(newMember, grp.Id)
	if errIns != nil {
		fmt.Println(errIns)
		return nil, nil, false
	}

	Group, _ := s.groupRepo.GetNoGroupIndex(dataNewMember.GroupID)

	nMembers = dataNewMember

	return nMembers, Group, true

}

func generateNumber(categoryType string) int {

	min := 1
	max := 0
	if categoryType == "Ayat" {
		max = 6236
	} else if categoryType == "Halaman" {
		max = 604
	} else {
		max = 30
	}

	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min+1) + min
}

func (s *groupService) GetAvailableGroupByType(category string) *models.Group {

	groups, err := s.groupRepo.ListAvailableGroupByType(category)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	//create group if not exist
	if len(groups) == 0 {

		getgroups, err := s.groupRepo.GetLastNoIndexGroup(category)

		if err != nil {
			fmt.Println(err)
			return nil
		}

		return s.CreateGroup(category, getgroups[0].NoGroupIndex)
	}

	g := groups[0]

	return &g
}

// func (s *groupService) UpdateInactiveMember(groupId int, oldUserId int, newUserId int) (*models.GroupMember, error) {
// 	newMember, err := s.groupRepo.UpdateInactiveMember(groupId, oldUserId, newUserId)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return newMember, nil
// }

func (s *groupService) CreateGroup(category string, noindex int) *models.Group {

	now := time.Now()
	groupName := ""
	limitUser := 0
	switch category {
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
		NoGroupIndex:  noindex + 1,
		GroupName:     groupName,
		MaxMember:     limitUser,
		CurrentMember: 0,
		GroupType:     category,
		CreatedAt:     now,
	}

	groupID, err := s.groupRepo.CreateGroup(newGroup)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	newGroup.Id = groupID
	return &newGroup
}

func (s *groupService) IsMemberExist(groupId int, userId int) bool {
	code, err := s.groupRepo.IsMemberExist(groupId, userId)
	if err != nil {
		return false
	}
	if code == 2 {
		return false
	}
	return true
}

func (s *groupService) GetListUserGroup(userId int) ([]models.ResponsGroupeGroupMembersListWithNoUrut, error) {
	data, err := s.groupRepo.ListUserGroup(userId)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return data, nil
}

func (s *groupService) GetListGroupMembersByUserIDAndGroupID(UserID int, groupID int) int {
	Count, err := s.groupRepo.GetListGroupMembersByUserIDAndGroupID(UserID, groupID)
	if err != nil {
		return 0
	}
	return Count
}

func (s *groupService) GetListAllGroup() ([]models.ResponseAllgroup, error) {

	data, err := s.groupRepo.ListAllGroup()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return data, nil

}

func (s *groupService) GetListGroupType(UserID int, GroupType string, Offset int, Limit int) ([]models.ResponseListGroupTypeWithNoUrut, error) {
	datalist, err := s.groupRepo.GetListGroupType(UserID, GroupType, Offset, Limit)
	if err != nil {
		fmt.Println(err)
		return nil, nil
	}

	return datalist, nil
}

func (s *groupService) GetListGroupMemberWithName(groupId int) []models.ResponseGroupMemberList {
	members, err := s.groupRepo.ListGroupMemberWithName(groupId)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return members
}

func (s *groupService) GetListGroupMemberWithNameAndOffsetLimit(groupId int, offset int, limit int) []models.ResponseGroupMemberList {
	members, err := s.groupRepo.ListGroupMemberWithNameAndOffsetLimit(groupId, offset, limit)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return members
}

func (s *groupService) GetListMemberGroupByStatusAndPaging(groupId int, offset int, limit int, status int) []models.ResponseGroupMemberList {
	members, err := s.groupRepo.ListGroupMemberByStatusAndPaging(groupId, offset, limit, status)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return members
}

func (s *groupService) TotalGroupMemberByStatus(groupId int, status int) int {
	count, err := s.groupRepo.TotalGroupMemberByStatus(groupId, status)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	return count
}

func (s *groupService) TotalGroupMember(groupId int) int {
	count, err := s.groupRepo.TotalGroupMember(groupId)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	return count
}

func (s *groupService) LeaveReadingGroup(params models.RequestLeaveReading) error {

	err := s.groupRepo.LeaveReadingGroupAndInsertActivity(params)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil

}

func (s *groupService) ExitGroup(groupID int, userID int) error {

	err := s.groupRepo.ExitGroupAndInsertActivity(groupID, userID)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil

	// e1 := s.SetGroupMemberAsInactive(groupId, userId)
	// if e1 != nil {
	// 	return e1
	// }

	// return nil
}

func (s *groupService) SetGroupMemberAsInactive(groupId int, userId int) error {
	err := s.groupRepo.SetGroupMemberAsInactive(groupId, userId)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (s *groupService) GetAllGroupMembersWithType() []models.GroupMemberWithType {
	data, err := s.groupRepo.GetAllGroupMembersWithType()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return data
}

func (s *groupService) UpdateGroupMemberReadingStatusToZero() error {
	err := s.groupRepo.UpdateGroupMemberReadingStatusToZero()
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (s *groupService) HandleUpdateMemberReadingIndex() {

	// err1 := s.UpdateGroupMemberReadingStatusToZero()
	// if err1 != nil {
	// 	return
	// }

	members := s.GetAllGroupMembersWithType()

	maxIndex := 0
	for _, member := range members {
		switch strings.Title(member.GroupType) {

		case "Halaman":
			maxIndex = 604
			break
		case "Juz":
			maxIndex = 30
			break
		case "Ayat":
			maxIndex = 6236
			break
		}

		if member.CurrentIndex == maxIndex {
			s.UpdateGroupMemberReadingIndex(member.Id, 1)
		} else {
			s.UpdateGroupMemberReadingIndex(member.Id, member.CurrentIndex+1)
		}
	}

}

func (s *groupService) UpdateGroupMemberReadingIndex(ID int, index int) {
	err := s.groupRepo.UpdateGroupMemberReadingIndex(ID, index)
	if err != nil {
		fmt.Println(err)
	}
}

func (s *groupService) GetAllReadNotIsDone() []models.GroupMemberReadNotIsDone {

	data, err := s.groupRepo.GetAllGroupMembersReadNotIsDone()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return data

}

func (s *groupService) HandleNotifUserReadIsNotDone() {

	member_check_read := s.GetAllReadNotIsDone()

	fmt.Println("Count Member Beloman Baca => ", len(member_check_read))

	if len(member_check_read) > 0 {

		//GetToken := []string{}

		url_fcm := viper.GetString(`fcm_key.url`)
		serverKey := viper.GetString(`fcm_key.server_key`)

		notification := models.NotificationSendFcm{
			Title:    "Emurojaah Pengingat ..",
			Priority: "high",
			Body:     "Assalamualaikum..! Jangan lupa untuk melanjutkan bacaannya ya..",
		}

		for _, member := range member_check_read {

			//GetToken = append(GetToken, ""+member.TokenFirebase+"")

			PostJSON := models.RequestSingleSendFcm{
				To:           member.TokenFirebase,
				Notification: notification,
			}

			payloadByte, _ := json.Marshal(PostJSON)

			// fmt.Println(string(payloadByte))
			// fmt.Print("\n")

			var payload = bytes.NewReader(payloadByte)
			req, errReq := http.NewRequest("POST", url_fcm, payload)
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

}
