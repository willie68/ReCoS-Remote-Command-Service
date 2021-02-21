# Remote Command Service, ReCoS

Remote Command Service, for executing scripts remotely but secure. This project consists of 2 components. One is the service, a small go based microservice, which delivers all the funktionality behind the UI. This service is responsible for the execution of the different actions. The other project is the UI Web Component, builded with VUE as the main UI Framework. 

# Service

## Profileconfiguration

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



### Page

A page is a view component mainly showing with rows and columns. Each of this cells will than visualize an action. The cells list will link to an action by name in the action list. The index of an action of a cell is calculated as 

`index = (cell.row * page.rows) + cell.column`

parameters:

`name`: The name of the page
`columns`: Number of columns on the page
`rows`: Number of rows on the page
`toolbar`: `show`, this page will appear in the toolbar to directly switch to, `hide`, if this page should appear in the toolbar. With the page command you can switch to this page.
`cells`: List of names of the action per cell

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



### Action

An action is the part which defines, what to do if a cell is triggered. 

The following parameters are used:

`type`: **SINGLE** is a single shot action. The action is always starting the command list. 
**DISPLAY** is a display only cell. It will only show Text, Icons, or images, but you can't interact with it.
**MULTI** is the third option. Here you can define 3 or more stages, and you every stage you can define the status and a command list, which is fired on activating this stage. As you can see, a simple on/off switch is a Multi with 2 Stages.
`name`: s the name of the action
title: the title of the action used by the UI
`description`: user defined description of this action
`runone`: is true or false. On true, if the action is fired twice, all commands of the first execution must be finished before the second execution will take place. On false, the execution will start directly without checking the action execution state.
`icon`: the icon which will be displayed on the cell
`fontsize`: the size of the title and the text, defaults to 14
`fontcolor`: the color of the title and the text, defaults to black
`outlined`: true or false, sometime reading a black text on a black ground is a little bit difficult. lining out can help.
`commands`: list of commands for execution of this action
`actions`: only apply in a MULTI action. For every stage there should be the name of the action which will be called, when the stage is executed. If a stage is executed, the icon of the last executed action (stage) will be displayed as the icon of the multi action and the title will be displayed on the text line. Example for a multi action:

```yaml
type: MULTI
name: action2
title: Multi Action
description: description for action
icon: hand_share.png
actions:
 - action2_1
 - action2_2
 - action2_3
```

Example for a single action:

```yaml
type: SINGLE
name: action1
title: Action Title Sync
description: description for action
fontsize: 10
fontcolor: "#FF0000"
outlined: true
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



#### Command

This is the command, which should be executed

type: the type of the command
name: names the command
icon: should be the icon that should be displayed when running this command
title: should be the text that should be displayed when running this command
parameters: parameters defers from command to command

##### No Operation

Do nothing.

`Type`: NOOP

`Parameter`:  none

Example

```yaml
type: DELAY
name: delay
icon: accesibility.png 
title: Do Nothing
```

##### Delay

`Type`: DELAY

`Parameter`: 

`time`: time to delay in Seconds

Example

```yaml
type: DELAY
name: delay
parameters:
  time: 2
```

##### Timer

Starting a timer with a response every second. You can define the format of the timer message and the message on finish.

`Type`: TIMER

`Parameter`: 

`time`: time to delay in Seconds
`format`: the message for the response, defaults %d seconds
`finished`: the message at the end of the timer, defaults: finished

Example

```yaml
type: TIMER
name: timer
parameters:
  time: 10
  format: noch %ds
  finished: Fertig
```

##### Clock

Just a simple textual clock.

`Type`: CLOCK

`Parameter`: 

`format`: the format of the clock in Golang format syntax, defaults: 15:04:05
`analog`: true or false, shows an analog clock

Example

```yaml
type: CLOCK
name: clock
parameters:
  format: "15:04:05 02 Jan 06"
```

##### Execute

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

##### Page

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

##### Keys

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

##### Controlling Application Main Window

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

##### Screenshot, making a screenshot

With this command, you can take a screenshot. 

type: SCREENSHOT

Parameter:
saveto: the folder, where the screen shot will be saved. Format is `screen_<#number>_<display>.png`
display: optional, the number of the display, if you want to store screen shot of every display please use -1. Getting the right display, simply do a screen shot with display = -1. Than look at the screen shots and look in the name at the last number of the right image. That is your display.

```yaml
type: SCREENSHOT
name: screenshot
parameters:
  saveto: e:/temp/screenshot
  display: 1
```

##### Hardware monitor

This command connects to the openhardwaremonitor application on windows. With this you can get different sensors of your computer. For using the webserver of the openhardwaremonitor app, you have to add another external configurationinto the main service configuration. The url is the url to the app webserver added with data.json. the `updateperiod` is the update time in seconds. 

```yaml
extconfig:
  openhardwaremonitor:
	url: http://127.0.0.1:12999/data.json
	updateperiod: 5
```

If you have configured this, the service will evaluate on startup the connection and all possible sensor names. This lsit of names you will see in the log. The sensor name starts with the category, like CPU, GPU or Memory, followed by the hardware component. After that there is the sensor type like Clocks, Temperatures or Load, followed by the sensor name. To use a sensor you have to copy the whole name: like `"CPU/Intel Core i7-6820HQ/Load/CPU Total"`

