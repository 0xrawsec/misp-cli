# Description

Command line utility to query your favorite **MISP** instance.

Combined with the command `jq` you can achieve very nice things from command line.

## Configuration

Take example on the `config.json.example` found in this project.
By default the tool takes a configuration file named `config.json` in the directory where the binary is located.
If you want to change this behaviour, you can use the `-c` option to specify the path of another configuration file.

## Usage

```
Usage of misp-cli: misp-cli [OPTIONS]
  -a	Flag to search for attributes
  -c string
    	Configuration file to connect to MISP
  -cat string
    	Category to query
  -d	Enable debugging
  -e	Flag to search for events
  -eventid string
    	Event ID to look for
  -from string
    	Query events from date
  -l string
    	Last event query
  -org string
    	Organisation
  -tags string
    	Tags argument for query
  -to string
    	Query events until 'to' parameter
  -type string
    	Type argument for query
  -v string
    	Value to search for
  -version
    	Print version information
```

# Limitations

So far the CLI only supports queries to MISP.
