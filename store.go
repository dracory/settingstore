package settingstore

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"os"
	"time"

	"github.com/dracory/neat"
	contractsorm "github.com/dracory/neat/contracts/database/orm"
	contractsschema "github.com/dracory/neat/contracts/database/schema"
	"github.com/dromara/carbon/v2"
)

// StoreInterface defines the interface for a setting store.
type StoreInterface interface {
	MigrateDown(ctx context.Context, tx ...*sql.Tx) error
	MigrateUp(ctx context.Context, tx ...*sql.Tx) error
	EnableDebug(debug bool)

	SettingCount(ctx context.Context, query SettingQueryInterface) (int64, error)
	SettingCreate(ctx context.Context, setting SettingInterface) error
	SettingDelete(ctx context.Context, setting SettingInterface) error
	SettingDeleteByID(ctx context.Context, settingID string) error
	SettingFindByID(ctx context.Context, settingID string) (SettingInterface, error)
	SettingFindByKey(ctx context.Context, settingKey string) (SettingInterface, error)
	SettingList(ctx context.Context, query SettingQueryInterface) ([]SettingInterface, error)
	SettingSoftDelete(ctx context.Context, setting SettingInterface) error
	SettingSoftDeleteByID(ctx context.Context, settingID string) error
	SettingUpdate(ctx context.Context, setting SettingInterface) error

	Delete(ctx context.Context, settingKey string) error
	Get(ctx context.Context, settingKey string, valueDefault string) (string, error)
	GetAny(ctx context.Context, key string, valueDefault any) (any, error)
	GetMap(ctx context.Context, key string, valueDefault map[string]any) (map[string]any, error)
	Has(ctx context.Context, settingKey string) (bool, error)
	MergeMap(ctx context.Context, key string, mergeMap map[string]any) error
	Set(ctx context.Context, settingKey string, value string) error
	SetAny(ctx context.Context, key string, value interface{}, seconds int64) error
	SetMap(ctx context.Context, key string, value map[string]any) error
	SettingDeleteByKey(ctx context.Context, settingKey string) error
}

var _ StoreInterface = (*storeImplementation)(nil)

// == TYPE =====================================================================

// Store defines a setting store
type storeImplementation struct {
	settingTableName   string
	db                 *neat.Database
	automigrateEnabled bool
	debugEnabled       bool
	logger             *slog.Logger
}

// PUBLIC METHODS ==============================================================

// MigrateUp creates the settings table if it does not exist
func (store *storeImplementation) MigrateUp(ctx context.Context, tx ...*sql.Tx) error {
	if store.db.Schema().HasTable(store.settingTableName) {
		if store.debugEnabled {
			store.logger.Info("MigrateUp: table already exists", "table", store.settingTableName)
		}
		return nil
	}

	err := store.db.Schema().Create(store.settingTableName, func(table contractsschema.Blueprint) {
		table.String(COLUMN_ID, 40)
		table.Primary(COLUMN_ID)
		table.String(COLUMN_SETTING_KEY, 255)
		table.Text(COLUMN_SETTING_VALUE)
		table.DateTime(COLUMN_CREATED_AT)
		table.DateTime(COLUMN_UPDATED_AT)
		table.DateTime(COLUMN_SOFT_DELETED_AT)
	})

	if err != nil {
		if store.debugEnabled {
			store.logger.Error("MigrateUp failed", "error", err)
		}
		return err
	}

	return nil
}

// MigrateDown drops the settings table
func (store *storeImplementation) MigrateDown(ctx context.Context, tx ...*sql.Tx) error {
	if !store.db.Schema().HasTable(store.settingTableName) {
		if store.debugEnabled {
			store.logger.Info("MigrateDown: table does not exist", "table", store.settingTableName)
		}
		return nil
	}

	err := store.db.Schema().Drop(store.settingTableName)
	if err != nil {
		if store.debugEnabled {
			store.logger.Error("MigrateDown failed", "error", err)
		}
		return err
	}
	return nil
}

// EnableDebug - enables the debug option
func (st *storeImplementation) EnableDebug(debug bool) {
	st.debugEnabled = debug
	if debug {
		st.db.EnableDebug()
		st.logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	} else {
		st.db.DisableDebug()
		st.logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	}
}

