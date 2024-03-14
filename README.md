# router-location-connector

## Table of Contents

- [Project Description](#project-description)
- [Design Details](#design-details)
- [Other Considerations](#other-considerations)
- [Running Tests](#running-tests)
- [Running Application](#running-application)
- [Still To Implement](#still-to-implement)

## Project Description:
This golang application accesses JSON data via a public  REST API about routers and their locations, stores the data in 
a redis instance and outputs to stdout a list of connections between their locations.

The public REST API(https://my-json-server.typicode.com/marcuzh/router_location_test_api/db)
returns JSON data on routers and their locations. Each router record contains the ID of its location and a list of links, which
represent the IDs of other routers it is connected to. e.g.

```json
    {
      "id": 1,
      "name": "citadel-01",
      "location_id": 1,
      "router_links": [
        1
      ]
    }
```

Each location record has the id, postcode and name of the location. e.g.

```json
    {
      "id": 1,
      "postcode": "BE12 2ND",
      "name": "Birmingham Motorcycle Museum"
    },
```

The links between routers are The links are bidirectional, meaning if router 1 is connected to router 2, then router 2 
is also connected to router 1.

Output should be one line per location, in the format [location name] <-> [location name], e.g.:

```shell
Adastral <-> London
London <-> Birmingham
Birmingham <-> Adastral
```
## Design Details:
The code is separated by packages based on mostly single responsibility apart from the `app` package which is the driver 
of processing the local flow of the program. 

Our `main` package is our entry point for executing the program. When you build a Go program into an executable binary, 
it must have a main package, and within that package, there must be a main function. This main function is where the 
execution of the program begins. package typically handles the setup of dependencies and dependency injection, although 
it doesn't follow the conventional patterns of dependency injection seen in other languages like Java or C#.
In Go, dependencies are typically imported directly within the main package or its subpackages. The init function, 
if present in the package, can be used to perform initialization tasks, including setting up dependencies.

In our app, we have a dependencies on Storage(Redis), and calling out to A REST API, so we instantiate clients
via their respective packages `New` constructor methods to interact with them. Thereafter, we inject them into an instance 
of the `app` object which has a `process` method which will dictate the flow of execution.

Once in the process function we firstly request the router location data via the api client and function call
`GetRouterLocationData`. The api package contains an `API` interface which the client is an instantiation of. This is an
example of the OO principal abstract as we define a set of methods in the interface that a type(client in this instance)
must implement. The interesting thing about the API. The beauty of Go is that while we use OO prinpals, we can also leverage
some benefits or functional programing as functions are first-class citizens. If you look the `ClientOptions` function
in the api package and the `options.go` file, you will see that we use functional options as a way to provide flexible and 
extensible configuration to Go structs. They allow users to customize the behavior of a struct by passing functional options 
to its constructor.

Once we retrieve the data from the REST API we then store it in Redis. There are many reasons for doing this but the main
are offline access, reducing network calls by caching data as well as storing raw data before transforming it along a 
processing pipeline. Redis was chosen due to its high performance, simplicity and caching and session storage.

The redis package also takes advantage of interfaces and type that implement it(`Redis`). There are also aspects of data 
encapsulation for the fields of the `Redis` struct. It contains a client for redis as well as a handler for a 
ReJson (client for json data type in Redis)

After storage, the next step is to process the data and calculate the links between router locations. An approach of 
crawling was taken using the recursive function `processLinkedRouter`

Its logic is simple in that it works its way through router links in the following way...

Given data for router id 3 and router id 15
```json
    {
      "id": 3,
      "name": "core-07",
      "location_id": 7,
      "router_links": [
        15
      ]
    },
    {
      "id": 15,
      "name": "cdn20",
      "location_id": 7,
      "router_links": [
        3,
        9
      ]
    }
```
router 3, has router links to router 15, but the above won't print a link for the link from router 3 to router 15 as they are in the same location. 
```location_id : 7 "Lancaster Castle"```
So then we focus on the next router links from router 15 to router 9.

```json
{
  "id": 9,
  "name": "edgesrv-01",
  "location_id": 8,
  "router_links": [
    14,
    15
  ]
}
```
router 9 has location id 8(`Loughborough University`)

so we have the link `"Lancaster Castle" <-> "Loughborough University"`

Since we are now at router 9, we then look for its other connected links, in this instance router 15.

```json
{
  "id": 15,
  "name": "cdn20",
  "location_id": 7,
  "router_links": [
    3,
    9
  ]
}
```

router 15 has a location_id of 7 ("Lancaster Castle") but we already have the link between location_id 7("Lancaster Castle") and location_id 8 ("Loughborough University")

but since we are already at router 15, we look to its links which we have already covered in looking at router 3 and router 9. So the linking for these 3 routers is now complete

As we process a router and its links and calculate the connections between the links we save ids of routers processed
to a map `processedRouters := make(map[int]struct{}, 0)` which allows for quick lookup access as we loop through all the routers to
see if we have already processed the router following a crawl of all the router links. 
We also store all location links to redis as this is a way to identify if the link has been processed before.

The routerLink is not taken from the raw JSON data but rather a new data model that we defined in order to ensure
bidirectional links are stored only once and have uniqueness to each pair. e.g:

```go
type RouterLocationLink struct {
	UniqueID   string // sort two locations alphabetically and concatenate them to ensure id always the same
	Connection string // in the format of `Location 1` <-> `Location 2`
}
```
### note
In its current implementation, data does not persist after each run of the application unless the 
run flag `persist-data` is set to true. The default of this flag is set to false as to allow printing of the locations as if it
is set to true the locations won't print as the logic to read from already fetched data is not set up as it always 
calls from the same API, and we should not assume data returned is always the same without a check, but if it is the same 
it will trigger a check that it already exists in the database and not print to stdout as we need to establish per run it 
printed to stdout already in that run.

## Other Considerations:
- Structured logging was added by the go package zerolog,
due to ease of use and flexibility.
- Gomock for mocking api and storage dependencies in unit tests


## Running Tests:

A docker-compose file, makefile and run_tests.sh script has been provided for your convenience

the easiest automated way is by running the following command in your terminal

```shell
 sh run_tests.sh
```

This script spins up containers as per the configuration in the docker-compose file and runs it in detached mode.
It then runs the unit and integration test commands respectively which can be found in the `Makefile`

NOTE: The services spun up by docker-compose are redis and [MockServer](https://www.mock-server.com/). MockServer was used due to its easy configuration and simplicity.
See `config/initializerJson.json` for the configuration as an example of the data returned from a request. This server is run in a docker container.


## Running Application:

A script per os environment has been provided, either `start_linux.sh` or `start_macos.sh`
e.g.

```shell
 sh start_macos.sh.sh
```

This runs a command from the Makefile to build an executable depending on the OS. It then spins up redis as one
of the applications dependencies and runs the executable built in the previous step

```shell
➜  router-location-connecter sh start_macos.sh
Building macOS executable
env GOOS=darwin GOARCH=amd64 go build -o router-location-connector ./cmd/router-location-connector
[+] Running 2/2
 ✔ Network router-location-connecter_default    Created                                                                                             0.0s 
 ✔ Container router-location-connecter-redis-1  Started                                                                                             0.2s 
2024/03/14 23:16:33 [DEBUG] GET https://my-json-server.typicode.com/marcuzh/router_location_test_api/db
[Williamson Park] <-> [Birmingham Hippodrome]
[Loughborough University] <-> [Lancaster Brewery]
[Lancaster Brewery] <-> [Lancaster University]
[Loughborough University] <-> [Lancaster Castle]
➜  router-location-connecter 

```

if you wish to change any of the run flags parsed in main you can run the go executable manually after 
running docker-compose redis -d and make build-macos or make build-linux

`./router-location-connector -base-url=https://my-json-server.typicode.com/marcuzh/router_location_test_api/db -retries=1 -timeout=10 -persist-data=false`

```shell
➜  router-location-connecter docker compose up redis -d                                                                                                                               
[+] Running 1/0
 ✔ Container router-location-connecter-redis-1  Running 

➜  router-location-connecter make build-macos                            
Building macOS executable
env GOOS=darwin GOARCH=amd64 go build -o router-location-connector ./cmd/router-location-connector
➜  router-location-connecter ./router-location-connector -base-url=https://my-json-server.typicode.com/marcuzh/router_location_test_api/db -retries=1 -timeout=10 -persist-data=false 

2024/03/14 23:38:46 [DEBUG] GET https://my-json-server.typicode.com/marcuzh/router_location_test_api/db
[Williamson Park] <-> [Birmingham Hippodrome]
[Loughborough University] <-> [Lancaster Brewery]
[Lancaster Brewery] <-> [Lancaster University]
[Loughborough University] <-> [Lancaster Castle]



```

## Still to implement
- Better test coverages as mostly happy paths were covered due to time constraints
- Concurrency: When writing to redis we could write locations and routers concurrently but need to change to a 
redis worker pool implementation and spin up goroutines and use waitgroups to wait before processing the data
- Logic for persisting data and reprinting routes (See [note](#note))
- Error monitoring with metrics

