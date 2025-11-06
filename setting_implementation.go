package settingstore

import (
	"github.com/dracory/dataobject"
	"github.com/dracory/sb"
	"github.com/dracory/uid"
	"github.com/dromara/carbon/v2"
)

var _ SettingInterface = (*settingImplementation)(nil)

// Setting type
type settingImplementation struct {
	dataobject.DataObject
}

// == CONSTRUCTORS ============================================================

func NewSetting() SettingInterface {
	createdAt := carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)
	updatedAt := carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)
	deletedAt := sb.MAX_DATETIME

	o := (&settingImplementation{})

	o.SetID(uid.HumanUid()).
		SetKey("").
		SetValue("").
		SetCreatedAt(createdAt).
		SetUpdatedAt(updatedAt).
		SetSoftDeletedAt(deletedAt)

	return o
}

func NewSettingFromExistingData(data map[string]string) SettingInterface {
	o := &settingImplementation{}
	o.Hydrate(data)
	return o
}

// == METHODS =================================================================

func (o *settingImplementation) IsSoftDeleted() bool {
	return o.GetSoftDeletedAtCarbon().Compare("<", carbon.Now(carbon.UTC))
}

// == SETTERS AND GETTERS =====================================================

func (setting *settingImplementation) GetID() string {
	return setting.Get(COLUMN_ID)
}

func (setting *settingImplementation) SetID(id string) SettingInterface {
	setting.Set(COLUMN_ID, id)
	return setting
}

func (setting *settingImplementation) GetKey() string {
	return setting.Get(COLUMN_SETTING_KEY)
}

func (setting *settingImplementation) SetKey(key string) SettingInterface {
	setting.Set(COLUMN_SETTING_KEY, key)
	return setting
}

func (setting *settingImplementation) GetValue() string {
	return setting.Get(COLUMN_SETTING_VALUE)
}

func (setting *settingImplementation) SetValue(value string) SettingInterface {
	setting.Set(COLUMN_SETTING_VALUE, value)
	return setting
}

func (setting *settingImplementation) GetCreatedAt() string {
	return setting.Get(COLUMN_CREATED_AT)
}

func (setting *settingImplementation) GetCreatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(setting.GetCreatedAt(), carbon.UTC)
}

func (setting *settingImplementation) SetCreatedAt(createdAt string) SettingInterface {
	setting.Set(COLUMN_CREATED_AT, createdAt)
	return setting
}

func (setting *settingImplementation) GetUpdatedAt() string {
	return setting.Get(COLUMN_UPDATED_AT)
}

func (setting *settingImplementation) GetUpdatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(setting.GetUpdatedAt(), carbon.UTC)
}

func (setting *settingImplementation) SetUpdatedAt(updatedAt string) SettingInterface {
	setting.Set(COLUMN_UPDATED_AT, updatedAt)
	return setting
}

func (setting *settingImplementation) GetSoftDeletedAt() string {
	return setting.Get(COLUMN_SOFT_DELETED_AT)
}

func (setting *settingImplementation) GetSoftDeletedAtCarbon() *carbon.Carbon {
	return carbon.Parse(setting.GetSoftDeletedAt(), carbon.UTC)
}

func (setting *settingImplementation) SetSoftDeletedAt(deletedAt string) SettingInterface {
	setting.Set(COLUMN_SOFT_DELETED_AT, deletedAt)
	return setting
}
