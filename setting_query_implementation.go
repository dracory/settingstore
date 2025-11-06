package settingstore

import "errors"

// SettingQuery is a shortcut version of NewSettingQuery to create a new query
func SettingQuery() SettingQueryInterface {
	return NewSettingQuery()
}

// NewSettingQuery creates a new setting query
func NewSettingQuery() SettingQueryInterface {
	return &settingQueryImplementation{
		properties: make(map[string]interface{}),
	}
}

var _ SettingQueryInterface = (*settingQueryImplementation)(nil)

type settingQueryImplementation struct {
	properties map[string]interface{}
}

func (q *settingQueryImplementation) Validate() error {
	if q.HasCreatedAtGte() && q.CreatedAtGte() == "" {
		return errors.New("setting query. created_at_gte cannot be empty")
	}

	if q.HasCreatedAtLte() && q.CreatedAtLte() == "" {
		return errors.New("setting query. created_at_lte cannot be empty")
	}

	if q.HasID() && q.ID() == "" {
		return errors.New("setting query. id cannot be empty")
	}

	if q.HasIDIn() && len(q.IDIn()) < 1 {
		return errors.New("setting query. id_in cannot be empty array")
	}

	if q.HasKey() && q.Key() == "" {
		return errors.New("setting query. key cannot be empty")
	}

	if q.HasLimit() && q.Limit() < 0 {
		return errors.New("setting query. limit cannot be negative")
	}

	if q.HasOffset() && q.Offset() < 0 {
		return errors.New("setting query. offset cannot be negative")
	}

	return nil
}

func (q *settingQueryImplementation) Columns() []string {
	if !q.hasProperty("columns") {
		return []string{}
	}

	return q.properties["columns"].([]string)
}

func (q *settingQueryImplementation) SetColumns(columns []string) SettingQueryInterface {
	q.properties["columns"] = columns
	return q
}

func (q *settingQueryImplementation) HasCountOnly() bool {
	return q.hasProperty("count_only")
}

func (q *settingQueryImplementation) IsCountOnly() bool {
	return q.hasProperty("count_only") && q.properties["count_only"].(bool)
}

func (q *settingQueryImplementation) SetCountOnly(countOnly bool) SettingQueryInterface {
	q.properties["count_only"] = countOnly
	return q
}

func (q *settingQueryImplementation) HasCreatedAtGte() bool {
	return q.hasProperty("created_at_gte")
}

func (q *settingQueryImplementation) CreatedAtGte() string {
	return q.properties["created_at_gte"].(string)
}

func (q *settingQueryImplementation) SetCreatedAtGte(createdAtGte string) SettingQueryInterface {
	q.properties["created_at_gte"] = createdAtGte
	return q
}

func (q *settingQueryImplementation) HasCreatedAtLte() bool {
	return q.hasProperty("created_at_lte")
}

func (q *settingQueryImplementation) CreatedAtLte() string {
	return q.properties["created_at_lte"].(string)
}

func (q *settingQueryImplementation) SetCreatedAtLte(createdAtLte string) SettingQueryInterface {
	q.properties["created_at_lte"] = createdAtLte
	return q
}

func (q *settingQueryImplementation) HasID() bool {
	return q.hasProperty("id")
}

func (q *settingQueryImplementation) ID() string {
	return q.properties["id"].(string)
}

func (q *settingQueryImplementation) SetID(id string) SettingQueryInterface {
	q.properties["id"] = id
	return q
}

func (q *settingQueryImplementation) HasIDIn() bool {
	return q.hasProperty("id_in")
}

func (q *settingQueryImplementation) IDIn() []string {
	return q.properties["id_in"].([]string)
}

func (q *settingQueryImplementation) SetIDIn(idIn []string) SettingQueryInterface {
	q.properties["id_in"] = idIn
	return q
}

func (q *settingQueryImplementation) HasKey() bool {
	return q.hasProperty("key")
}

func (q *settingQueryImplementation) Key() string {
	return q.properties["key"].(string)
}

func (q *settingQueryImplementation) SetKey(key string) SettingQueryInterface {
	q.properties["key"] = key
	return q
}

func (q *settingQueryImplementation) HasLimit() bool {
	return q.hasProperty("limit")
}

func (q *settingQueryImplementation) Limit() int {
	return q.properties["limit"].(int)
}

func (q *settingQueryImplementation) SetLimit(limit int) SettingQueryInterface {
	q.properties["limit"] = limit
	return q
}

func (q *settingQueryImplementation) HasOffset() bool {
	return q.hasProperty("offset")
}

func (q *settingQueryImplementation) Offset() int {
	return q.properties["offset"].(int)
}

func (q *settingQueryImplementation) SetOffset(offset int) SettingQueryInterface {
	q.properties["offset"] = offset
	return q
}

func (q *settingQueryImplementation) HasOrderBy() bool {
	return q.hasProperty("order_by")
}

func (q *settingQueryImplementation) OrderBy() string {
	return q.properties["order_by"].(string)
}

func (q *settingQueryImplementation) SetOrderBy(orderBy string) SettingQueryInterface {
	q.properties["order_by"] = orderBy
	return q
}

func (q *settingQueryImplementation) HasSoftDeletedIncluded() bool {
	return q.hasProperty("soft_deleted_included")
}

func (q *settingQueryImplementation) SoftDeletedIncluded() bool {
	if !q.HasSoftDeletedIncluded() {
		return false
	}

	return q.properties["soft_deleted_included"].(bool)
}

func (q *settingQueryImplementation) SetSoftDeletedIncluded(softDeletedIncluded bool) SettingQueryInterface {
	q.properties["soft_deleted_included"] = softDeletedIncluded
	return q
}

func (q *settingQueryImplementation) HasSortOrder() bool {
	return q.hasProperty("sort_order")
}

func (q *settingQueryImplementation) SortOrder() string {
	return q.properties["sort_order"].(string)
}

func (q *settingQueryImplementation) SetSortOrder(sortOrder string) SettingQueryInterface {
	q.properties["sort_order"] = sortOrder
	return q
}

func (q *settingQueryImplementation) hasProperty(key string) bool {
	_, ok := q.properties[key]
	return ok
}
