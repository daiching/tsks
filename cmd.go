package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

const (
	cmdErrorHeader           = "CmdError : "
	subCmdNotSelectedError   = cmdErrorHeader + "Specified sub cmd."
	wrongSubCmdError         = cmdErrorHeader + "Selected subcommand is wrong."
	notEnoughSubCmdArgsError = cmdErrorHeader + "The number of args in sub command is not enough."

	commentAdd = "Add new task. (Syntax : tsks add [option] <task content>)"
	commentLs  = "Display task list. (Syntax : tsks ls [option])"
	commentFin = "Finish any task. (Syntax : tsks fin [option] <the number of task>)"
	commentRev = "Revival any task. (Syntax : tsks rev [option] <the number of task>)"
	commentFav = "Register a task for favorite one. (Syntax : tsks fav [option] <the number of task>)"
	commentFls = "Delete tasks before a day from taskList. (Syntax : tsks fls <a day>)"

	commentAddName     = "Specify a name of favorite task."
	commentAddDay      = "Date option using t(today) and YYYY-MM-DD. ex: t => today, t-1 => yesterday, 2019-01-01"
	commentLsdayRange  = "Range of date option using t(today), w(wild card) and YYYY-MM-DD. ex: w => all day, t => today, t-1:W => all day after yesterday, 2019-01-01:2019-02-28"
	commentLsAll       = "Display tasks including finished task."
	commentLsYaml      = "Display taskList.yaml."
	commentFinOrRevDay = "Date option using t(today) and YYYY-MM-DD. ex: t => today, t-1 => yesterday, 2019-01-01"
	commentFavIsDelete = "Delete option of favorite tasks."
)

var debugNow bool = false

// pairs of sub command function and description.
type Subcmd struct {
	Cmd         func([]string) error
	Description string
}

// map of "Subcmd" struct. This is used in determing sub command from arguments
// and display help of sub command when any sub command is not selected.
var subcmds = map[string]Subcmd{
	"add": {cmdAdd, commentAdd},
	"ls":  {cmdLs, commentLs},
	"fin": {cmdFinOrRev, commentFin},
	"rev": {cmdFinOrRev, commentRev},
	"fav": {cmdFav, commentFav},
	"fls": {cmdFls, commentFls},
}

// command entry point.
func cmdMain() error {
	dn := flag.Bool("d", false, "Debug mode on/off.")
	flag.Parse()
	debugNow = *dn
	mainArgs := flag.Args()
	if len(mainArgs) < 1 {
		dspSubCmdList()
		return nil
	}
	var err error = nil
	if s, exist := subcmds[mainArgs[0]]; exist {
		err = s.Cmd(mainArgs)
	} else {
		err = errors.New(wrongSubCmdError)
	}
	return err
}

// add new task to list.
func cmdAdd(mainArgs []string) error {
	add := flag.NewFlagSet("add", flag.ExitOnError)
	name := add.Bool("n", false, commentAddName)
	day := add.String("d", "t", commentAddDay)
	add.Parse(mainArgs[1:])
	args := add.Args()
	if len(args) < 1 {
		return errors.New(notEnoughSubCmdArgsError)
	}
	newTask := Task{}
	var c *Favorite
	var err error
	if *name {
		c, err = getFavoriteByName(args[0])
		if err != nil {
			return err
		}
	}
	if c != nil {
		newTask.Name = c.Name
		newTask.Content = c.Content
		newTask.IsFin = false
	} else {
		newTask.Name = ""
		newTask.Content = strings.Join(args, " ")
		newTask.IsFin = false
	}
	return addTask(&newTask, *day)
}

// display task list.
func cmdLs(mainArgs []string) error {
	ls := flag.NewFlagSet("ls", flag.ExitOnError)
	day := ls.String("d", "t", commentLsdayRange)
	isIncludeFin := ls.Bool("a", false, commentLsAll)
	isTextOfYaml := ls.Bool("t", false, commentLsYaml)
	ls.Parse(mainArgs[1:])
	if *isTextOfYaml {
		txt, err := readTaskListYaml()
		if err != nil {
			return err
		}
		fmt.Println(txt)
		return nil
	}
	// only particular day (default is today.)
	if err := lsTasks(*day, *isIncludeFin); err != nil {
		return err
	}
	return nil
}

// finish a task in a specified day.
func cmdFinOrRev(mainArgs []string) error {
	var finOrRev *flag.FlagSet
	if mainArgs[0] == "fin" {
		finOrRev = flag.NewFlagSet("fin", flag.ExitOnError)
	} else {
		finOrRev = flag.NewFlagSet("rev", flag.ExitOnError)
	}
	day := finOrRev.String("d", "t", commentFinOrRevDay)
	finOrRev.Parse(mainArgs[1:])
	subArgs := finOrRev.Args()
	if len(subArgs) < 1 {
		return errors.New(notEnoughSubCmdArgsError)
	}
	nums := subArgs[0]
	var err error
	if mainArgs[0] == "fin" {
		err = finTask(*day, nums)
	} else {
		err = revTask(*day, nums)
	}
	if err != nil {
		return err
	}
	return nil
}

func cmdFav(mainArgs []string) error {
	fav := flag.NewFlagSet("fav", flag.ExitOnError)
	isDelete := fav.Bool("d", false, commentFavIsDelete)
	fav.Parse(mainArgs[1:])
	subArgs := fav.Args()

	if len(subArgs) == 0 {
		if err := writeFavorites(); err != nil {
			return err
		}
		return nil
	}

	if len(subArgs) == 1 {
		if *isDelete {
			if err := deleteFavorite(subArgs[0]); err != nil {
				return err
			}
			return nil
		}
		if err := writeFavoriteByName(subArgs[0]); err != nil {
			return err
		}
		return nil
	}

	if err := addOrModFavorite(subArgs[0], strings.Join(subArgs[1:], " ")); err != nil {
		return err
	}
	return nil
}

func cmdFls(mainArgs []string) error {
	fav := flag.NewFlagSet("fls", flag.ExitOnError)
	fav.Parse(mainArgs[1:])
	subArgs := fav.Args()
	if len(subArgs) < 1 {
		return errors.New(notEnoughSubCmdArgsError)
	}
	if err := flsTsks(subArgs[0]); err != nil {
		return err
	}
	return nil
}

func dspSubCmdList() {
	fmt.Println("Select next sub commands.")
	for k, s := range subcmds {
		fmt.Print("  ", k, "\n")
		fmt.Print("    ", s.Description, "\n")
	}
	fmt.Println("\nIf you want to see sub command help, input next syntax : tsks <sub command> -h")
}

func readTaskListYaml() (string, error) {
	f, err := os.Open(config.TaskListPath)
	if err != nil {
		return "", err
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	return string(b), err
}
