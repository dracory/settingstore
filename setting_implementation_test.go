package settingstore

import (
    "testing"

    "github.com/dracory/sb"
)

func TestNewSetting_Defaults(t *testing.T) {
    setting := NewSetting()

    if setting == nil {
        t.Fatal("NewSetting returned nil")
    }

    if id := setting.GetID(); id == "" {
        t.Fatal("expected ID to be generated")
    } else if len(id) != 32 {
        t.Fatalf("unexpected ID length: %d", len(id))
    }

    if key := setting.GetKey(); key != "" {
        t.Fatalf("expected empty key, got %q", key)
    }

    if value := setting.GetValue(); value != "" {
        t.Fatalf("expected empty value, got %q", value)
    }

    if created := setting.GetCreatedAt(); created == "" {
        t.Fatal("expected created_at to be set")
    }

    if updated := setting.GetUpdatedAt(); updated == "" {
        t.Fatal("expected updated_at to be set")
    }

    if deleted := setting.GetSoftDeletedAt(); deleted != sb.MAX_DATETIME {
        t.Fatalf("expected soft_deleted_at to default to %q, got %q", sb.MAX_DATETIME, deleted)
    }

    if setting.IsSoftDeleted() {
        t.Fatal("expected new setting not to be soft deleted")
    }
}

func TestNewSettingFromExistingData(t *testing.T) {
    data := map[string]string{
        COLUMN_ID:              "foo-id",
        COLUMN_SETTING_KEY:     "feature",
        COLUMN_SETTING_VALUE:   "enabled",
        COLUMN_CREATED_AT:      "2024-01-02 03:04:05",
        COLUMN_UPDATED_AT:      "2024-01-03 04:05:06",
        COLUMN_SOFT_DELETED_AT: "2024-02-01 00:00:00",
    }

    setting := NewSettingFromExistingData(data)

    if setting.GetID() != data[COLUMN_ID] {
        t.Fatalf("unexpected id: %s", setting.GetID())
    }

    if setting.GetKey() != data[COLUMN_SETTING_KEY] {
        t.Fatalf("unexpected key: %s", setting.GetKey())
    }

    if setting.GetValue() != data[COLUMN_SETTING_VALUE] {
        t.Fatalf("unexpected value: %s", setting.GetValue())
    }

    if setting.GetCreatedAt() != data[COLUMN_CREATED_AT] {
        t.Fatalf("unexpected created_at: %s", setting.GetCreatedAt())
    }

    if setting.GetUpdatedAt() != data[COLUMN_UPDATED_AT] {
        t.Fatalf("unexpected updated_at: %s", setting.GetUpdatedAt())
    }

    if setting.GetSoftDeletedAt() != data[COLUMN_SOFT_DELETED_AT] {
        t.Fatalf("unexpected soft_deleted_at: %s", setting.GetSoftDeletedAt())
    }
}

func TestSetting_IsSoftDeleted(t *testing.T) {
    setting := NewSetting()

    if setting.IsSoftDeleted() {
        t.Fatal("expected default setting to not be soft deleted")
    }

    past := "2000-01-01 00:00:00"
    setting.SetSoftDeletedAt(past)

    if !setting.IsSoftDeleted() {
        t.Fatal("expected setting with past soft_deleted_at to be soft deleted")
    }
}
