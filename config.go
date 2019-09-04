package main

import (
	"errors"
)

const (
	configErrorHeader = "ConfigError : "
)

type Config struct {
	TaskListPath   string
	FavoritesPath  string
	TaskListPath1  string
	TaskListPath2  string
	FavoritesPath1 string
	FavoritesPath2 string
}

var config Config = Config{
	"",
	"",
	"$HOME/tsks/taskList.yaml",
	"$GOPATH/src/github.com/daiching/tsks/conf/taskList.yaml",
	"$HOME/tsks/favorites.yaml",
	"$GOPATH/src/github.com/daiching/tsks/conf/favorites.yaml",
}

func getConfigError(eBody string) error {
	return errors.New(configErrorHeader + eBody)
}

func readConfig() error {
	if exists(config.TaskListPath1) {
		config.TaskListPath = transEnvPath(config.TaskListPath1)
	} else if exists(config.TaskListPath2) {
		config.TaskListPath = transEnvPath(config.TaskListPath2)
	} else {
		return getUtilError("taskList.yaml is not found.")
	}

	if exists(config.FavoritesPath1) {
		config.FavoritesPath = transEnvPath(config.FavoritesPath1)
	} else if exists(config.FavoritesPath2) {
		config.FavoritesPath = transEnvPath(config.FavoritesPath2)
	} else {
		return getUtilError("favorites.yaml is not found.")
	}

	return nil
}
