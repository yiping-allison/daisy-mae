package models

import (
	"sync"

	"github.com/bwmarrin/discordgo"
)

const (
	// MaxQueue limits the amount of queues a user can join
	MaxQueue int = 3

	// MaxEvent limits the amount of events a user can create
	MaxEvent int = 1
)

// UserService wraps to the User interface
type UserService interface {
	User
}

// User defines all the methods we can use to interact with
// the User Service
type User interface {
	UserExists(user *discordgo.User) bool
	AddUser(user *discordgo.User) error
	AddEvent(user *discordgo.User, eventID string)
	LimitEvent(user *discordgo.User) bool
	RemoveEvent(user *discordgo.User, eventID string)
}

type userService struct {
	User
}

type userStore struct {
	user map[string]*UserData
	m    *sync.RWMutex
}

// UserData represents user data bot needs to keep track of
// in order to perform event and queue services
type UserData struct {
	events []string
	queues []string
}

// internal check to see if interface is implemented correctly
var _ User = &userStore{}

// RemoveEvent removes an event from a user's tracking state
func (us userStore) RemoveEvent(user *discordgo.User, eventID string) {
	us.m.Lock()
	defer us.m.Unlock()
	us.user[user.ID].events = removeEvent(eventID, us.user[user.ID].events)
}

// removeEvent will remove a specified event from the user tracking
// event list
//
// This is a helper func to be used with RemoveEvent
func removeEvent(eventID string, events []string) []string {
	var ret []string
	for _, i := range events {
		if i == eventID {
			continue
		}
		ret = append(ret, i)
	}
	return ret
}

// LimitEvent returns true if the user has an event list equal to the max event
//
// In otherwords, if this is true, the user should now be able to
// make more events
func (us userStore) LimitEvent(user *discordgo.User) bool {
	us.m.RLock()
	defer us.m.RUnlock()
	if !us.UserExists(user) {
		return false
	}
	return len(us.user[user.ID].events) == MaxEvent
}

// UserExists returns a bool indicating if the requested user is
// already in server tracking or not
func (us userStore) UserExists(user *discordgo.User) bool {
	us.m.RLock()
	defer us.m.RUnlock()
	if _, ok := us.user[user.ID]; ok {
		return true
	}
	return false
}

// AddEvent adds a new event to the user tracking events
func (us userStore) AddEvent(user *discordgo.User, eventID string) {
	us.m.Lock()
	defer us.m.Unlock()
	val := us.user[user.ID].events
	val = append(val, eventID)
	us.user[user.ID].events = val
}

// AddUser will add a new user to the map along with initialized slices
// for events and queues tracking
func (us userStore) AddUser(user *discordgo.User) error {
	us.m.Lock()
	defer us.m.Unlock()
	e := make([]string, 0)
	q := make([]string, 0)
	data := UserData{
		events: e,
		queues: q,
	}
	us.user[user.ID] = &data
	return nil
}

// NewUserService initializes a new user service
func NewUserService() UserService {
	return userService{
		User: userStore{
			user: make(map[string]*UserData),
			m:    &sync.RWMutex{},
		},
	}
}
