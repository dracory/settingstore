package settingstore

type SettingQueryInterface interface {
	Validate() error

	IsCountOnly() bool

	Columns() []string
	SetColumns(columns []string) SettingQueryInterface

	HasCreatedAtGte() bool
	CreatedAtGte() string
	SetCreatedAtGte(createdAtGte string) SettingQueryInterface

	HasCreatedAtLte() bool
	CreatedAtLte() string
	SetCreatedAtLte(createdAtLte string) SettingQueryInterface

	HasID() bool
	ID() string
	SetID(id string) SettingQueryInterface

	HasIDIn() bool
	IDIn() []string
	SetIDIn(idIn []string) SettingQueryInterface

	HasKey() bool
	Key() string
	SetKey(key string) SettingQueryInterface

	HasOffset() bool
	Offset() int
	SetOffset(offset int) SettingQueryInterface

	HasLimit() bool
	Limit() int
	SetLimit(limit int) SettingQueryInterface

	HasSortOrder() bool
	SortOrder() string
	SetSortOrder(sortOrder string) SettingQueryInterface

	HasOrderBy() bool
	OrderBy() string
	SetOrderBy(orderBy string) SettingQueryInterface

	HasCountOnly() bool
	SetCountOnly(countOnly bool) SettingQueryInterface

	HasSoftDeletedIncluded() bool
	SoftDeletedIncluded() bool
	SetSoftDeletedIncluded(withSoftDeleted bool) SettingQueryInterface
}
