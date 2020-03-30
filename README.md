# parcel
Simple go template renderer for nomad and other generic groups of templates (i.e. a light-weight Helm for purely templating).

## Command usage
### pull
```
parcel pull [command options] [SOURCE]
```
### render
```
parcel render [command options] [SOURCE]
```

## Documentation
### Parcel structure
Each parcel requires the following directory structure:
``` bash
folder_name
|_ templates/ # place all your templates here
...|_ job_1.hcl # go template
...|_ job_2.hcl # go template
|_ values.yaml # default values to feed template
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
### Injected variables
TODO
