package settingstore

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"testing"

	_ "modernc.org/sqlite"
)

func initDB() (*sql.DB, error) {
	dsn := ":memory:?parseTime=true"
	db, err := sql.Open("sqlite", dsn)

	if err != nil {
		return nil, err
	}

	return db, nil
}

func initStore() (StoreInterface, error) {
	db, err := initDB()

	if err != nil {
		return nil, err
	}

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		SettingTableName:   "setting",
		AutomigrateEnabled: true,
	})

	if err != nil {
		return nil, err
	}

	if store == nil {
		return nil, errors.New("unexpected nil store")
	}

	return store, nil
}

func TestStore_Create(t *testing.T) {
	store, err := initStore()

	if err != nil {
		t.Fatal("Store could not be created: ", err.Error())
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}
}

func TestStore_Automigrate(t *testing.T) {
	store, err := initStore()

	if err != nil {
		t.Fatal("Store could not be created: ", err.Error())
	}

	err = store.MigrateUp(context.Background())

	if err != nil {
		t.Fatal("Automigrate failed: " + err.Error())
	}
}

func TestStore_EnableDebug(t *testing.T) {
	store, err := initStore()

	if err != nil {
		t.Fatal("Store could not be created: ", err.Error())
	}

	store.EnableDebug(true)

	err = store.MigrateUp(context.Background())

	if err != nil {
		t.Fatal("Automigrate failed: " + err.Error())
	}
}

func TestStore_SettingCreate(t *testing.T) {
	store, err := initStore()

	if err != nil {
		t.Fatal("Store could not be created: ", err.Error())
	}

	setting := NewSetting().
		SetKey("1").
		SetValue("one two three")

	if setting == nil {
		t.Fatal("unexpected nil setting")
	}

	if setting.GetID() == "" {
		t.Fatal("unexpected empty id:", setting.GetID())
	}

	if len(setting.GetID()) < 8 {
		t.Fatal("unexpected id length:", len(setting.GetID()))
	}

	if setting.GetKey() == "" {
		t.Fatal("unexpected empty key:", setting.GetKey())
	}

	if len(setting.GetKey()) != 1 {
		t.Fatal("unexpected key length:", len(setting.GetKey()))
	}

	err = store.SettingCreate(context.Background(), setting)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}
}

func TestStore_SettingDelete(t *testing.T) {
	store, err := initStore()

	if err != nil {
		t.Fatal("Store could not be created: ", err.Error())
	}

	setting := NewSetting().
		SetKey("1").
		SetValue("one two three")

	err = store.SettingCreate(context.Background(), setting)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.SettingDeleteByID(context.Background(), setting.GetID())

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	settingFindWithDeleted, err := store.SettingList(context.Background(), SettingQuery().
		SetID(setting.GetID()).
		SetLimit(1).
		SetSoftDeletedIncluded(true))

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if len(settingFindWithDeleted) != 0 {
		t.Fatal("Setting MUST be deleted, but it is not")
	}
}

func TestStore_SettingDeleteByID(t *testing.T) {
	store, err := initStore()

	if err != nil {
		t.Fatal("Store could not be created: ", err.Error())
	}

	setting := NewSetting().
		SetKey("1").
		SetValue("one two three")

	if setting == nil {
		t.Fatal("unexpected nil setting")
	}

	if setting.GetID() == "" {
		t.Fatal("unexpected empty id:", setting.GetID())
	}

	err = store.SettingCreate(context.Background(), setting)

	if err != nil {
		t.Error("unexpected error:", err)
	}

	err = store.SettingDeleteByID(context.Background(), setting.GetID())

	if err != nil {
		t.Error("unexpected error:", err)
	}

	settingFindWithDeleted, err := store.SettingList(context.Background(), SettingQuery().
		SetID(setting.GetID()).
		SetLimit(1).
		SetSoftDeletedIncluded(true))

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if len(settingFindWithDeleted) != 0 {
		t.Fatal("Setting MUST be deleted, but it is not")
	}
}

func TestStore_SettingDeleteByKey(t *testing.T) {
	store, err := initStore()

	if err != nil {
		t.Fatal("Store could not be created: ", err.Error())
	}

	setting := NewSetting().
		SetKey("1").
		SetValue("one two three")

	if setting == nil {
		t.Fatal("unexpected nil setting")
	}

	if setting.GetID() == "" {
		t.Fatal("unexpected empty id:", setting.GetID())
	}

	err = store.SettingCreate(context.Background(), setting)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.SettingDeleteByKey(context.Background(), setting.GetKey())

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	settingFindWithDeleted, err := store.SettingList(context.Background(), SettingQuery().
		SetKey(setting.GetKey()).
		SetLimit(1).
		SetSoftDeletedIncluded(true))

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if len(settingFindWithDeleted) != 0 {
		t.Fatal("Setting MUST be deleted, but it is not")
	}
}