e.g.: 

```
2021/02/18 08:47:17 found sensor with name: CPU/Intel Core i7-6820HQ/Clocks/Bus Speed
2021/02/18 08:47:17 found sensor with name: CPU/Intel Core i7-6820HQ/Clocks/CPU Core #1
2021/02/18 08:47:17 found sensor with name: CPU/Intel Core i7-6820HQ/Clocks/CPU Core #2
2021/02/18 08:47:17 found sensor with name: CPU/Intel Core i7-6820HQ/Clocks/CPU Core #3
2021/02/18 08:47:17 found sensor with name: CPU/Intel Core i7-6820HQ/Clocks/CPU Core #4
2021/02/18 08:47:17 found sensor with name: CPU/Intel Core i7-6820HQ/Temperature/CPU Core #1
2021/02/18 08:47:17 found sensor with name: CPU/Intel Core i7-6820HQ/Temperature/CPU Core #2 
2021/02/18 08:47:17 found sensor with name: CPU/Intel Core i7-6820HQ/Temperature/CPU Core #3 
2021/02/18 08:47:17 found sensor with name: CPU/Intel Core i7-6820HQ/Temperature/CPU Core #4 
2021/02/18 08:47:17 found sensor with name: CPU/Intel Core i7-6820HQ/Temperature/CPU Package 
2021/02/18 08:47:17 found sensor with name: CPU/Intel Core i7-6820HQ/Load/CPU Total 
2021/02/18 08:47:17 found sensor with name: CPU/Intel Core i7-6820HQ/Load/CPU Core #1 
2021/02/18 08:47:17 found sensor with name: CPU/Intel Core i7-6820HQ/Load/CPU Core #2 
2021/02/18 08:47:17 found sensor with name: CPU/Intel Core i7-6820HQ/Load/CPU Core #3 
2021/02/18 08:47:17 found sensor with name: CPU/Intel Core i7-6820HQ/Load/CPU Core #4 
2021/02/18 08:47:17 found sensor with name: CPU/Intel Core i7-6820HQ/Powers/CPU Package 
2021/02/18 08:47:17 found sensor with name: CPU/Intel Core i7-6820HQ/Powers/CPU Cores 
2021/02/18 08:47:17 found sensor with name: CPU/Intel Core i7-6820HQ/Powers/CPU Graphics 
2021/02/18 08:47:17 found sensor with name: CPU/Intel Core i7-6820HQ/Powers/CPU DRAM 
2021/02/18 08:47:17 found sensor with name: Memory/Generic Memory/Load/Memory 
2021/02/18 08:47:17 found sensor with name: Memory/Generic Memory/Data/Used Memory 
2021/02/18 08:47:17 found sensor with name: Memory/Generic Memory/Data/Available Memory 
2021/02/18 08:47:17 found sensor with name: GPU/AMD FirePro W5170M/Voltages/GPU Core 
2021/02/18 08:47:17 found sensor with name: GPU/AMD FirePro W5170M/Clocks/GPU Core 
2021/02/18 08:47:17 found sensor with name: GPU/AMD FirePro W5170M/Clocks/GPU Memory 
2021/02/18 08:47:17 found sensor with name: GPU/AMD FirePro W5170M/Temperature/GPU Core 
2021/02/18 08:47:17 found sensor with name: GPU/AMD FirePro W5170M/Load/GPU Core 
2021/02/18 08:47:17 found sensor with name: GPU/AMD FirePro W5170M/Fans/GPU Fan 
2021/02/18 08:47:17 found sensor with name: GPU/AMD FirePro W5170M/Data/GPU Fan 
2021/02/18 08:47:17 found sensor with name: Storage/HGST HTS721010A9E630/Temperature/Temperature 
2021/02/18 08:47:17 found sensor with name: Storage/HGST HTS721010A9E630/Load/Used Space 
2021/02/18 08:47:17 found sensor with name: Storage/PC300 NVMe SK hynix 256GB/Temperature/Temperature 2021/02/18 08:47:17 found sensor with name: Storage/PC300 NVMe SK hynix 256GB/Load/Used Space 
2021/02/18 08:47:17 found sensor with name: Storage/Generic Hard Disk/Load/Used Space
```

On the action side you have to configure this:

type: HARDWAREMONITOR

Parameter:
sensor: the sensor name like given above.
format: the format string for the textual representation
display: text, graph,  text shows only the textual representation, graph shows both
ymin: the value for the floor of the graph
ymax: the value for the bottom of the graph
color: color of the graph

```yaml
type: HARDWAREMONITOR
name: cpu
parameters:
  sensor: "CPU/Intel Core i7-6820HQ/Temperature/CPU Package"
  format: "%0.1f °C"
  display: text
  ymin: 30
  ymax: 80
  color: "#ff0000"
```

# UI

As this project is in an early development stage, i have nothing to show here. But if you like, you can simply checkout the develop branch, starting 2 Instances of Visual Studio Code, one with the frontend, one with the service backend. Than start the backend server in VS, and same for the frontend (npm run serve). After some hacking with configuration you may be able to see something. :-) Good Luck to you.