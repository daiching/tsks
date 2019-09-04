package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/go-yaml/yaml"
)

const (
	favoriteErrorHeader    = "FavoriteError : "
	notExistFavoriteOfName = "That favorite of name is not exists."
	notExistAnyFavorites   = "There aren't any favorites."
)

func getFavoriteError(eBody string) error {
	return errors.New(favoriteErrorHeader + eBody)
}

type Favorite struct {
	Name    string
	Content string
}

type Favorites []Favorite

func (fs *Favorites) readFavorites() error {
	bytes, err := ioutil.ReadFile(config.FavoritesPath)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(bytes, fs)
	if err != nil {
		return err
	}
	return nil
}

func (fs *Favorites) saveFavorites() error {
	bytes, err := yaml.Marshal(*fs)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(config.FavoritesPath, bytes, os.ModeExclusive)
	if err != nil {
		return err
	}
	return nil
}

func addOrModFavorite(name string, content string) error {
	fs := Favorites{}
	err := fs.readFavorites()
	if err != nil {
		return err
	}

	var mi int = -1
	for i, f := range fs {
		if f.Name == name {
			mi = i
		}
	}
	if mi == -1 {
		nf := Favorite{
			name,
			content,
		}
		fs = append(fs, nf)
	} else {
		fs[mi].Content = content
	}

	err = fs.saveFavorites()
	if err != nil {
		return err
	}
	return nil
}

func deleteFavorite(name string) error {
	fs := Favorites{}
	err := fs.readFavorites()
	if err != nil {
		return err
	}

	var di int = -1
	for i, f := range fs {
		if f.Name == name {
			di = i
			break
		}
	}
	if di == -1 {
		return errors.New(favoriteErrorHeader + notExistFavoriteOfName)
	}
	fs = append(fs[:di], fs[di+1:]...)

	err = fs.saveFavorites()
	if err != nil {
		return err
	}
	return nil
}

func getFavoriteByName(name string) (*Favorite, error) {
	fs := Favorites{}
	err := fs.readFavorites()
	if err != nil {
		return nil, err
	}
	for _, f := range fs {
		if f.Name == name {
			return &f, nil
		}
	}
	return nil, getFavoriteError(notExistFavoriteOfName)
}

func writeFavorites() error {
	fs := Favorites{}
	err := fs.readFavorites()
	if err != nil {
		return err
	}
	if len(fs) == 0 {
		return errors.New(favoriteErrorHeader + notExistAnyFavorites)
	}
	for _, f := range fs {
		fmt.Println(" -", f.Name, ":", f.Content)
	}
	return nil
}

func writeFavoriteByName(name string) error {
	f, err := getFavoriteByName(name)
	if err != nil {
		return err
	}
	fmt.Println(" -", f.Name, ":", f.Content)
	return nil
}