// Delete is a shortcut method to delete a value by key
func (st *storeImplementation) Delete(ctx context.Context, settingKey string) error {
	return st.SettingDeleteByKey(ctx, settingKey)
}

// Get is a shortcut method to get a value by key, or a default, if not found
func (st *storeImplementation) Get(ctx context.Context, settingKey string, valueDefault string) (string, error) {
	setting, errFindByKey := st.SettingFindByKey(ctx, settingKey)

	if errFindByKey != nil {
		return "", errFindByKey
	}

	if setting != nil {
		return setting.GetValue(), nil
	}

	return valueDefault, nil
}

// GetAny is a shortcut method to get a value by key as an interface, or a default if not found
func (st *storeImplementation) GetAny(ctx context.Context, key string, valueDefault any) (any, error) {
	setting, errFindByKey := st.SettingFindByKey(ctx, key)

	if errFindByKey != nil {
		return valueDefault, errFindByKey
	}

	if setting != nil {
		jsonValue := setting.GetValue()
		var val interface{}
		jsonError := json.Unmarshal([]byte(jsonValue), &val)
		if jsonError != nil {
			return valueDefault, jsonError
		}

		return val, nil
	}

	return valueDefault, nil
}

// GetMap is a shortcut method to get a value by key as a map, or a default if not found
func (st *storeImplementation) GetMap(ctx context.Context, key string, valueDefault map[string]any) (map[string]any, error) {
	setting, errFindByKey := st.SettingFindByKey(ctx, key)

	if errFindByKey != nil {
		return valueDefault, errFindByKey
	}

	if setting != nil {
		jsonValue := setting.GetValue()
		var val map[string]any
		jsonError := json.Unmarshal([]byte(jsonValue), &val)
		if jsonError != nil {
			return valueDefault, jsonError
		}

		return val, nil
	}

	return valueDefault, nil
}

// Has is a shortcut method to check if a setting exists by key
func (store *storeImplementation) Has(ctx context.Context, settingKey string) (bool, error) {
	if settingKey == "" {
		return false, errors.New("setting store > find by key: setting key is required")
	}

	query := SettingQuery().
		SetKey(settingKey).
		SetLimit(1)

	count, err := store.SettingCount(ctx, query)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// MergeMap is a shortcut method to merge a map with an existing map
func (st *storeImplementation) MergeMap(ctx context.Context, key string, mergeMap map[string]any) error {
	currentMap, err := st.GetMap(ctx, key, nil)

	if err != nil {
		return err
	}

	if currentMap == nil {
		return errors.New("settingstore. nil found")
	}

	for mapKey, mapValue := range mergeMap {
		currentMap[mapKey] = mapValue
	}

	return st.SetMap(ctx, key, currentMap)
}

// Set is a shortcut method to save a value by key, use Get to extract
func (st *storeImplementation) Set(ctx context.Context, settingKey string, value string) error {
	setting, errFindByKey := st.SettingFindByKey(ctx, settingKey)

	if errFindByKey != nil {
		return errFindByKey
	}

	if setting == nil {
		setting = NewSetting().SetKey(settingKey).SetValue(value)
		return st.SettingCreate(ctx, setting)
	}

	setting.SetValue(value)

	return st.SettingUpdate(ctx, setting)
}

// SetAny is a shortcut method to save any value by key, use GetAny to extract
func (st *storeImplementation) SetAny(ctx context.Context, key string, value interface{}, seconds int64) error {
	jsonValue, jsonError := json.Marshal(value)

	if jsonError != nil {
		return jsonError
	}

	return st.Set(ctx, key, string(jsonValue))
}

// SetMap is a shortcut method to save a map by key, use GetMap to extract
func (st *storeImplementation) SetMap(ctx context.Context, key string, value map[string]any) error {
	jsonValue, jsonError := json.Marshal(value)

	if jsonError != nil {
		return jsonError
	}

	return st.Set(ctx, key, string(jsonValue))
}

// SettingCount counts the settings based on the provided query.
func (store *storeImplementation) SettingCount(ctx context.Context, options SettingQueryInterface) (int64, error) {
	if options == nil {
		return 0, errors.New("setting query: cannot be nil")
	}

	if err := options.Validate(); err != nil {
		return 0, err
	}

	q := store.buildQuery(options)

	var count int64
	err := q.Table(store.settingTableName).Count(&count)
	return count, err
}

// SettingCreate creates a new setting
func (st *storeImplementation) SettingCreate(ctx context.Context, setting SettingInterface) error {
	if setting == nil {
		return errors.New("settingstore > setting create. setting cannot be nil")
	}

	if setting.GetKey() == "" {
		return errors.New("settingstore > setting create. key cannot be empty")
	}

	if setting.GetCreatedAt() == "" {
		setting.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString())
	}

	if setting.GetUpdatedAt() == "" {
		setting.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString())
	}

	if setting.GetSoftDeletedAt() == "" {
		setting.SetSoftDeletedAt(MAX_DATETIME)
	}

	row := map[string]any{
		COLUMN_ID:              setting.GetID(),
		COLUMN_SETTING_KEY:     setting.GetKey(),
		COLUMN_SETTING_VALUE:   setting.GetValue(),
		COLUMN_CREATED_AT:      setting.GetCreatedAtCarbon().StdTime(),
		COLUMN_UPDATED_AT:      setting.GetUpdatedAtCarbon().StdTime(),
		COLUMN_SOFT_DELETED_AT: setting.GetSoftDeletedAtCarbon().StdTime(),
	}

	return st.db.Query().Table(st.settingTableName).Create(row)
}

