# Recursive gotpl

This is a helper tool, that reads a data structure from files or 
environment variables and steps recursively through a folder structure rendering 
all templates it finds with the read in data, either in place or recreating the
folder structure in a target location. The templates have to follow the [go 
templagting syntax](https://golang.org/pkg/text/template/).


## Config

The tool itself takes it's config from a file relative to the executable 
in .rgotpl/config.(yaml|json|toml). The config file location can be overwritten by specifying the 
environment variable __RGTPL_CONFIG__ with the full path to the config file.
All fields from the config file can be overridden by an environment variable 
with an uppercase field name, with underscore delimiters for nested values 
and prefixed with __RGTPL___, e.g. 
`RGTPL_LOG_LEVEL=7`.


### Example config

```yaml
log:
    level: 3 # loglevel 0-10 with 10 most verbose, defaults to 0
    logDir: /var/log/rgtpl/ # directory where to write log files if empty logging to stderr
    vmodule: "*config*.go=10" # pattern to match modules files and set different log level 
      
template:
    missingKey: zero # one of zero, invalid, error, defaults to error as described [here](https://golang.org/pkg/text/template/#Template.Option) 
    sourcePath: /var/data/ci-templates #
    targetPath: /usr/local/ci/pipelines # if empty templates get rendered in place defaults to ""
``` 


## Example data and template
The data can be provided as yaml, which is parsed into a go
data structure and then passed to the go templating engine. Also all environment
variables are added to the data structure under the field __env__ and then the 
variable name camelCased. Data from the input files under the key env is overridden as a whole.

#### Example
the template file:

```yaml
platform: linux

image_resource:
  type: docker-image
  source: {{ .docker.repo }}
  tag: {{ .docker.tag }}
  username: {{ .env.user }}
  password: {{ .env.password }}
  client_certs: 
  - domain: private.registry.com
    cert: {{ .env.cert }}     
    key: {{ .env.key }}
```
with the data file

```json
{
  "docker": {
    "repo": "hmuendel/recursive-gotpl",
    "tag": "latest" 
  },
  "env": {
    "user": "defaultUser",
     "cert": "dummyCert",
     "key": "dummyKey"
  }
}

```

and environment variables set like this

```bash
export HOME=/home/templater
export USER=hmuendel
export PASSWORD=supersecret

```

would be rendered like this, if the missingKey property would be set to _invalid_

```yaml
platform: linux

image_resource:
  type: docker-image
  source: hmuendel/recursive-gotpl
  tag: latest
  username: hmuendel
  password: supersecret
  client_certs: 
  - domain: private.registry.com
    cert: <no value>    
    key: <no value>
```
## verbose logging

Since I could not find a strict definition,
I interpreted google's verbose logging (glog)[https://github.com/google/glog]
in the following way.

Verbosity increases with increased v level where an increase to an odd level shows
a more detailed control flow of the program while an increase to an even level
logs data structures in the detail of the previous log level.
All higher v levels log everything the lower levels did.


|v  |verbosity|
|---|---------|
|0  | unrecoverable situations, startup and finish
|1  | errors and malformed conditions with maybe unwanted but recoverable situations
|2  | showing erroneous and malformed data
|3  | logging at general control points of the code,like entering and leaving modules
|4  | showing important data at control points
|5  | logging branching of control flow, like function calls
|6  | showing important data from function calls and return values
|7  | logging in function conditional branches and such
|8  | showing data that lead to decisions in functions and results of transformations
|9  | logging single events of everything, each request 
|10 | showing all data of every 
    


