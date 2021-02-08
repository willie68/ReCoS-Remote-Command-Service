# go-micro

go-micro microservice framework
For project structure see: https://github.com/golang-standards/project-layout

# Profileconfiguration

## Action

### Command

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

