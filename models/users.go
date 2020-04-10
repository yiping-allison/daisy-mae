package models

import (
	"fmt"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

const (
	// MaxQueue limits the amount of queues a user can join
	MaxQueue int = 3

	// MaxEvent limits the amount of events a user can create
	MaxEvent int = 1

	// MaxTrade limits the amount of trades a user can create
	MaxTrade int = 2
)

// UserService wraps to the User interface
type UserService interface {
	User
}

// User defines all the methods we can use to interact with
// the User Service
type User interface {
	AddTrade(user *discordgo.User, tradeID string)
	UserExists(user *discordgo.User) bool
	AddUser(user *discordgo.User) error
	AddEvent(user *discordgo.User, eventID string, expire time.Time)
	LimitEvent(user *discordgo.User) bool
	RemoveEvent(user *discordgo.User, eventID string)
	Clean()
	AddQueue(eventID string, user *discordgo.User, expire time.Time)
	LimitQueue(user *discordgo.User) bool
	RemoveQueue(eventID string, user *discordgo.User)
	RemoveAllQueue(eventID string)
	LimitTrade(userID string) bool
}

type userService struct {
	User
}

type userStore struct {
	user map[string]*UserData
	m    *sync.RWMutex
}

type tEvent struct {
	eventID string
	expire  time.Time
}

type tQueue struct {
	eventID string
	expire  time.Time
}

// Trades defines a trade event
type Trades struct {
	tradeID string
}

// UserData represents user data bot needs to keep track of
// in order to perform event and queue services
type UserData struct {
	events []tEvent
	queues []tQueue
	trades []Trades
}

// internal check to see if interface is implemented correctly
var _ User = &userStore{}

func (us userStore) AddTrade(user *discordgo.User, tradeID string) {
	// TODO: Finish this
}

// LimitTrade returns true when the user has reached
// the max amount of trade creation
func (us userStore) LimitTrade(userID string) bool {
	us.m.Lock()
	defer us.m.Unlock()
	fmt.Println(userID)
	val := us.user[userID]
	if len(val.trades) == 0 {
		return false
	}
	return len(val.trades) == MaxTrade
}

// RemoveAllQueue will remove all users with a certain eventID
func (us userStore) RemoveAllQueue(eventID string) {
	us.m.Lock()
	defer us.m.Unlock()
	var ret []tQueue
	for k, v := range us.user {
		for _, id := range v.queues {
			if id.eventID == eventID {
				continue
			}
			ret = append(ret, id)
		}
		us.user[k].queues = ret
	}
}

// LimitQueue returns true when user hit max queue amount
//
// Otherwise, it return false
func (us userStore) LimitQueue(user *discordgo.User) bool {
	us.m.RLock()
	defer us.m.RUnlock()
	if !us.UserExists(user) {
		return false
	}
	return len(us.user[user.ID].queues) == MaxQueue
}

// RemoveQueue will remove an item from the queue slice
func (us userStore) RemoveQueue(eventID string, user *discordgo.User) {
	us.m.Lock()
	defer us.m.Unlock()
	var ret []tQueue
	for _, v := range us.user[user.ID].queues {
		if v.eventID == eventID {
			continue
		}
		ret = append(ret, v)
	}
	us.user[user.ID].queues = ret
}

// AddQueue adds an item to the queue
func (us userStore) AddQueue(eventID string, user *discordgo.User, expire time.Time) {
	us.m.Lock()
	defer us.m.Unlock()
	new := tQueue{
		eventID: eventID,
		expire:  expire,
	}
	val := us.user[user.ID]
	val.queues = append(val.queues, new)
}

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
func removeEvent(eventID string, events []tEvent) []tEvent {
	var ret []tEvent
	for _, i := range events {
		if i.eventID == eventID {
			continue
		}
		ret = append(ret, i)
	}
	return ret
}

// LimitEvent returns true if the user has an event list equal to the max event
//
// In otherwords, if this is true, the user should not be able to
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
func (us userStore) AddEvent(user *discordgo.User, eventID string, expire time.Time) {
	us.m.Lock()
	defer us.m.Unlock()
	val := us.user[user.ID].events
	new := tEvent{
		eventID: eventID,
		expire:  expire,
	}
	val = append(val, new)
	us.user[user.ID].events = val
}

// AddUser will add a new user to the map along with initialized slices
// for events and queues tracking
func (us userStore) AddUser(user *discordgo.User) error {
	us.m.Lock()
	defer us.m.Unlock()
	e := make([]tEvent, 0)
	q := make([]tQueue, 0)
	t := make([]Trades, 0)
	data := UserData{
		events: e,
		queues: q,
		trades: t,
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

// Clean will remove event listings from tracking that have exceeded time limit
//
// DO NOT CALL THIS RANDOMLY!!
//
// This should only be called in the goroutine in main (ticker to check expiration)
func (us userStore) Clean() {
	us.m.Lock()
	defer us.m.Unlock()
	var newE []tEvent
	var newQ []tQueue
	for k, v := range us.user {
		for _, event := range v.events {
			if time.Now().Sub(event.expire) > 0 {
				// remove this event
				continue
			}
			newE = append(newE, event)
		}
		us.user[k].events = newE
		for _, queue := range v.queues {
			if time.Now().Sub(queue.expire) > 0 {
				// remove this queue
				continue
			}
			newQ = append(newQ, queue)
		}
		us.user[k].queues = newQ
	}
}
