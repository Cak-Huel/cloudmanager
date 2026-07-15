package storage

import (
	"github.com/cak-huel/cloudmanager/v2/auth"
	"github.com/cak-huel/cloudmanager/v2/settings"
	"github.com/cak-huel/cloudmanager/v2/share"
	"github.com/cak-huel/cloudmanager/v2/users"
)

// Storage is a storage powered by a Backend which makes the necessary
// verifications when fetching and saving data to ensure consistency.
type Storage struct {
	Users    users.Store
	Share    *share.Storage
	Auth     *auth.Storage
	Settings *settings.Storage
}
