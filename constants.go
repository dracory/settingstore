package settingstore

const (
	COLUMN_ID              = "id"
	COLUMN_SETTING_KEY     = "setting_key"
	COLUMN_SETTING_VALUE   = "setting_value"
	COLUMN_CREATED_AT      = "created_at"
	COLUMN_UPDATED_AT      = "updated_at"
	COLUMN_SOFT_DELETED_AT = "soft_deleted_at"
)

// MAX_DATETIME is a far-future datetime used as the default soft-delete sentinel.
const MAX_DATETIME = "9999-12-31 23:59:59"
