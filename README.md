# sm-xsuaa-poc

## Overview

This repository contains the POC of the communication with the ServiceManager. The goal of this POC was to check if it is possible to reuse the client library that comes from the ServiceManager CLI. The example flow implemented here creates the XSUAA service instance and then creates a binding for that service. 

The ServiceManager CLI contains the package that exposes a client library for the communications with the ServiceManager API. This library secures the communication using OAuth tokens but the platform registered in the ServiceManager gives you the basic credentials. The POC shows the possibility to create the ServiceManager client with the HTTP client that uses the BasicAuthTransport.

To run the POC you must specify username, password and the ServiceManager URL:

```sh
APP_USERNAME=username APP_PASSWORD=passwd APP_BASE_URL=https://smurl.domain.local go run main.go
```
