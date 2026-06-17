package settingstore

import (
	"github.com/dracory/neat/database/orm"
	"github.com/dracory/neat/database/soft_delete"
	neatuid "github.com/dracory/neat/support/uid"
	"github.com/dromara/carbon/v2"
)

// SettingInterface defines the interface for a setting record.
type SettingInterface interface {
	// Methods

	IsSoftDeleted() bool

	// Setters and Getters

	GetCreatedAt() string
	GetCreatedAtCarbon() *carbon.Carbon
	SetCreatedAt(createdAt string) SettingInterface

	GetID() string
	SetID(id string) SettingInterface

	GetKey() string
	SetKey(key string) SettingInterface

	GetSoftDeletedAt() string
	GetSoftDeletedAtCarbon() *carbon.Carbon
	SetSoftDeletedAt(deletedAt string) SettingInterface

	GetUpdatedAt() string
	GetUpdatedAtCarbon() *carbon.Carbon
	SetUpdatedAt(updatedAt string) SettingInterface

	GetValue() string
	SetValue(value string) SettingInterface
}

var _ SettingInterface = (*settingImplementation)(nil)

// == TYPE =====================================================================

type settingImplementation struct {
	orm.ShortID

	KeyField   string `db:"setting_key"`
	ValueField string `db:"setting_value"`

	CreatedAtField orm.CreatedAt
	UpdatedAtField orm.UpdatedAt
	soft_delete.SoftDeletesMaxDate
}

// == CONSTRUCTORS =============================================================

func NewSetting() SettingInterface {
	o := &settingImplementation{}
	o.SetID(neatuid.GenerateShortID())
	o.SetKey("")
	o.SetValue("")
	o.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	o.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	o.SetSoftDeletedAt(MAX_DATETIME)
	return o
}

func NewSettingFromExistingData(data map[string]string) SettingInterface {
	o := &settingImplementation{}
	o.SetID(data[COLUMN_ID])
	o.SetKey(data[COLUMN_SETTING_KEY])
	o.SetValue(data[COLUMN_SETTING_VALUE])
	if v, ok := data[COLUMN_CREATED_AT]; ok {
		o.SetCreatedAt(v)
	}
	if v, ok := data[COLUMN_UPDATED_AT]; ok {
		o.SetUpdatedAt(v)
	}
	if v, ok := data[COLUMN_SOFT_DELETED_AT]; ok {
		o.SetSoftDeletedAt(v)
	}
	return o
}

// == METHODS ==================================================================

func (o *settingImplementation) IsSoftDeleted() bool {
	return o.SoftDeletesMaxDate.IsSoftDeleted()
}

// == SETTERS AND GETTERS ======================================================

func (o *settingImplementation) GetID() string {
	return o.ShortID.ID
}

func (o *settingImplementation) SetID(id string) SettingInterface {
	o.ShortID.ID = id
	return o
}

func (o *settingImplementation) GetKey() string {
	return o.KeyField
}

func (o *settingImplementation) SetKey(key string) SettingInterface {
	o.KeyField = key
	return o
}

func (o *settingImplementation) GetValue() string {
	return o.ValueField
}

func (o *settingImplementation) SetValue(value string) SettingInterface {
	o.ValueField = value
	return o
}

func (o *settingImplementation) GetCreatedAt() string {
	if o.CreatedAtField.CreatedAt.IsZero() {
		return ""
	}
	return carbon.CreateFromStdTime(o.CreatedAtField.CreatedAt).ToDateTimeString()
}

func (o *settingImplementation) GetCreatedAtCarbon() *carbon.Carbon {
	return carbon.CreateFromStdTime(o.CreatedAtField.CreatedAt)
}

func (o *settingImplementation) SetCreatedAt(createdAt string) SettingInterface {
	if createdAt == "" {
		return o
	}
	o.CreatedAtField.CreatedAt = carbon.Parse(createdAt, carbon.UTC).StdTime()
	return o
}

func (o *settingImplementation) GetUpdatedAt() string {
	if o.UpdatedAtField.UpdatedAt.IsZero() {
		return ""
	}
	return carbon.CreateFromStdTime(o.UpdatedAtField.UpdatedAt).ToDateTimeString()
}

func (o *settingImplementation) GetUpdatedAtCarbon() *carbon.Carbon {
	return carbon.CreateFromStdTime(o.UpdatedAtField.UpdatedAt)
}

func (o *settingImplementation) SetUpdatedAt(updatedAt string) SettingInterface {
	if updatedAt == "" {
		return o
	}
	o.UpdatedAtField.UpdatedAt = carbon.Parse(updatedAt, carbon.UTC).StdTime()
	return o
}

func (o *settingImplementation) GetSoftDeletedAt() string {
	if o.SoftDeletesMaxDate.SoftDeletedAt.IsZero() {
		return ""
	}
	return carbon.CreateFromStdTime(o.SoftDeletesMaxDate.SoftDeletedAt).ToDateTimeString()
}

func (o *settingImplementation) GetSoftDeletedAtCarbon() *carbon.Carbon {
	return carbon.CreateFromStdTime(o.SoftDeletesMaxDate.SoftDeletedAt)
}

func (o *settingImplementation) SetSoftDeletedAt(deletedAt string) SettingInterface {
	if deletedAt == "" {
		return o
	}
	o.SoftDeletesMaxDate.SoftDeletedAt = carbon.Parse(deletedAt, carbon.UTC).StdTime()
	return o
}
