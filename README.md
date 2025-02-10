# Settings Store

[![Tests Status](https://github.com/gouniverse/settingstore/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/gouniverse/settingstore/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/gouniverse/settingstore)](https://goreportcard.com/report/github.com/gouniverse/settingstore)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/gouniverse/settingstore)](https://pkg.go.dev/github.com/gouniverse/settingstore)

Saves settings in an SQL database. 

## 🌏  Open in the Cloud 
Click any of the buttons below to start a new development environment to demo or contribute to the codebase without having to install anything on your machine:

[![Open in VS Code](https://img.shields.io/badge/Open%20in-VS%20Code-blue?logo=visualstudiocode)](https://vscode.dev/github/gouniverse/settingstore)
[![Open in Glitch](https://img.shields.io/badge/Open%20in-Glitch-blue?logo=glitch)](https://glitch.com/edit/#!/import/github/gouniverse/settingstore)
[![Open in GitHub Codespaces](https://github.com/codespaces/badge.svg)](https://codespaces.new/gouniverse/settingstore)
[![Edit in Codesandbox](https://codesandbox.io/static/img/play-codesandbox.svg)](https://codesandbox.io/s/github/gouniverse/settingstore)
[![Open in StackBlitz](https://developer.stackblitz.com/img/open_in_stackblitz.svg)](https://stackblitz.com/github/gouniverse/settingstore)
[![Open in Repl.it](https://replit.com/badge/github/withastro/astro)](https://replit.com/github/gouniverse/settingstore)
[![Open in Codeanywhere](https://codeanywhere.com/img/open-in-codeanywhere-btn.svg)](https://app.codeanywhere.com/#https://github.com/gouniverse/settingstore)
[![Open in Gitpod](https://gitpod.io/button/open-in-gitpod.svg)](https://gitpod.io/#https://github.com/gouniverse/settingstore)


## Description

Every application needs to preserve settings between multiple restarts. This package helps save the setting represented as key-value pairs in an SQL database.

## Features

- Saves settings data as key-value pairs
- Supports SQLite, MySQL and Postgres
- Uses sql.DB directly
- Automigration

## Installation
```
go get -u github.com/gouniverse/settingstore
```

## Setup

```
// as one line
settingStore = settingstore.NewStore(settingstore.WithDb(databaseInstance), settingstore.WithTableName("settings"), entitystore.WithAutoMigrate(true))


// as multiple lines
settingStore = settingstore.NewStore(settingstore.WithDb(databaseInstance), settingstore.WithTableName("settings"))
settingStore.AutoMigrate()

```

## Usage

1. Create a new key value setting pair
```
settingsStore.Set("app.name", "My Web App")
settingsStore.Set("app.url", "http://localhost")
settingsStore.Set("server.ip", "127.0.0.1")
settingsStore.Set("server.port", "80")
```

2. Retrieve a setting value (or default value if not exists)
```
appName = settingsStore.Get("app.name", "Default Name")
appUrl = settingsStore.Get("app.url", "")
serverIp = settingsStore.Get("server.ip", "")
serverPort = settingsStore.Get("server.port", "")
```

3. Check if required setting is setup
```
if serverIp == "" {
    log.Panic("server ip not setup")
}
```

## Methods

These methods may be subject to change as still in development

### Store Methods

- NewStore(opts ...StoreOption) *Store - creates a new setting store
- WithAutoMigrate(automigrateEnabled bool) StoreOption - a store option. sets whether database migration will run on setup
- WithDb(db *sql.DB) StoreOption - a store option (required). sets the DB to be used by the store
- WithTableName(settingsTableName string) StoreOption - a store option (required). sets the table name for the setting store in the DB

- AutoMigrate() error - auto migrate (create the tables in the database) the settings store tables
- DriverName(db *sql.DB) string - the name of the driver used for SQL strings (you may use this if you need to debug)
- SqlCreateTable() string - SQL string for creating the tables (you may use this string if you want to set your own migrations)

### Setting Methods

- Delete() bool - deletes the entity
- FindByKey(key string) *Setting - finds a Setting by key
- Get(key string, valueDefault string) string - gets a value from key-value setting pair
- GetJSON(key string, valueDefault interface{}) interface{} - gets a value as JSON from key-value setting pair
- Keys() ([]string, error) - gets all keys sorted alphabetically (useful if you want to list these in admin panel)
- Remove(key string) error - removes a setting from store
- Set(key string, value string) (bool, error) - sets new key value pair
- SetJSON(key string, value interface{}) (bool, error) - sets new key JSON value pair
