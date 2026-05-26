# NullCloud - backend

NullCloud is a fake cloud service provider. Although you can provision resources through its API, the resources do not really exist. NullCloud backend and its terraform provider are useful to build demonstrations, run tests. It behaves as a real cloud service provider but with no real cloud resources.

## About

* Backend server for NullCloud.
* It implements the NullCloud API.
* It is written in Go.
* It can run on macOS, Linux, Windows

## Supported resources

NullCloud supports the following resources:
* VPC - virtual private cloud
* Subnet
* Virtual server instance

## API

NullCloud exposes REST APIs for its supported resources.
* NullCloud REST APIs require the use of an Authorization token
* The token can be any string, but can not be empty
* Resources will be linked to the given token
* Two different tokens do not see each other resources

## Persistence

NullCloud API state is persisted with every write. The backend offers different persistent mode:
* in-memory (useful for testing) - the default if no other persistent mode specified
* JSON file-based - the destination file is passed during launched as a command line parameter to the backend

API implementation does not know which persistent mode is used so an abstraction of the storage exists.

## Build

Use the provided Dockerfile to compile the backend.

## Test

Use the provider Dockerfile to compile and run the tests.
