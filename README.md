# GoLib
golang library functions

## pkg/db
Contains helper functions for interacting with a postgreSQL database via gorm.
These functions expect a gorm query already started and provide some form of modification of the query.

functions include: debugging, filtering, paging, schemas, and sorting.

## pkg/io
Contains parsers for yaml and json to aid loading a json/yaml file into a struct.

## pkg/net
Contains network library for a generic server that performs a set function from parameters and needs a cache to speed up repeated requests.

## pkg/env
Contains helper functions for reading environment variables for an application

## pkg/mod
Provides ability to load generic functions from multiple modules
