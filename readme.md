
# Introduction
N/A


# Additions
I've added a *--colors *or *-c* flag to output the text version in color. Directory names will be blue, and symbolic links will be yellow.  Note that *--ouput text* is requred for colors to show up. 


# Required External Libraries
|Library|Github location|Description|
|-----|------|-------|
|Go Yaml | https://github.com/go-yaml/yaml | Because I'm not going to write a YAML parser|
|CodeGansta | https://github.com/codegangsta/cli | I'm really lazy and not going to write command line api when there looks to be four or five good ones out there|


At the command line, fetch the two libraries via the standard *go get* command

```
go get gopkg.in/yaml.v2
go get github.com/codegangsta/cli
```

# Samples to Run