// SettingDelete deletes a setting
func (store *storeImplementation) SettingDelete(ctx context.Context, setting SettingInterface) error {
	if setting == nil {
		return errors.New("setting is nil")
	}

	return store.SettingDeleteByID(ctx, setting.GetID())
}

// SettingDeleteByID deletes a setting by id
func (store *storeImplementation) SettingDeleteByID(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("setting id is empty")
	}

	_, err := store.db.Query().
		Table(store.settingTableName).
		Where(COLUMN_ID+" = ?", id).
		Delete()

	return err
}

// SettingDeleteByKey deletes a setting by key
func (store *storeImplementation) SettingDeleteByKey(ctx context.Context, settingKey string) error {
	if settingKey == "" {
		return errors.New("setting key is empty")
	}

	_, err := store.db.Query().
		Table(store.settingTableName).
		Where(COLUMN_SETTING_KEY+" = ?", settingKey).
		Delete()

	return err
}

// SettingFindByID finds a setting by id
func (store *storeImplementation) SettingFindByID(ctx context.Context, settingID string) (SettingInterface, error) {
	if settingID == "" {
		return nil, errors.New("setting store > find by id: setting id is required")
	}

	query := SettingQuery().
		SetID(settingID).
		SetLimit(1)

	list, err := store.SettingList(ctx, query)

	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return list[0], nil
	}

	return nil, nil
}

// SettingFindByKey finds a setting by key
func (store *storeImplementation) SettingFindByKey(ctx context.Context, settingKey string) (SettingInterface, error) {
	if settingKey == "" {
		return nil, errors.New("setting store > find by key: setting key is required")
	}

	query := SettingQuery().
		SetKey(settingKey).
		SetLimit(1)

	list, err := store.SettingList(ctx, query)

	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return list[0], nil
	}

	return nil, nil
}

// SettingList retrieves a list of settings
func (store *storeImplementation) SettingList(ctx context.Context, query SettingQueryInterface) ([]SettingInterface, error) {
	if query == nil {
		return []SettingInterface{}, errors.New("at setting list > setting query is nil")
	}

	if err := query.Validate(); err != nil {
		return []SettingInterface{}, err
	}

	q := store.buildQuery(query)

	type settingRow struct {
		ID            string    `db:"id"`
		Key           string    `db:"setting_key"`
		Value         string    `db:"setting_value"`
		CreatedAt     time.Time `db:"created_at"`
		UpdatedAt     time.Time `db:"updated_at"`
		SoftDeletedAt time.Time `db:"soft_deleted_at"`
	}

	var rows []settingRow
	if err := q.Table(store.settingTableName).Get(&rows); err != nil {
		return []SettingInterface{}, err
	}

	list := make([]SettingInterface, 0, len(rows))
	for _, r := range rows {
		s := &settingImplementation{}
		s.SetID(r.ID)
		s.SetKey(r.Key)
		s.SetValue(r.Value)
		s.CreatedAtField.CreatedAt = r.CreatedAt
		s.UpdatedAtField.UpdatedAt = r.UpdatedAt
		s.SoftDeletesMaxDate.SoftDeletedAt = r.SoftDeletedAt
		list = append(list, s)
	}

	return list, nil
}