func TestStore_SettingFindByID(t *testing.T) {
	store, err := initStore()

	if err != nil {
		t.Fatal("Store could not be created: ", err.Error())
	}

	setting := NewSetting().
		SetKey("1").
		SetValue("one two three")

	if setting == nil {
		t.Fatal("unexpected nil setting")
	}

	if setting.GetID() == "" {
		t.Fatal("unexpected empty id:", setting.GetID())
	}

	err = store.SettingCreate(context.Background(), setting)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	settingFound, err := store.SettingFindByID(context.Background(), setting.GetID())

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if settingFound == nil {
		t.Fatal("Setting MUST be found, but it is not")
	}
}

func TestStore_SettingFindByKey(t *testing.T) {
	store, err := initStore()

	if err != nil {
		t.Fatal("Store could not be created: ", err.Error())
	}

	setting := NewSetting().
		SetKey("1").
		SetValue("one two three")

	if setting == nil {
		t.Fatal("unexpected nil setting")
	}

	if setting.GetID() == "" {
		t.Fatal("unexpected empty id:", setting.GetID())
	}

	err = store.SettingCreate(context.Background(), setting)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	settingFound, err := store.SettingFindByKey(context.Background(), setting.GetKey())

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if settingFound == nil {
		t.Fatal("Setting MUST be found, but it is not")
	}
}

func TestStore_SettingList(t *testing.T) {
	store, err := initStore()

	if err != nil {
		t.Fatal("Store could not be created: ", err.Error())
	}

	setting := NewSetting().
		SetKey("1").
		SetValue("one two three")

	if setting == nil {
		t.Fatal("unexpected nil setting")
	}

	if setting.GetID() == "" {
		t.Fatal("unexpected empty id:", setting.GetID())
	}

	err = store.SettingCreate(context.Background(), setting)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	settingFound, err := store.SettingList(context.Background(), SettingQuery().
		SetKey(setting.GetKey()).
		SetLimit(1))

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if len(settingFound) != 1 {
		t.Fatal("Setting MUST be found, but it is not")
	}

	if settingFound[0].GetID() != setting.GetID() {
		t.Fatalf("Setting ID is wrong, expected %s, got %s", setting.GetID(), settingFound[0].GetID())
	}
}

func TestStore_SettingUpdate(t *testing.T) {
	store, err := initStore()

	if err != nil {
		t.Fatal("Store could not be created: ", err.Error())
	}

	setting := NewSetting().
		SetKey("1").
		SetValue("one two three")

	if setting == nil {
		t.Fatal("unexpected nil setting")
	}

	if setting.GetID() == "" {
		t.Fatal("unexpected empty id:", setting.GetID())
	}

	err = store.SettingCreate(context.Background(), setting)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	setting.SetValue("four five six")

	err = store.SettingUpdate(context.Background(), setting)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	settingFound, err := store.SettingFindByID(context.Background(), setting.GetID())

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if settingFound == nil {
		t.Fatal("Setting MUST be found, but it is not")
	}

	if settingFound.GetValue() != "four five six" {
		t.Fatalf("Setting Value is wrong, expected %s, got %s", "four five six", settingFound.GetValue())
	}
}

func TestStore_SettingCount(t *testing.T) {
	store, err := initStore()

	if err != nil {
		t.Fatal("Store could not be created: ", err.Error())
	}

	setting := NewSetting().
		SetKey("1").
		SetValue("one two three")

	if setting == nil {
		t.Fatal("unexpected nil setting")
	}

	if setting.GetID() == "" {
		t.Fatal("unexpected empty id:", setting.GetID())
	}

	err = store.SettingCreate(context.Background(), setting)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	count, err := store.SettingCount(context.Background(), SettingQuery().
		SetKey(setting.GetKey()))

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if count != 1 {
		t.Fatalf("Setting count is wrong, expected %d, got %d", 1, count)
	}
}

