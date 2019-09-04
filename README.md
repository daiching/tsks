# tsks
tsks is simple cli task management tool.

## Commands
- tsks add : add new task.
- tsks ls  : preview task list.
- tsks fin : finish task.
- tsks rev : revival task.
- tsks fav : register task often used.
- tsks flv : remove task before a specified day
* sub commands' help is displayed by next syntax.
tsks <sub commands> -h

## How to use (examples)
- senario 1: add new task -> review task list -> finish the task -> review task list including deleted task
``` bash
$ tsks add a example task.
$ tsks ls
# 2019-09-04
 [wip] 1. a example task.
$ tsks fin 1
$ tsks ls -a
# 2019-09-04
 [fin] 1. a example task.
```

- senario 2: add new favorite task (often used) -> add new task specifying the favorite task -> review task list
``` bash
$ tsks fav favoriteTask this is sample favorite task.
$ tsks add -n favoriteTask
$ tsks ls
# 2019-09-04
 [wip] 1. this is sample favorite task. (favoriteTask)
```

- senario 3: add new task specifying a day -> add new task specifying tommorow -> review task list from yesterday to tomorrow
today is 2019/09/04
``` bash
$ tsks add -d 2019-09-03 I develop tsks.
$ tsks add -d t+1 I work hard.
$ tsks ls -d t-1:t+1
# 2019-09-05
 [wip] 1. I work hard.
# 2019-09-04
 [wip] 1. this is sample favorite task. (favoriteTask)
# 2019-09-03
 [wip] 1. I develop tsks.
 ```
 
- senario 4: remove all task before tomorrow -> review task list all days.
 ``` bash
 $ tsks ls -d t-1:t+1
# 2019-09-05
 [wip] 1. I work hardtsks ls -d t-1!
# 2019-09-04
 [wip] 1. this is sample favorite task. (favoriteTask)
# 2019-09-03
 [wip] 1. I develop tsks.
$ tsks ls -d w
(no line)
 ```
