# Recursive gotpl

This is a helper tool, that reads a data structure from files or 
environment variables and steps recursively through a folder structure rendering 
all templates it finds with the read in data, either in place or recreating the
folder structure in a target location. The templates have to follow the [go 
templagting syntax](https://golang.org/pkg/text/template/).


## Config

The tool itself takes it's config from a file in 
$HOME/.rgotpl/config.(yaml|json|toml) or relative to the executable 
in .rgotpl/config.(yaml|json|toml) while the later takes precedence over the 
first. The config file location can be overwritten by specifying the 
environment variable __RGTPL_CONFIG__ with the full path to the config file.
All fields from the config file can be overridden by an environment variable 
with an uppercase field name, with underscore delimiters 
instead of camelCase and prefixed with __RGTPL__, e.g. 
__RGTPL_LOGGING.LOG_LEVEL=debug__.


### Example config

```yaml
logging:
    logLevel: debug # one of debug, error defaults to error
    logFormat: text # one of json or text, defaults to json

templates:
    missingKey: zero # one of zero, invalid, error, defaults to error as described [here](https://golang.org/pkg/text/template/#Template.Option) 
    sourcePath: /var/data/ci-templates #
    targetPath: /usr/local/ci/pipelines # if empty templates get rendered in place defaults to ""
``` 


## Example data and template
The data can be provided as yaml, json or toml files which is parsed into a go
data structure and then passed to the go templating engine. Also all environment
variables are added to the data structure under the field __env__ and then the 
variable name camelCased. Data from the input files under the key env is overriden.
This could be used to provide defaults but the data structures are not merged but
overridden as a whole.

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
    "" ""
  }
}

```

and environment variables set like this

```bash
export DOCKER.REPO=somerepo
export HOME=/home/templater
export USER=hmuendel
export PASSWORD=supersecret

```

would be rendered like this, it the missingKey property would be invalid

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


