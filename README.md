# Service Manager POC

## Overview

This repository contains the POC of the communication with the ServiceManager. The goal of this POC was to check if it is possible to reuse the client library that comes from the Service Manager CLI. The example flow implemented here creates the XSUAA service instance and a binding for that service. 

The Service Manager CLI contains the package that exposes a client library for communication with the Service Manager API. This library secures the communication using OAuth tokens but the platform registered in the Service Manager gives you the basic credentials. The POC shows the possibility to create the Service Manager client with the HTTP client that uses the BasicAuthTransport.

To run the POC, specify **APP_USERNAME**, **APP_PASSWORD**, and the Service Manager **APP_BASE_URL** and run this command:

```sh
APP_USERNAME=username APP_PASSWORD=passwd APP_BASE_URL=https://smurl.domain.local go run main.go
```
