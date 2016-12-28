# GO Deploy

## This is just an exercise, not considered totally reliable

### A tool for building and deploying GO projects ###

#### Setup

* Install Git
* Install Docker
* Install Docker compose
* Generate an [ssh key](https://help.github.com/articles/generating-an-ssh-key/), then add it to git repository host and to all remote deploy destinations you will define in config.
* Clone this repository
* Create config file: cp config.yml.dist config.yml
* Install tool: ./build.sh

#### Configure projects

Define your projects in config.yml

* Each project will have a name
* repository: HTTP or SSH link to git repository
* branch: Branch to build project from
* destinations: Full paths (including filename) to deploy binary to
    * Pattern: user@host:/full/deploy/path/binary_filename

#### Directory structure

* bin: All binaries will be built here
* pkg: Go installed packages
* src: All projects sources will be cloned here
    * src/deploy: The deploy tool source code

#### Usage

From local machine: ./run.sh project_name

#### Workflow

On each run, the following actions occur:
* If project is deployed for the first time, it will be cloned (only the specified branch), else hard reset and rebase pull will be performed
* Dependencies are downloaded and installed
* Unit tests are ran
* The binary is built
* The binary is synced to all destinations

#### Other

Will work only with github and bitbucket
If you want to get into the container, run the following in the project's root directory: docker-compose run --rm deploy bash
Any changes can be done in src/deploy. Then build the deploy binary: ./build.sh

#### Nice to have

* Build from specific path in repository (multi projects repository)
* Logs
* Concurrent sync to all destinations
* Multiple ssh keys
* Repository hosts auto ssh scan (instead of scan in Dockerfile)
* Go workspace for each project
* Multiple go versions
