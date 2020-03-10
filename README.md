<img align="right" width="290px" src=gopher.png>

```text
             _   _     _       
 ___ ___ ___| |_|_|___| |_ ___ 
| . | .'|  _|   | |   | '_| . |
|  _|__,|___|_|_|_|_|_|_,_|___|
|_|

modular, plugin based media sorter  
created by @rbtr  
```
---

[![Build Status](https://cloud.drone.io/api/badges/rbtr/pachinko/status.svg)](https://cloud.drone.io/rbtr/pachinko)
[![Go Report Card](https://goreportcard.com/badge/github.com/rbtr/pachinko)](https://goreportcard.com/report/github.com/rbtr/pachinko)
[![Release](https://img.shields.io/github/release/rbtr/pachinko.svg)](https://github.com/rbtr/pachinko/releases/latest)
[![Docker](https://img.shields.io/docker/pulls/rbtr/pachinko)](https://hub.docker.com/r/rbtr/pachinko)
[![License](https://img.shields.io/github/license/rbtr/pachinko)](/LICENSE)


### what it is
pachinko is a media sorter. it integrates with the tvdb and the moviedb to, given a directory of reasonably named mix media, organize that media into a clean hierarchal directory structure ideal for use in media servers like plex, kodi/xbmc, etc.

unlike some of the prior implementations of this idea, pachinko was designed from inception to be automation and container-friendly.  
it has no heavy gui - configure it through the config file or via flags, then execute it and walk away.

it is written in go so that it is approachable for anyone interested in contributing without sacrificing performance.  
the plugin-style architecture keeps the core codebase clear and efficient.

### design

#### plugins
pachinko has a plugin based pipeline design. the base plugin types are:
- input - add data from a datasource to the stream
- processor - modify the datastream in-flight
- output - deal with the processed data by e.g. moving files from their original dir to their sorted path.

these base plugin types make pachinko flexible. composing a pipeline of many combination of plugins is possible.

additionally there are subtypes of `processor` plugins:
- preprocessor - parse data already present in the datastream to classify, clean, or add information to the data before main processor plugins run. the preprocessors make modifications to the datastream based only on the data already present in the objects in the datastream.
- (intra)processor - the main working processors, where external datasources may be queried to enrich the datastream and significant modifications made. 
- postprocessor - last chance to modify the datastream before it is sent to outputs, but after the rich data has been added. the postprocessors make final modifications that shouldn't be the responsibility of the intraprocessors but may depend on the data enrichments that those have added.

they subtypes exist mainly to allow ordering of plugin flow. 

#### datatypes
pachinko currently supports these data types: 
- tv and 
- movie video files

other datatypes planned include: images (and whatever you would like to contribute!)

#### inputs
pachinko currently supports these inputs: 
- local filesystem (`path`). 

other datastore types planned include : s3 (and whatever you would like to contribute!)

#### outputs
pachinko currently supports these outputs:
- local filesystem (`path_mover`)
- stdout (`logger`)
- [trakt collector (`trakt_collector`)](docs/plugins/outputs/trakt.md)

#### processors
pachinko has the following required processors:
- categorizer (internal) 

pachinko has the following optional processors:
- tv identifier (pre-tv)
- movie identifier (pre-movie)
- tvdb (intra-tvdb)
- tmdb (intra-tmdb)
- tv path solver (post-tv_path_solver)
- movie path solver (post-movie_path_solver)

### how to run it
pachinko is distributed as a container and as a cross-platform binary.  

the container is recommended:
```bash
$ docker run -v /path/to/source:/src:z -v /path/to/dest:/dest:z -v /path/to/cfg:/cfg rbtr/pachinko:latest --config /cfg
```

to run the binary:
```bash
$ ./pachinko sort --config /path/to/config
```

### options
pachinko is configurable via file (yaml, toml), cli flags, or env vars.

the config file is recommended:
```yaml
dry-run: true
log-level: debug
inputs: []
outputs: []
processors: {}
```

the full, current list of options is available by running `./pachinko genconf` on the commandline.  
the core pachinko options are:

| option | inputs | usage |
| - | - | - |
| conf | string | full path to config file - ignored in the config file | 
| dry-run | bool | dry-runs print only, pachkino will not make changes |
| log-level | string | one of (trace,debug,info,warn,error) for logging verbosity |
| log-format | string | one of (json,text) | 


inputs, outputs, and processors are lists of plugins objects and look generally like:

```yaml
inputs:
- name: path
  src-dir: /path/to/source
outputs:
- name: stdout
```

note that each plugin may have its own independent config options; refer to that plugin's docs for details on configuring that specific plugin. here, the `path` input plugin has a `src-dir` parameter that we configure in the plugin list item.

the plugin list is processed in the written order and repeats are allowed. all loaded plugins are guaranteed to see each of the items in the datastream at least once. if the order that your datastream is processed by each plugin matters, make sure to load your plugins in the correct order!


### testimonials

here's what users had to say when asked what they thought about `pachinko`:

> Ew. _Pachinko_? Why would you name it _pachinko_? _Pachinko_ just makes me think of flashing lights and cigarette smoke. - a Japanese user

### license

`pachinko` is licensed under MPL-2 which generally means you can use it however you like with few exceptions.  

read the full license terms [here](https://www.mozilla.org/en-US/MPL/2.0/FAQ/).  

---

created by @rbtr  
inspired by the functionality and frustrating user experience of: [sorttv by cliffe](https://sourceforge.net/projects/sorttv/), filebot, tinymediamanager, and others  
and the excellent architecture patterns of telegraf, caddy, coredns, and others
