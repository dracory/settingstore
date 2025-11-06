package settingstore

import (
    "reflect"
    "strings"
    "testing"
)

func TestNewSettingQuery_Defaults(t *testing.T) {
    query := NewSettingQuery()

    if query == nil {
        t.Fatal("expected query to be created")
    }

    if err := query.Validate(); err != nil {
        t.Fatalf("unexpected validation error: %v", err)
    }

    if cols := query.Columns(); len(cols) != 0 {
        t.Fatalf("expected no columns, got %v", cols)
    }

    if query.HasLimit() {
        t.Fatal("unexpected limit flag")
    }

    if query.HasOffset() {
        t.Fatal("unexpected offset flag")
    }

    if query.HasKey() {
        t.Fatal("unexpected key flag")
    }

    if query.IsCountOnly() {
        t.Fatal("expected count-only to default to false")
    }

    if query.SoftDeletedIncluded() {
        t.Fatal("expected soft deleted to be excluded by default")
    }
}

func TestSettingQuery_SettersAndGetters(t *testing.T) {
    query := NewSettingQuery()

    columns := []string{"id", "setting_key"}
    query.SetColumns(columns)

    if got := query.Columns(); !reflect.DeepEqual(columns, got) {
        t.Fatalf("unexpected columns: %v", got)
    }

    query.SetCountOnly(true)
    if !query.HasCountOnly() {
        t.Fatal("expected count-only flag to be set")
    }
    if !query.IsCountOnly() {
        t.Fatal("expected IsCountOnly to be true")
    }

    query.SetCreatedAtGte("2024-01-01 00:00:00")
    if !query.HasCreatedAtGte() {
        t.Fatal("expected created_at_gte flag to be set")
    }
    if query.CreatedAtGte() != "2024-01-01 00:00:00" {
        t.Fatalf("unexpected created_at_gte: %s", query.CreatedAtGte())
    }

    query.SetCreatedAtLte("2024-02-01 00:00:00")
    if !query.HasCreatedAtLte() {
        t.Fatal("expected created_at_lte flag to be set")
    }
    if query.CreatedAtLte() != "2024-02-01 00:00:00" {
        t.Fatalf("unexpected created_at_lte: %s", query.CreatedAtLte())
    }

    query.SetID("abc")
    if !query.HasID() {
        t.Fatal("expected id flag to be set")
    }
    if query.ID() != "abc" {
        t.Fatalf("unexpected id: %s", query.ID())
    }

    query.SetIDIn([]string{"a", "b"})
    if !query.HasIDIn() {
        t.Fatal("expected id_in flag to be set")
    }
    if got := query.IDIn(); !reflect.DeepEqual([]string{"a", "b"}, got) {
        t.Fatalf("unexpected id_in: %v", got)
    }

    query.SetKey("feature")
    if !query.HasKey() {
        t.Fatal("expected key flag to be set")
    }
    if query.Key() != "feature" {
        t.Fatalf("unexpected key: %s", query.Key())
    }

    query.SetLimit(10)
    if !query.HasLimit() {
        t.Fatal("expected limit flag to be set")
    }
    if query.Limit() != 10 {
        t.Fatalf("unexpected limit: %d", query.Limit())
    }

    query.SetOffset(5)
    if !query.HasOffset() {
        t.Fatal("expected offset flag to be set")
    }
    if query.Offset() != 5 {
        t.Fatalf("unexpected offset: %d", query.Offset())
    }

    query.SetOrderBy("created_at")
    if !query.HasOrderBy() {
        t.Fatal("expected order_by flag to be set")
    }
    if query.OrderBy() != "created_at" {
        t.Fatalf("unexpected order_by: %s", query.OrderBy())
    }

    query.SetSortOrder("ASC")
    if !query.HasSortOrder() {
        t.Fatal("expected sort_order flag to be set")
    }
    if query.SortOrder() != "ASC" {
        t.Fatalf("unexpected sort_order: %s", query.SortOrder())
    }

    query.SetSoftDeletedIncluded(true)
    if !query.HasSoftDeletedIncluded() {
        t.Fatal("expected soft_deleted_included flag to be set")
    }
    if !query.SoftDeletedIncluded() {
        t.Fatal("expected soft deleted to be included")
    }

    query.SetSoftDeletedIncluded(false)
    if query.SoftDeletedIncluded() {
        t.Fatal("expected soft deleted to be disabled when set to false")
    }
}

func TestSettingQueryValidateErrors(t *testing.T) {
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