// SettingSoftDelete soft deletes a setting
func (store *storeImplementation) SettingSoftDelete(ctx context.Context, setting SettingInterface) error {
	if setting == nil {
		return errors.New("setting is nil")
	}

	setting.SetSoftDeletedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	return store.SettingUpdate(ctx, setting)
}

// SettingSoftDeleteByID soft deletes a setting by id
func (store *storeImplementation) SettingSoftDeleteByID(ctx context.Context, id string) error {
	setting, err := store.SettingFindByID(ctx, id)

	if err != nil {
		return err
	}

	return store.SettingSoftDelete(ctx, setting)
}

// SettingUpdate updates a setting
func (store *storeImplementation) SettingUpdate(ctx context.Context, setting SettingInterface) error {
	if setting == nil {
		return errors.New("settingstore > setting update. setting cannot be nil")
	}

	setting.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	row := map[string]any{
		COLUMN_SETTING_VALUE:   setting.GetValue(),
		COLUMN_UPDATED_AT:      setting.GetUpdatedAtCarbon().StdTime(),
		COLUMN_SOFT_DELETED_AT: setting.GetSoftDeletedAtCarbon().StdTime(),
	}

	_, err := store.db.Query().
		Table(store.settingTableName).
		Where(COLUMN_ID+" = ?", setting.GetID()).
		Update(row)

	return err
}

// == QUERY BUILDER ============================================================

// buildQuery builds a neat query from the setting query interface.
func (store *storeImplementation) buildQuery(query SettingQueryInterface) contractsorm.Query {
	// Use Model() to enable neat's automatic soft delete handling via SoftDeletesMaxDate
	q := store.db.Query().Model(&settingImplementation{})

	if query == nil {
		return q
	}

	if query.HasID() && query.ID() != "" {
		q = q.Where(COLUMN_ID+" = ?", query.ID())
	}

	if query.HasIDIn() && len(query.IDIn()) > 0 {
		args := make([]any, len(query.IDIn()))
		for i, id := range query.IDIn() {
			args[i] = id
		}
		q = q.WhereIn(COLUMN_ID, args)
	}

	if query.HasKey() && query.Key() != "" {
		q = q.Where(COLUMN_SETTING_KEY+" = ?", query.Key())
	}

	if query.HasLimit() && query.Limit() > 0 {
		q = q.Limit(query.Limit())
	}

	if query.HasOffset() && query.Offset() > 0 {
		q = q.Offset(query.Offset())
	}

	if query.HasOrderBy() && query.OrderBy() != "" {
		sortOrder := "desc"
		if query.HasSortOrder() && query.SortOrder() != "" {
			sortOrder = query.SortOrder()
		}
		q = q.OrderBy(query.OrderBy(), sortOrder)
	}

	// Handle soft delete filtering via neat's automatic handling (SoftDeletesMaxDate)
	if query.HasSoftDeletedIncluded() && query.SoftDeletedIncluded() {
		q = q.WithSoftDeleted()
	}

	return q
}

// NewStoreOptions define the options for creating a new setting store
type NewStoreOptions struct {
	SettingTableName   string
	DB                 *sql.DB
	AutomigrateEnabled bool
	DebugEnabled       bool
}

// NewStore creates a new setting store
func NewStore(opts NewStoreOptions) (StoreInterface, error) {
	if opts.DB == nil {
		return nil, errors.New("setting store: DB is required")
	}

	if opts.SettingTableName == "" {
		return nil, errors.New("setting store: settingTableName is required")
	}

	neatDB, err := neat.NewFromSQLDB(opts.DB)
	if err != nil {
		return nil, err
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	store := &storeImplementation{
		settingTableName:   opts.SettingTableName,
		db:                 neatDB,
		automigrateEnabled: opts.AutomigrateEnabled,
		debugEnabled:       opts.DebugEnabled,
		logger:             logger,
	}

	if store.automigrateEnabled {
		if err := store.MigrateUp(context.Background()); err != nil {
			return nil, err
		}
	}

	return store, nil
}
