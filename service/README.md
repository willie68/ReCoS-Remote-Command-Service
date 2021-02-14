# Remote Desk

Desktop Remote software, for executing scripts remotely but secure.

# Profileconfiguration

Every profile has it's own configuration file. This file is written in yaml and has the following sections:

In the root you will find the following parameters

name: The name of the profile
description: a user readable description
pages: This are the different pages for showing up the actions
actions: a list of different actions. An action can appear on different pages. But as it's the same action, the status/result will be the same.

This is an example

```yaml
name: Default
description: description for default
pages:
    - name: page1
      columns: 5
      rows: 5
      cells:
        - action1
        - action2
        - action3
    - name: page2
      columns: 3
      rows: 3
      cells:
        - action1
        - action2
actions:
    - type: SINGLE
      name: action1
      title: Action Title Sync
      description: description for action
      runone: true
      commands:
        - type: DELAY
          name: delay
          parameters:
            time: 2
        - type: EXECUTE
          name: execute
          parameters:
            command: echo.bat 
            args:
              - "Hello world"
        - type: DELAY
          name: delay
          parameters:
            time: 1
    - type: SINGLE
      name: action2
      title: Action Title Async
      description: description for action
      runone: false
      commands:
        - type: DELAY
          name: delay
          parameters:
            time: 2
        - type: EXECUTE
          name: execute
          parameters:
            command: echo.bat 
            args: 
              - "Hello world"
        - type: DELAY
          name: delay
          parameters:
            time: 1
    - type: SINGLE
      name: action3
      title: Execute go version
      description: description for action
      runone: true
      commands:
        - type: EXECUTE
          name: execute
          parameters:
            command: go.exe 
            args:
              - "version"
```



## Page

A page is a view component mainly showing with rows and columns. Each of this cells will than visualize an action. The cells list will link to an action by name in the action list. The index of an action of a cell is calculated as 

`index = (cell.row * page.rows) + cell.column`

parameters:

name: The name of the page
columns: Number of columns on the page
rows: Number of rows on the page
cells: List of names of the action per cell

Example:

```yaml
name: page1
columns: 5
rows: 5
cells:
  - action1
  - action2
  - action3
```



## Action

An action is the part which defines, what to do if a cell is triggered. 

The following parameters are used:

type: SINLGE is a single shot action. The action is always starting the command list. TOGGLE is an action with two states, just like a on/off switch. For each transition you can define an own command list. MULTISTAGE is the third option. Here you can define 3 or more stages, and you every stage you can define the status and a command list, which is fired on activating this stage. As you can see, TOGGLE is a Multiswitch with 2 Stages.
name: s the name of the action
title: the title of the action used by the UI
description: user defined description of this action
runone: is true or false. On true, if the action is fired twice, all commands of the first execution must be finished before the second execution will take place. On false, the execution will start directly without checking the action execution state.
icon: the icon which will be displayed on the cell
commands: list of commands for execution of this action

```yaml
type: SINGLE
name: action1
title: Action Title Sync
description: description for action
runone: true
icon: trash_can.png
commands:
  - type: DELAY
    name: delay
    parameters:
      time: 2
  - type: EXECUTE
    name: execute
    parameters:
      command: echo.bat 
      args:
        - "Hello world"
  - type: DELAY
    name: delay
    parameters:
      time: 1
```



### Command

This is the command, which should be executed

type: the type of the command
name: names the command
icon: should be the icon that should be displayed when running this command
title: should be the text that should be displayed when running this command
parameters: parameters defers from command to command

#### No Operation

Do nothing.

Type: NOOP

Parameter:  none

Example

```yaml
type: DELAY
name: delay
icon: accesibility.png 
title: Do Nothing
```

#### Delay

Type: DELAY

Parameter: 

time: time to delay in Seconds

Example

```yaml
type: DELAY
name: delay
parameters:
  time: 2
```

#### Timer

Starting a timer with a response every second. You can define the format of the timer message and the message on finish.

Type: TIMER

Parameter: 

time: time to delay in Seconds
format: the message for the response, defaults %d seconds
finished: the message at the end of the timer, defaults: finished

Example

```yaml
type: TIMER
name: timer
parameters:
  time: 10
  format: noch %ds
  finished: Fertig
```

#### Execute

Type: EXECUTE

Parameter:

command: the executable or shell script to execute, with or without path
args: list of string arguments to this executable

Example

```yaml
type: EXECUTE
name: execute
parameters:
  command: go.exe 
  args:
    - "version"
```

#### Page

Switch to another page.

type: PAGE

Parameter:
page: the name of the page to switch to

```yaml
type: PAGE
name: page
parameters:
  page: page2
```

#### Keys

Sending keys to the active application. This command is emulating a keyboard input by sending key strokes of a keyboard to the active application. You can use different keyboard layouts and there are some macros defining special keys.

type: KEYS

layout: defining the layout of the keyboard used to send the data. en for English (us) "qwerty" and de for a German "qwertz" keyboard layout. Default is "de"

keys: are the string with the keys used to send. For special keys there are defined special macros. Every macro starts with an "{" and ends than with "}". If you want to send the "{" as a character simply double this. ("{" -> "{{"). 

Another specialized character is the "~" char. It will lead into a 1 second delay between the typing. To get the "~" Character, simple double it.

The following macros are defined: 

| Macro               | Keyboard key |
| ------------------- | ------------ |
| backspace, bs, bksp | backspace    |
| break               | break        |
| capslock            | caps lock    |
| del, delete         | delete       |
| down                | arrow down   |
| end                 | end          |
| enter               | enter        |
| esc                 | esc          |
| help                | help         |
| home                | home         |
| ins, insert         | insert       |
| left                | arrow left   |
| num                 | num lock     |
| pgdn                | page down    |
| pgup                | page up      |
| prtsc               | print screen |
| right               | arrow right  |
| scrolllock          | scroll lock  |
| tab                 | tab          |
| up                  | arrow up     |
| f1 .. f12           | function key 1 ... 12 |


```yaml
type: KEYS
name: sendkeys
parameters:
  layout: de
  keys: "akteon00{enter}"
```

#### Controlling Application Main Window

With this command, you can control the main window of an application.

type: WINDOWCTRL

Parameter:
caption: the caption of the application window
command: the command to execute on this window. Possible values are:
minimize: for minimizing the application window
activate: for activating the application window again. (restore it if minimized and active/bring it to front) 
move  x y: moving the window to the new position x,y

```yaml
# activate the german calculator program
- type: WINDOWCTRL
  name: control window
  parameters:
    caption: Rechner 
    command: activate
# move it to it's new location
- type: WINDOWCTRL
  name: control window
  parameters:
    caption: Rechner 
    command: move 700 300 
# minimize it
- type: WINDOWCTRL
  name: control window
  parameters:
    caption: Rechner 
    command: minimize 
```

​	