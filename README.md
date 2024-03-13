# Project Styx - POC

## What is Styx?

Styx currently is in Proof-of-Concept status and connects container-based services using a streamlined API.

## Features

### REST API

Container applications are required to implement the Styx REST API in order for Styx to send events and get status information.
See `client` directory for an example in TypeScript using Bun and Elysia.

* Styx REST API always starts with `/_styx`
* The default port is 8844 (but it can be reconfigured, see "Service discovery" below)
* The following endpoints need to be implemented:
  * `/_styx/init`
    * Called when Styx loads the service for the first time to get information.
    * Cached on main server for version emitted by "version" endpoint; or if reinitialize endpoint is triggered via deployment (TBD)
    * Expected response:
      ```json
      {
        "subscribedEvents": ["my-event"]
      }
      ```
    * "subscribedEvents" is a list of freely definable strings.
      Styx will pass the events that are triggered by other containers via Styx to the container, if it is in this list.
  *  `/_styx/event`
    * Receiving endpoint for events sent by Styx.
    * Expected request body:
    * ```json
      {"event": "my-event", "data": {"foo": "any data passed by the event"}}
      ```
  * `/_styx/status`
    * Status checked every 15 minutes to see if the service is still alive
    * Must return `{"status": "ready"}` if the service is operational. Else, no events etc will be passed to the service
    * Expected response:
      ```json
      {"status": "ready"}
      ```
  * `/_styx/version`
    * Return the current version of your application
    * Styx will cache the config from "init" step until new version is set. Make sure to set a new version number, when your list of subscribed events changes.
      * Note: In the future we may provide an endpoint for reinitializing the service directly
    * The value and format of the version does not matter. It may be an incremental number, a semantic version or anything else, as long as it can identify a release.
    * Expected response:
      ```json
      {"version": "1"}
      ```

### Service discovery

Services need to implement the Styx REST API. When they do, they can be discovered by Styx when they are on the same Docker network.

* If Styx runs in Docker, it will look for containers on the same networks
* If Styx runs outside of Docker, it will look for containers on the *bridge* network (Docker default)

There are two means of discovery: direct reference via environment variable, or automagic discovery via Docker labels

#### Environment variables

The environment variable `STYX_SERVICES` takes a comma-separated list of service names.
Use the same name as the service can be accessed on the Docker network.
You may add a port number after the service name if your service does not expose the API to default port 8844.

#### Docker labels

The containers need to be accessible for Styx, i.e. it needs access to the Docker socket.

All container with the following label are discovered automatically:
```dockerfile
LABEL me.saitho.styx.service=1
```

If your application exposes the Styx API on a different port, there is a label for that as well.
```dockerfile
LABEL me.saitho.styx.port=8843
```

*Nice to know:* Styx accesses the internal IP address, so there is no need to actually EXPOSE the port.

## Try

There are two setups available for testing:

Configuration via environment: `docker compose up backend-environment`
Then, open http://localhost:3322

Configuration via Docker Labels: `docker compose up backend-labels`
Then, open http://localhost:3322