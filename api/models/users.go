package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	DiscordID       string                  `gorm:"unique;not null" json:"discord_id"`
	BitField        uint                    `json:"bit_field"`
	DiscordIdentity *map[string]interface{} `gorm:"type:jsonb;default:null" json:"discord_identity"`

	Coins  uint     `json:"coins"`
	Guilds []Guilds `gorm:"foreignKey:UserID"`

	GuildsJoined uint `json:"guilds_joined"`
	GuildsLeft   uint `json:"guilds_left"`

	History []History `gorm:"foreignKey:UserID"`
	Reviews []Reviews `gorm:"foreignKey:UserID"`
}

func (u *User) TableName() string {
	return "users.identity"
}

func (u *User) GetDiscordIdentity() map[string]interface{} {
	return *u.DiscordIdentity
}

type Guilds struct {
	gorm.Model
	UserID uint `json:"user_id"`
	User   User `gorm:"foreignKey:UserID"`

	GuildID uint                    `json:"guild_id"`
	Guild   *map[string]interface{} `gorm:"type:jsonb;default:null" json:"guild"`

	BitField uint `json:"bit_field"` // Will be used to determine if the state of the guild when showing it.

	InviteURL string `json:"invite_url"`

	LastMessageTime time.Time `json:"last_message_time"`

	Statistics []Statistics `gorm:"foreignKey:GuildID"`

	Reviews []Reviews `gorm:"foreignKey:GuildID"`
}

func (g *Guilds) TableName() string {
	return "users.guilds"
}

func (g *Guilds) GetGuild() map[string]interface{} {
	return *g.Guild
}

type Statistics struct {
	gorm.Model
	GuildID uint   `json:"guild_id"`
	Guild   Guilds `gorm:"foreignKey:GuildID"`

	Messages    uint `json:"messages"`
	Users       uint `json:"users"`
	ActiveUsers uint `json:"active_users"`
}

func (s *Statistics) TableName() string {
	return "users.statistics"
}

type Platform string

const (
	Paypal Platform = "paypal"
	Stripe Platform = "stripe"
)

type Action string

const (
	Payment          Action = "payment"
	AskingRefund     Action = "asking_refund"
	Refund           Action = "refund"
	GuildJoin        Action = "guild_join"
	GuildLeave       Action = "guild_leave"
	GuildPromote     Action = "guild_promote"
	UserJoined       Action = "user_joined"
	UserLeft         Action = "user_left"
	Review           Action = "review"
	Authentification Action = "authentification"
	Disconnect       Action = "disconnect"
	CreateGuild      Action = "create_guild"
	DeleteGuild      Action = "delete_guild"
	EditGuild        Action = "edit_guild"
	CoinAdd          Action = "coin_add"
	CoinRemove       Action = "coin_remove"
	AdminAction      Action = "admin_action"
)

type History struct {
	gorm.Model
	UserID uint `json:"user_id"`
	User   User `gorm:"foreignKey:UserID"`

	Action    Action                  `json:"action"`
	Reason    string                  `json:"reason"`
	Amount    *uint                   `json:"amount"`
	Platform  *Platform               `json:"platform"`
	OrderID   *string                 `json:"order_id"`
	RawObject *map[string]interface{} `gorm:"type:jsonb;default:null" json:"raw_object"`
}

func (h *History) TableName() string {
	return "users.history"
}

type Reviews struct {
	gorm.Model
	UserID uint `json:"user_id"`
	User   User `gorm:"foreignKey:UserID"`

	GuildID uint   `json:"guild_id"`
	Guild   Guilds `gorm:"foreignKey:GuildID"`

	Stars   uint   `json:"stars"`
	Comment string `json:"comment"`
}

func (r *Reviews) TableName() string {
	return "users.reviews"
}
