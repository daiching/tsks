package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/go-yaml/yaml"
)

const (
	taskErrorHeader      = "TaskError : "
	dayTaskNotExistError = "Specified day of task is not exist."
	dayOptionError       = "String format of '-d' option is uncorrect."
	minDay               = "0001-01-01"
	maxDay               = "3000-12-31"
)

type Task struct {
	Name    string
	Content string
	IsFin   bool
}

// output line of a task in a day.
func (t *Task) writeTask(i int, isIncludeFin bool) {
	if !isIncludeFin && t.IsFin {
		return
	}
	var line string = strconv.Itoa(i) + ". " + t.Content
	// when isIncludeFin is true, finished tasks are also displayed with marked *.
	if isIncludeFin && t.IsFin {
		line = " [fin] " + line
	} else {
		line = " [wip] " + line
	}
	if t.Name == "" {
		fmt.Println(line)
	} else {
		fmt.Println(line + " (" + t.Name + ")")
	}
}

type Tasks map[string][]Task

func (ts *Tasks) readTasks() error {
	bytes, err := ioutil.ReadFile(config.TaskListPath)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(bytes, ts)
	if err != nil {
		return err
	}
	return nil
}

func (ts *Tasks) saveTasks() error {
	bytes, err := yaml.Marshal(*ts)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(config.TaskListPath, bytes, os.ModeExclusive)
	if err != nil {
		return err
	}
	return nil
}

type DayOfTask []string

func (d DayOfTask) Len() int {
	return len(d)
}

func (d DayOfTask) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

func (d DayOfTask) Less(i, j int) bool {
	return d[i] > d[j]
}

func getStartAndEndDay(dayOption string) (map[string]string, error) {
	var getDay = func(r string, isStart bool) (string, error) {
		// today or wild card
		switch r {
		case "t":
			return today(), nil
		case "w":
			if isStart {
				return minDay, nil
			}
			return maxDay, nil
		}

		// day format yyyy-mm-dd
		if checkDayFormat(r) {
			return r, nil
		}

		// incorrect format or t+?, t-?
		var tmp []string
		var isPlus bool
		if t := strings.Split(r, "+"); len(t) == 2 {
			tmp = t
			isPlus = true
		}
		if t := strings.Split(r, "-"); len(t) == 2 {
			tmp = t
			isPlus = false
		}
		if len(tmp) != 2 {
			return "", getTaskError(dayOptionError)
		}
		if tmp[0] != "t" {
			return "", getTaskError(dayOptionError)
		}
		i, err := strconv.Atoi(tmp[1])
		if err != nil {
			return "", getTaskError(dayOptionError)
		}
		if isPlus {
			return aday((1) * i), nil
		}
		return aday((-1) * i), nil
	}

	tmp := strings.Split(dayOption, ":")
	dayRange := make(map[string]string)
	l := len(tmp)
	for i := 0; i < l; i++ {
		if l == 1 {
			if dayOption == "w" {
				dayRange["start"] = ""
				dayRange["end"] = maxDay
				return dayRange, nil
			}
			d, err := getDay(dayOption, true)
			if err != nil {
				return nil, err
			}
			dayRange["day"] = d
			return dayRange, nil
		}
		if i == 0 {
			d, err := getDay(tmp[0], true)
			if err != nil {
				return nil, err
			}
			dayRange["start"] = d
		}
		d, err := getDay(tmp[1], false)
		if err != nil {
			return nil, err
		}
		dayRange["end"] = d
		return dayRange, nil
	}
	return nil, getTaskError(dayOptionError)
}

// output line of a day.
func writeDayLine(day string) {
	fmt.Printf("# %v \n", day)
}

func getTaskError(eBody string) error {
	return errors.New(taskErrorHeader + eBody)
}

