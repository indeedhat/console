# Console
a simple too to quickly build up a collection of scripts into a console application.

## How To
- fill out the .console.yml file
- build the binary
- put them in the same directory
- run it

## Config
```yaml
# give your console app a name
title: Console test
# Provides a brief description about the console
usage: this is an example console application for demonstratino purposes

# A list of commands available to the application
commands:
    # Key is what you need to type on the console to run
  - key: github
    # brief description of the command
    usage: Prints the github token from env
    # the body of the shell script
    # this will be written to a file and then executed
    cmd: |
      #!/bin/sh
      echo $GITHUB_TOKEN
    # Optional directory to run the script in
    # can be relative or absoulete
    workDir: ".git"

  - key: uptime
    usage: Prints the github token from env
    cmd: uptime

  - key: uname
    usage: Prints the full uname string
    cmd: |
      #!/bin/sh
      uname -a
```

## Usage
```console
./console
Console test
this is an example console application for demonstratino purposes

USAGE:
    ./console <command>

OPTIONS:

COMMANDS:
    github    Prints the github token from env
    uptime    Prints the github token from env
    uname     Prints the full uname string
```
