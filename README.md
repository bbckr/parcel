# parcel
Simple go template renderer for nomad and other generic groups of templates (i.e. a light-weight Helm for purely templating).

## Command usage
### pull
Pull a parcel given the git source URI
``` bash
parcel pull [options] [SOURCE]

-f --force # (optional) force pull parcel
```
Examples:
``` bash
# pull from remote repository
parcel pull git::ssh://git@github.com/bbckr/parcel//example

# pull from local repository
parcel pull git::file://${PROJECT_ROOT}/.git//example

# pull from repository ref
parcel pull git::ssh://git@github.com/bbckr/parcel//example?ref=${GIT_REF}
```
### render
Render template output for a parcel given the parcel identifiers. The parcel must be pulled before it can be referenced.
``` bash
parcel render [options] [OWNER] [NAME] [VERSION]
-v --values # (optional) path to values yaml file
-o --output # directory to output rendered files (default: .)
```

Examples:
``` bash
# render to output directory
parcel render -o .output bbckr example 1.0.0

# override default values
parcel render -v myvalues.yaml bbckr example 1.0.0
```

## Documentation
### Parcel structure
Each parcel requires the following directory structure:
``` bash
folder_name
|_ static/ # place static files to reference in templates here
|_ templates/ # place all your templates here
...|_ job_1.hcl # go template, any name/ext
...|_ job_2.hcl # go template, any name/ext
|_ values.yaml # default values to feed parcel templates
|_ manifest.yaml # parcel manifest
```
### Referencing values in templates
Parcel uses basic Go templating with no added helper functions.

For a full example of how the `values.yaml` should look and how to reference the values and other injected variables in a template, checkout [this example](/example).

For quick reference:

``` ruby
# templates/sample-job.hcl
job "sample-jon" {
    group "sample-group" {
        task "sample-task" {
            config {
                image = "{{ .Values.image.name }}:{{ .Values.image.tag }}"
            }
        }
    }
}
```

``` yaml
# values.yaml
image:
  name: nginx
  tag: latest
```
Output:
``` ruby
job "sample-jon" {
    group "sample-group" {
        task "sample-task" {
            config {
                image = "nginx:latest"
            }
        }
    }
}
```
### Template Variables
- `.Meta.name` Parcel name, taken from manifest
- `.Meta.owner` Parcel owner, taken from manifest
- `.Meta.version` Parcel version, taken from manifest
- `.Values.*` Merged values from default Parcel values and values passed during render
- `.Static.path` Path to static directory to reference static resources in templates

### Template Functions
Uses [Go template](https://golang.org/pkg/text/template/) as a base and adds the following functions:
- `toJSON` Marshal to JSON string
- `replace` An alias for strings.Replace
- `replaceAll` An alias for strings.ReplaceAll
- `quote` Surround with double quotation marks
- `b64encode` Base64 encode string
- `b64decode` Base64 decode string

