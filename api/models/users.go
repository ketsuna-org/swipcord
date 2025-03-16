package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	DiscordID       string       `gorm:"unique;not null" json:"discord_id"`
	BitField        uint         `gorm:"default:0" json:"bit_field"`
	AccessToken     string       `json:"access_token"`
	RefreshToken    string       `json:"refresh_token"`
	ExpiresIn       time.Time    `json:"expires_in"`
	DiscordIdentity *DiscordUser `gorm:"type:jsonb;default:null" json:"discord_identity"`
	Coins           uint         `json:"coins"`
	Guilds          []Guilds     `gorm:"foreignKey:UserID"`

	GuildsJoined uint `json:"guilds_joined"`
	GuildsLeft   uint `json:"guilds_left"`

	History []History `gorm:"foreignKey:UserID"`
	Reviews []Reviews `gorm:"foreignKey:UserID"`
}

func (u *User) TableName() string {
	return "users.identity"
}

type DiscordUser struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	Discriminator string `json:"discriminator"`
	Avatar        string `json:"avatar"`
	Bot           bool   `json:"bot"`
	System        bool   `json:"system"`
	MFA           bool   `json:"mfa_enabled"`
	Locale        string `json:"locale"`
	Verified      bool   `json:"verified"`
	Email         string `json:"email"`
	Flags         uint   `json:"flags"`
	PremiumType   uint   `json:"premium_type"`
	PublicFlags   uint   `json:"public_flags"`
}

func (du *DiscordUser) Value() (driver.Value, error) {
	if du == nil {
		return nil, nil
	}
	return json.Marshal(du)
}

// Scan implémente l'interface sql.Scanner pour la désérialisation depuis JSONB
func (du *DiscordUser) Scan(value interface{}) error {
	if value == nil {
		*du = DiscordUser{}
		return nil
	}

	data, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(data, du)
}

func (u *User) GetDiscordIdentity() DiscordUser {
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