// this function is called if sub command 'add' is used.
func addTask(newTask *Task, dayOption string) error {
	tasks := Tasks{}
	err := tasks.readTasks()
	if err != nil {
		return err
	}
	dayRange, err := getStartAndEndDay(dayOption)
	day, ok := dayRange["day"]
	if !ok {
		return getTaskError(dayOptionError)
	}
	if _, exist := tasks[dayRange["day"]]; exist {
		tasks[day] = append(tasks[day], *newTask)
	} else {
		t := []Task{*newTask}
		tasks[day] = t
	}
	err = tasks.saveTasks()
	if err != nil {
		return err
	}
	return nil
}

// this function is firstlly called if sub command 'fin' is used.
func finTask(day string, numbers string) error {
	return finOrRevivalTask(day, numbers, false)
}

// this function is firstlly called if sub command 'rev' is used.
func revTask(day string, numbers string) error {
	return finOrRevivalTask(day, numbers, true)
}

// this function is called if sub command 'fin' or 're' is used.
// if number is plus, a task is finished and if not, revival.
func finOrRevivalTask(dayOption string, numbers string, isRev bool) error {
	numList := strings.Split(numbers, ",")
	tasks := Tasks{}
	err := tasks.readTasks()
	if err != nil {
		return err
	}
	dayRange, err := getStartAndEndDay(dayOption)
	day, ok := dayRange["day"]
	if !ok {
		return getTaskError(dayOptionError)
	}
	if _, exist := tasks[day]; exist {
		if numList[0] == "w" {
			for i := 0; i < len(tasks[day]); i++ {
				if isRev {
					tasks[day][i].IsFin = false
				} else {
					tasks[day][i].IsFin = true
				}
			}
			err := tasks.saveTasks()
			if err != nil {
				return err
			}
			return nil
		}
		for _, strnum := range numList {
			num, err := strconv.Atoi(strnum)
			if err != nil {
				return err
			}
			if len(tasks[day]) < num-1 {
				return getTaskError("there aren't the number of task at the specified " + day + ".")
			}
			if isRev {
				tasks[day][num-1].IsFin = false
			} else {
				tasks[day][num-1].IsFin = true
			}
		}
		err := tasks.saveTasks()
		if err != nil {
			return err
		}
	} else {
		return getTaskError(dayTaskNotExistError)
	}
	return nil
}

// this function is firstlly called if sub command 'ls' is used.
func lsTasks(dayOption string, isIncludeFin bool) error {
	dayRange, err := getStartAndEndDay(dayOption)
	if err != nil {
		return err
	}
	tasks := Tasks{}
	err = tasks.readTasks()
	if err != nil {
		return err
	}
	switch len(dayRange) {
	case 1:
		if tasksInDay, ok := tasks[dayRange["day"]]; ok {
			writeDayLine(dayRange["day"])
			for i, task := range tasksInDay {
				task.writeTask(i+1, isIncludeFin)
			}
			return nil
		} else {
			if dayOption == "t" {
				return nil
			}
			return getTaskError("There aren't tasks at the specified " + dayRange["day"] + ".")
		}
	case 2:
		// days of tasks are sorted by day descending order with.
		var days DayOfTask
		for day, taskInDay := range tasks {
			if !(dayRange["start"] <= day && day <= dayRange["end"]) {
				continue
			}
			if isIncludeFin {
				days = append(days, day)
				continue
			}
			var isAllFin = true
			// check if all tasks is already finished.
			for _, task := range taskInDay {
				if !task.IsFin {
					isAllFin = false
					break
				}
			}
			if !isAllFin {
				days = append(days, day)
			}
		}
		sort.Sort(days)
		for _, day := range days {
			writeDayLine(day)
			for i, task := range tasks[day] {
				task.writeTask(i+1, isIncludeFin)
			}
		}
	}
	return nil
}

func flsTsks(day string) error {
	m, err := getStartAndEndDay(day)
	if err != nil {
		return err
	}
	oTasks := Tasks{}
	err = oTasks.readTasks()
	nTasks := oTasks
	if err != nil {
		return err
	}
	for d, _ := range oTasks {
		if len(m) == 1 {
			if d <= m["day"] {
				delete(nTasks, d)
			}
		} else {
			if d >= m["start"] && d <= m["end"] {
				delete(nTasks, d)
			}
		}
	}
	return nTasks.saveTasks()
}