func TestStore_SettingSoftDelete(t *testing.T) {
	store, err := initStore()

	if err != nil {
		t.Fatal("Store could not be created: ", err.Error())
	}

	setting := NewSetting().
		SetKey("1").
		SetValue("one two three")

	if setting == nil {
		t.Fatal("unexpected nil setting")
	}

	if setting.GetID() == "" {
		t.Fatal("unexpected empty id:", setting.GetID())
	}

	err = store.SettingCreate(context.Background(), setting)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.SettingSoftDelete(context.Background(), setting)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	settingFindWithDeleted, err := store.SettingList(context.Background(), SettingQuery().
		SetID(setting.GetID()).
		SetLimit(1).
		SetSoftDeletedIncluded(true))

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if len(settingFindWithDeleted) != 1 {
		t.Fatal("Setting MUST be found, but it is not")
	}

	if !settingFindWithDeleted[0].IsSoftDeleted() {
		t.Fatal("Setting MUST be soft deleted")
	}
}

func TestStore_SettingSoftDeleteByID(t *testing.T) {
	store, err := initStore()

	if err != nil {
		t.Fatal("Store could not be created: ", err.Error())
	}

	setting := NewSetting().
		SetKey("1").
		SetValue("one two three")

	if setting == nil {
		t.Fatal("unexpected nil setting")
	}

	if setting.GetID() == "" {
		t.Fatal("unexpected empty id:", setting.GetID())
	}

	err = store.SettingCreate(context.Background(), setting)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.SettingSoftDeleteByID(context.Background(), setting.GetID())

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	settingFindWithDeleted, err := store.SettingList(context.Background(), SettingQuery().
		SetID(setting.GetID()).
		SetLimit(1).
		SetSoftDeletedIncluded(true))

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if len(settingFindWithDeleted) != 1 {
		t.Fatal("Setting MUST be found, but it is not")
	}

	if !settingFindWithDeleted[0].IsSoftDeleted() {
		t.Fatal("Setting MUST be soft deleted")
	}
}

func TestStore_Get(t *testing.T) {
	store, err := initStore()

	if err != nil {
		t.Fatal("Store could not be created: ", err.Error())
	}

	setting := NewSetting().
		SetKey("1").
		SetValue("one two three")

	if setting == nil {
		t.Fatal("unexpected nil setting")
	}

	if setting.GetID() == "" {
		t.Fatal("unexpected empty id:", setting.GetID())
	}

	err = store.SettingCreate(context.Background(), setting)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	value, err := store.Get(context.Background(), setting.GetKey(), "")

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if value != "one two three" {
		t.Fatalf("Setting Value is wrong, expected %s, got %s", "one two three", value)
	}
}

func TestStore_GetAny(t *testing.T) {
	store, err := initStore()

	if err != nil {
		t.Fatal("Store could not be created: ", err.Error())
	}

	setting := NewSetting().
		SetKey("1").
		SetValue(`{"key1": "value1", "key2": "value2"}`)

	if setting == nil {
		t.Fatal("unexpected nil setting")
	}

	if setting.GetID() == "" {
		t.Fatal("unexpected empty id:", setting.GetID())
	}

	err = store.SettingCreate(context.Background(), setting)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	value, err := store.GetAny(context.Background(), setting.GetKey(), nil)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if value == nil {
		t.Fatal("Setting Value is wrong, expected not nil, got nil")
	}
}

func TestStore_GetMap(t *testing.T) {
	store, err := initStore()

	if err != nil {
		t.Fatal("Store could not be created: ", err.Error())
	}

	setting := NewSetting().
		SetKey("1").
		SetValue(`{"key1": "value1", "key2": "value2"}`)

	if setting == nil {
		t.Fatal("unexpected nil setting")
	}

	if setting.GetID() == "" {
		t.Fatal("unexpected empty id:", setting.GetID())
	}

	err = store.SettingCreate(context.Background(), setting)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	value, err := store.GetMap(context.Background(), setting.GetKey(), nil)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if value == nil {
		t.Fatal("Setting Value is wrong, expected not nil, got nil")
	}

	if value["key1"] != "value1" {
		t.Fatalf("Setting Value is wrong, expected %s, got %s", "value1", value["key1"])
	}
}

func TestStore_Set(t *testing.T) {
	store, err := initStore()

	if err != nil {
		t.Fatal("Store could not be created: ", err.Error())
	}

	setting := NewSetting().
		SetKey("1").
		SetValue("one two three")

	if setting == nil {
		t.Fatal("unexpected nil setting")
	}

	if setting.GetID() == "" {
		t.Fatal("unexpected empty id:", setting.GetID())
	}

	err = store.SettingCreate(context.Background(), setting)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.Set(context.Background(), setting.GetKey(), "four five six")

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	value, err := store.Get(context.Background(), setting.GetKey(), "")

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if value != "four five six" {
		t.Fatalf("Setting Value is wrong, expected %s, got %s", "four five six", value)
	}
}

