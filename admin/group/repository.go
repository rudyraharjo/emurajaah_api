package group

import (
	"github.com/rudyraharjo/emurojaah/models"
)

// Repository Group
type Repository interface {
	GetListGroups() ([]models.GroupWithStatus, error)
	GetListGroupMember(groupID int) ([]models.ResponseGroupMember, error)
	GetLastNoIndexGroup(category string) ([]models.Group, error)
	CreateGroup(g models.Group) ([]models.Group, error)
	DeleteDuplicateGroupMembers() (int, error)
	GetUserBelomBaca(groupID int) ([]models.GroupMemberReadNotIsDone, error)
	InsertHistory(History models.History) error
}
