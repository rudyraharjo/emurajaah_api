package group

import (
	"github.com/rudyraharjo/emurojaah/models"
)

// Service Group
type Service interface {
	GetListGroups() []models.GroupWithStatus
	GetListGroupMember(groupID int) []models.ResponseGroupMember
	GenerateAddGroup(Type string) []models.Group
	DeleteDuplicateGroupMembers() (int, error)
	SendNotifBelomBaca(GroupID int) (int, error)
}