func TestStore_SetAny(t *testing.T) {
	store, err := initStore()

	if err != nil {
		t.Fatal("Store could not be created: ", err.Error())
	}

	setting := NewSetting().
		SetKey("1").
		SetValue(`{"key1": "value1", "key2": "value2"}`)

	if setting == nil {
		t.Fatal("unexpected nil setting")
	}

	if setting.GetID() == "" {
		t.Fatal("unexpected empty id:", setting.GetID())
	}

	err = store.SettingCreate(context.Background(), setting)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.SetAny(context.Background(), setting.GetKey(), map[string]string{
		"key1": "value1",
		"key2": "value2",
	}, 0)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	value, err := store.GetAny(context.Background(), setting.GetKey(), nil)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if value == nil {
		t.Fatal("Setting Value is wrong, expected not nil, got nil")
	}
}

func TestStore_SetMap(t *testing.T) {
	store, err := initStore()

	if err != nil {
		t.Fatal("Store could not be created: ", err.Error())
	}

	setting := NewSetting().
		SetKey("1").
		SetValue(`{"key1": "value1", "key2": "value2"}`)

	if setting == nil {
		t.Fatal("unexpected nil setting")
	}

	if setting.GetID() == "" {
		t.Fatal("unexpected empty id:", setting.GetID())
	}

	err = store.SettingCreate(context.Background(), setting)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.SetMap(context.Background(), setting.GetKey(), map[string]any{
		"key1": "value1",
		"key2": "value2",
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	value, err := store.GetMap(context.Background(), setting.GetKey(), nil)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if value == nil {
		t.Fatal("Setting Value is wrong, expected not nil, got nil")
	}

	if value["key1"] != "value1" {
		t.Fatalf("Setting Value is wrong, expected %s, got %s", "value1", value["key1"])
	}
}

func TestStore_MergeMap(t *testing.T) {
	store, err := initStore()

	if err != nil {
		t.Fatal("Store could not be created: ", err.Error())
	}

	setting := NewSetting().
		SetKey("1").
		SetValue(`{"key1": "value1", "key2": "value2"}`)

	if setting == nil {
		t.Fatal("unexpected nil setting")
	}

	if setting.GetID() == "" {
		t.Fatal("unexpected empty id:", setting.GetID())
	}

	err = store.SettingCreate(context.Background(), setting)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.MergeMap(context.Background(), setting.GetKey(), map[string]any{
		"key2": "value22",
		"key3": "value3",
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	value, err := store.GetMap(context.Background(), setting.GetKey(), nil)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if value == nil {
		t.Fatal("Setting Value is wrong, expected not nil, got nil")
	}

	if value["key1"] != "value1" {
		t.Fatalf("Setting Value is wrong, expected %s, got %s", "value1", value["key1"])
	}

	if value["key2"] != "value22" {
		t.Fatalf("Setting Value is wrong, expected %s, got %s", "value22", value["key2"])
	}

	if value["key3"] != "value3" {
		t.Fatalf("Setting Value is wrong, expected %s, got %s", "value3", value["key3"])
	}
}

func TestStore_Has(t *testing.T) {
	store, err := initStore()

	if err != nil {
		t.Fatal("Store could not be created: ", err.Error())
	}

	setting := NewSetting().
		SetKey("1").
		SetValue("one two three")

	if setting == nil {
		t.Fatal("unexpected nil setting")
	}

	if setting.GetID() == "" {
		t.Fatal("unexpected empty id:", setting.GetID())
	}

	err = store.SettingCreate(context.Background(), setting)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	has, err := store.Has(context.Background(), setting.GetKey())

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if !has {
		t.Fatal("Setting MUST be found, but it is not")
	}
}

func TestStore_HasNotFound(t *testing.T) {
	store, err := initStore()

	if err != nil {
		t.Fatal("Store could not be created: ", err.Error())
	}

	has, err := store.Has(context.Background(), "not-found")

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if has {
		t.Fatal("Setting MUST NOT be found, but it is")
	}
}

func TestStore_Delete(t *testing.T) {
	store, err := initStore()

	if err != nil {
		t.Fatal("Store could not be created: ", err.Error())
	}

	setting := NewSetting().
		SetKey("1").
		SetValue("one two three")

	if setting == nil {
		t.Fatal("unexpected nil setting")
	}

	if setting.GetID() == "" {
		t.Fatal("unexpected empty id:", setting.GetID())
	}

	err = store.SettingCreate(context.Background(), setting)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.Delete(context.Background(), setting.GetKey())

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	has, err := store.Has(context.Background(), setting.GetKey())

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if has {
		t.Fatal("Setting MUST NOT be found, but it is")
	}
}

func TestStore_MigrateDown(t *testing.T) {
	store, err := initStore()

	if err != nil {
		t.Fatal("Store could not be created: ", err.Error())
	}

	err = store.MigrateDown(context.Background())

	if err != nil {
		t.Fatal("MigrateDown failed: " + err.Error())
	}

	// Table should no longer exist, so MigrateUp should recreate it
	err = store.MigrateUp(context.Background())

	if err != nil {
		t.Fatal("MigrateUp after MigrateDown failed: " + err.Error())
	}
}

func TestStore_SettingSoftDelete_ExcludesFromDefaultQuery(t *testing.T) {
	store, err := initStore()

	if err != nil {
		t.Fatal("Store could not be created: ", err.Error())
	}

	setting := NewSetting().
		SetKey("1").
		SetValue("one two three")

	err = store.SettingCreate(context.Background(), setting)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.SettingSoftDelete(context.Background(), setting)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	// Default query should NOT find soft-deleted records
	settingFound, err := store.SettingFindByID(context.Background(), setting.GetID())

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if settingFound != nil {
		t.Fatal("soft-deleted Setting MUST NOT be found by default query")
	}

	// Query with SoftDeletedIncluded should find it
	settingFindWithDeleted, err := store.SettingList(context.Background(), SettingQuery().
		SetID(setting.GetID()).
		SetLimit(1).
		SetSoftDeletedIncluded(true))

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if len(settingFindWithDeleted) != 1 {
		t.Fatalf("expected 1 soft-deleted setting, got %d", len(settingFindWithDeleted))
	}

	if !settingFindWithDeleted[0].IsSoftDeleted() {
		t.Fatal("Setting MUST be soft deleted")
	}
}

func TestStore_SettingList_MultipleWithFiltering(t *testing.T) {
	store, err := initStore()

	if err != nil {
		t.Fatal("Store could not be created: ", err.Error())
	}

	setting1 := NewSetting().SetKey("1").SetValue("one two three")
	setting2 := NewSetting().SetKey("2").SetValue("four five six")
	setting3 := NewSetting().SetKey("3").SetValue("seven eight nine")

	for _, setting := range []SettingInterface{setting1, setting2, setting3} {
		err = store.SettingCreate(context.Background(), setting)
		if err != nil {
			t.Fatal("unexpected error:", err)
		}
	}

	// Filter by key "2" should return exactly 1 result
	list, err := store.SettingList(context.Background(), SettingQuery().
		SetKey("2").
		SetLimit(2))

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if len(list) != 1 {
		t.Fatalf("expected 1 setting, got %d", len(list))
	}

	if list[0].GetKey() != "2" {
		t.Fatalf("expected key '2', got %s", list[0].GetKey())
	}
}

func TestStore_SettingList_OrderByAndOffset(t *testing.T) {
	store, err := initStore()

	if err != nil {
		t.Fatal("Store could not be created: ", err.Error())
	}

	setting1 := NewSetting().SetKey("a").SetValue("1")
	setting2 := NewSetting().SetKey("b").SetValue("2")
	setting3 := NewSetting().SetKey("c").SetValue("3")

	for _, setting := range []SettingInterface{setting1, setting2, setting3} {
		err = store.SettingCreate(context.Background(), setting)
		if err != nil {
			t.Fatal("unexpected error:", err)
		}
	}

	// Order by setting_key ASC with offset 1, limit 1
	list, err := store.SettingList(context.Background(), SettingQuery().
		SetOrderBy(COLUMN_SETTING_KEY).
		SetSortOrder("ASC").
		SetOffset(1).
		SetLimit(1))

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if len(list) != 1 {
		t.Fatalf("expected 1 setting, got %d", len(list))
	}

	if list[0].GetKey() != "b" {
		t.Fatalf("expected key 'b' with offset 1, got %s", list[0].GetKey())
	}
}

func TestStore_SettingList_IDIn(t *testing.T) {
	store, err := initStore()

	if err != nil {
		t.Fatal("Store could not be created: ", err.Error())
	}

	setting1 := NewSetting().SetKey("a").SetValue("1")
	setting2 := NewSetting().SetKey("b").SetValue("2")
	setting3 := NewSetting().SetKey("c").SetValue("3")

	for _, setting := range []SettingInterface{setting1, setting2, setting3} {
		err = store.SettingCreate(context.Background(), setting)
		if err != nil {
			t.Fatal("unexpected error:", err)
		}
	}

	list, err := store.SettingList(context.Background(), SettingQuery().
		SetIDIn([]string{setting1.GetID(), setting3.GetID()}).
		SetLimit(10))

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if len(list) != 2 {
		t.Fatalf("expected 2 settings, got %d", len(list))
	}
}

func TestStore_SettingCount_WithQuery(t *testing.T) {
	store, err := initStore()

	if err != nil {
		t.Fatal("Store could not be created: ", err.Error())
	}

	setting1 := NewSetting().SetKey("count-test").SetValue("1")
	setting2 := NewSetting().SetKey("count-test").SetValue("2")
	setting3 := NewSetting().SetKey("other").SetValue("3")

	for _, setting := range []SettingInterface{setting1, setting2, setting3} {
		err = store.SettingCreate(context.Background(), setting)
		if err != nil {
			t.Fatal("unexpected error:", err)
		}
	}

	count, err := store.SettingCount(context.Background(), SettingQuery().
		SetKey("count-test"))

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if count != 2 {
		t.Fatalf("expected count 2, got %d", count)
	}
}

func TestStore_SettingCreate_NilSetting(t *testing.T) {
	store, err := initStore()

	if err != nil {
		t.Fatal("Store could not be created: ", err.Error())
	}

	err = store.SettingCreate(context.Background(), nil)

	if err == nil {
		t.Fatal("expected error for nil setting")
	}
}

func TestStore_SettingCreate_EmptyKey(t *testing.T) {
	store, err := initStore()

	if err != nil {
		t.Fatal("Store could not be created: ", err.Error())
	}

	setting := NewSetting().SetKey("").SetValue("value")

	err = store.SettingCreate(context.Background(), setting)

	if err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestStore_SettingFindByID_EmptyID(t *testing.T) {
	store, err := initStore()

	if err != nil {
		t.Fatal("Store could not be created: ", err.Error())
	}

	_, err = store.SettingFindByID(context.Background(), "")

	if err == nil {
		t.Fatal("expected error for empty id")
	}
}

func TestStore_SettingFindByKey_EmptyKey(t *testing.T) {
	store, err := initStore()

	if err != nil {
		t.Fatal("Store could not be created: ", err.Error())
	}

	_, err = store.SettingFindByKey(context.Background(), "")

	if err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestStore_Has_EmptyKey(t *testing.T) {
	store, err := initStore()

	if err != nil {
		t.Fatal("Store could not be created: ", err.Error())
	}

	_, err = store.Has(context.Background(), "")

	if err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestStore_SettingQueryValidateErrors(t *testing.T) {
	testCases := []struct {
		name     string
		setup    func(SettingQueryInterface)
		contains string
	}{
		{
			name: "created_at_gte empty",
			setup: func(q SettingQueryInterface) {
				q.SetCreatedAtGte("")
			},
			contains: "created_at_gte cannot be empty",
		},
		{
			name: "created_at_lte empty",
			setup: func(q SettingQueryInterface) {
				q.SetCreatedAtLte("")
			},
			contains: "created_at_lte cannot be empty",
		},
		{
			name: "id empty",
			setup: func(q SettingQueryInterface) {
				q.SetID("")
			},
			contains: "id cannot be empty",
		},
		{
			name: "id_in empty",
			setup: func(q SettingQueryInterface) {
				q.SetIDIn([]string{})
			},
			contains: "id_in cannot be empty array",
		},
		{
			name: "key empty",
			setup: func(q SettingQueryInterface) {
				q.SetKey("")
			},
			contains: "key cannot be empty",
		},
		{
			name: "limit negative",
			setup: func(q SettingQueryInterface) {
				q.SetLimit(-1)
			},
			contains: "limit cannot be negative",
		},
		{
			name: "offset negative",
			setup: func(q SettingQueryInterface) {
				q.SetOffset(-1)
			},
			contains: "offset cannot be negative",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			query := NewSettingQuery()
			tc.setup(query)

			err := query.Validate()
			if err == nil {
				t.Fatalf("expected error containing %q", tc.contains)
			}

			if !strings.Contains(err.Error(), tc.contains) {
				t.Fatalf("unexpected error %q", err.Error())
			}
		})
	}
}
