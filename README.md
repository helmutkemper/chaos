# Chaos Test

> This is version 3.0, still under development

This project aims to create chaos testing for microservices, allowing to transform a simple golang test into a chaos 
test.

The focus of this project is to allow the chaos test still in the development of the project and try to solve the 
famous problem, on my machine it works!

The test consists of creating all the necessary infrastructure for the project to work on the developer's machine, 
using docker, and after that, pausing or dropping containers, stopping the data flow in the middle of the process.

Imagine making a microservice where three instances must keep data up to date with each other at all times.

The test allows you to create a container for each instance, simulate the data and leave the data flowing, while the 
containers are paused or restarted at random, pausing the transmission of data unexpectedly, allowing you to capture 
failures such as loss of connection or excessive delay in the transmission of data.


## Basic usage

### Using a git server

```go
package githubproject

import (
  "github.com/helmutkemper/chaos/factory"
  "testing"
  "time"
)

func TestLinear(t *testing.T) {

  primordial := factory.NewPrimordial().
    NetworkCreate("test_network", "10.0.0.0/16", "10.0.0.1").
    Test(t)

  factory.NewContainerFromGit(
    "server:latest",
    "https://github.com/helmutkemper/chaos.public.example.git",
  ).
    Ports("tcp", 3000, 3000).
    Create("server", 1).
    Start()

  if !primordial.Monitor(3 * time.Minute) {
    t.Fail()
  }
}
```

### Using a docker image

```go
package mongodbproject

import (
  "github.com/helmutkemper/chaos/factory"
  "testing"
  "time"
)

func TestLinear(t *testing.T) {

  primordial := factory.NewPrimordial().
    NetworkCreate("test_network", "10.0.0.0/16", "10.0.0.1").
    Test(t)

  factory.NewContainerFromImage(
    "mongo:latest",
  ).
    Ports("tcp", 27017, 27017).
    EnvironmentVar([]string{"--bind_ip_all"}).
    Create("mongo", 1).
    Start()

  if !primordial.Monitor(3 * time.Minute) {
    t.Fail()
  }
}
```

### Using a local folder

```go
package localFolderProject

import (
  "github.com/helmutkemper/chaos/factory"
  "testing"
  "time"
)

func TestDevOps_Linear(t *testing.T) {
  primordial := factory.NewPrimordial().
    NetworkCreate("test_network", "10.0.0.0/16", "10.0.0.1").
    Test(t)

  factory.NewContainerFromFolder(
    "folder:latest",
    "./project",
  ).
    MakeDockerfile().
    FailFlag("../bug", "panic:", "bug:", "error").
    Create("folder", 3).
    Start()

  if !primordial.Monitor(15 * time.Minute) {
    t.FailNow()
  }
}
```

## Example

```go
package localFolderProject

import (
  "github.com/helmutkemper/chaos/factory"
  "testing"
  "time"
)

func TestDevOps_Linear(t *testing.T) {
  primordial := factory.NewPrimordial().
    NetworkCreate("chaos_network", "10.0.0.0/16", "10.0.0.1").
    Test(t)

  factory.NewContainerFromImage("nats:latest").
    EnableChaos(2,2,2,0.0).
    FailFlag("./bug", "panic:", "bug:", "error").
    SaveStatistics("./").
    Ports("tcp", 4222, 4222, 4223, 4224).
    Create("nats", 3).
    Start()

  factory.NewContainerFromFolder(
    "folder:latest",
    "./project",
  ).
    MakeDockerfile().
    EnableChaos(2,2,2,0.0).
    FailFlag("./bug", "panic:", "bug:", "error").
    SaveStatistics("./").
    Create("folder", 3).
    Start()

  if !primordial.Monitor(60 * time.Minute) {
    t.FailNow()
  }
}
```

> The `main.go` and `go.mod` files must be contained in the root folder of the project

## Simulate network problems

```go
package mongodbproject

import (
	"context"
	"fmt"
	"github.com/helmutkemper/chaos/factory"
	"go.mongodb.org/mongo-driver/bson"
	"runtime/debug"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func TestLinear(t *testing.T) {

	primordial := factory.NewPrimordial().
		NetworkCreate("test_network", "10.0.0.0/16", "10.0.0.1").
		Test(t)

	factory.NewContainerFromImage(
		"mongo:latest",
	).
		Ports("tcp", 27017, 27017).
		EnvironmentVar([]string{"--bind_ip_all"}).
		Create("mongo", 1).
		Start()

	factory.NewContainerNetworkProxy(
		"delay",
		27016,
		"delete_mongo_0:27017",
		10, 100,
		0.0,
	)

	go mongoPopulate(t)

	if !primordial.Monitor(2 * time.Minute) {
		t.Fail()
	}
}

func mongoPopulate(t *testing.T) {
	var err error
	var mongoClient *mongo.Client
	var start = time.Now()

	fmt.Printf("conexão\n")

	mongoClient, err = mongo.NewClient(options.Client().ApplyURI("mongodb://0.0.0.0:27016"))
	if err != nil {
		panic(string(debug.Stack()))
	}

	err = mongoClient.Connect(context.Background())
	if err != nil {
		panic(string(debug.Stack()))
	}

	err = mongoClient.Ping(context.Background(), readpref.Primary())
	if err != nil {
		fmt.Printf("error: %v\n", err.Error())
		panic(string(debug.Stack()))
	}

	type Trainer struct {
		Name string
		Age  int
		City string
	}

	var collection *mongo.Collection
	var totalOfInserts int64 = 1000
	for i := int64(0); i != totalOfInserts; i += 1 {
		collection = mongoClient.Database("test").Collection("trainers")
		ash := Trainer{"Ash", 10, "Pallet Town"}
		_, err = collection.InsertOne(context.Background(), ash)
		if err != nil {
			panic(err)
		}
	}

	var total int64
	if total, err = collection.CountDocuments(context.Background(), bson.M{"name": "Ash", "age": 10, "city": "Pallet Town"}); err != nil {
		panic(err)
	}

	if total != totalOfInserts {
		t.Logf("total of inserts must be %v found %v", totalOfInserts, total)
		t.Fail()
	}

	fmt.Printf("fim\n")
	duration := time.Since(start)
	fmt.Printf("Duration: %v\n\n", duration)
}
```


## Important

* All names in docker elements, created by chaos test, start by `delete_` and end by `_` + a sequential number (0,1,2...);
* `NetworkCreate("test_network", "10.0.0.0/16", "10.0.0.1")` gives to the first container an IP `10.0.0.2`, but, you can use container name as address.

# Vulnerability Report Example

> Automatic function

This report is based on an open database and shows known vulnerabilities. Validity: Thu Dec 22 18:29:28 2022

## Path

Path: /scan/go.mod
Type: lockfile

### Packages

| Ecosystem | Package          | Version                           |
|-----------|------------------|-----------------------------------|
| Go        | golang.org/x/net | 0.0.0-20220225172249-27dd8689420f |

### Details:

HTTP/2 server connections can hang forever waiting for a clean shutdown that was preempted by a fatal error. This condition can be exploited by a malicious client to cause a denial of service.

### Affected:

| Ecosystem | Package          |
|-----------|------------------|
| Go        | stdlib           |
| Go        | golang.org/x/net |

| type   | URL                                                                                                                    |
|--------|------------------------------------------------------------------------------------------------------------------------|
| WEB    | [https://groups.google.com/g/golang-announce/c/x49AQzIVX-s](https://groups.google.com/g/golang-announce/c/x49AQzIVX-s) |
| REPORT | [https://go.dev/issue/54658](https://go.dev/issue/54658)                                                               |
| FIX    | [https://go.dev/cl/428735](https://go.dev/cl/428735)                                                                   |

### Details:

An attacker can cause excessive memory growth in a Go server accepting HTTP/2 requests.

HTTP/2 server connections contain a cache of HTTP header keys sent by the client. While the total number of entries in this cache is capped, an attacker sending very large keys can cause the server to allocate approximately 64 MiB per open connection.

### Affected:

| Ecosystem | Package          |
|-----------|------------------|
| Go        | stdlib           |
| Go        | golang.org/x/net |

| type   | URL                                                                                                                                                  |
|--------|------------------------------------------------------------------------------------------------------------------------------------------------------|
| REPORT | [https://go.dev/issue/56350](https://go.dev/issue/56350)                                                                                             |
| FIX    | [https://go.dev/cl/455717](https://go.dev/cl/455717)                                                                                                 |
| FIX    | [https://go.dev/cl/455635](https://go.dev/cl/455635)                                                                                                 |
| WEB    | [https://groups.google.com/g/golang-announce/c/L_3rmdT0BMU/m/yZDrXjIiBQAJ](https://groups.google.com/g/golang-announce/c/L_3rmdT0BMU/m/yZDrXjIiBQAJ) |
## Path

Path: /scan/go.mod
Type: lockfile

### Packages

| Ecosystem | Package           | Version |
|-----------|-------------------|---------|
| Go        | golang.org/x/text | 0.3.7   |

### Details:

An attacker may cause a denial of service by crafting an Accept-Language header which ParseAcceptLanguage will take significant time to parse.

### Affected:

| Ecosystem | Package           |
|-----------|-------------------|
| Go        | golang.org/x/text |

| type   | URL                                                                                                                                                  |
|--------|------------------------------------------------------------------------------------------------------------------------------------------------------|
| REPORT | [https://go.dev/issue/56152](https://go.dev/issue/56152)                                                                                             |
| FIX    | [https://go.dev/cl/442235](https://go.dev/cl/442235)                                                                                                 |
| WEB    | [https://groups.google.com/g/golang-announce/c/-hjNw559_tE/m/KlGTfid5CAAJ](https://groups.google.com/g/golang-announce/c/-hjNw559_tE/m/KlGTfid5CAAJ) |

# Memory and CPU log example

> Automatic function

| time                | state - running | state - dead | state - OOMKilled | state - paused | state - restarting | state - error | state - status | state - exitCode | state - health check | read      | pre read  | pids - current (linux) | pids - limit (linux) | num of process (windows) | storage - read count (windows) | storage - write count (windows) | cpu - online | cpu - system usage | cpu - usage in user mode | cpu - usage in kernel mode | cpu - total usage | cpu - throttled time | cpu - throttled periods | cpu - throttling periods | pre cpu - online | pre cpu - system usage | pre cpu - usage in user mode | pre cpu - usage in kernel mode | pre cpu - total usage | pre cpu - throttled time | pre cpu - throttled periods | pre cpu - throttling periods | memory - limit | memory - commit peak | memory - commit | memory - fail cnt | memory - usage | memory - max usage | memory - private working set |
|---------------------|-----------------|--------------|-------------------|----------------|--------------------|---------------|----------------|------------------|----------------------|-----------|-----------|------------------------|----------------------|--------------------------|--------------------------------|---------------------------------|--------------|--------------------|--------------------------|----------------------------|-------------------|----------------------|-------------------------|--------------------------|------------------|------------------------|------------------------------|--------------------------------|-----------------------|--------------------------|-----------------------------|------------------------------|----------------|----------------------|-----------------|-------------------|----------------|--------------------|------------------------------|
| 2022-12-22 18:05:17 | true            | false        | false             | false          | false              |               | running        |                  |                      | 177613586 | 166940794 | 6                      | -1                   | 0                        | 0                              | 0                               | 8            | 184040640000000    | 35303000                 | 39377000                   | 74681000          | 0                    | 0                       | 0                        | 8                | 184032720000000        | 32301000                     | 39031000                       | 71333000              | 0                        | 0                           | 0                            | 12544057344    | 0                    | 0               | 0                 | 2617344        | 0                  | 0                            |
| 2022-12-22 18:05:27 | true            | false        | false             | false          | false              |               | running        |                  |                      | 171544716 | 164254632 | 6                      | -1                   | 0                        | 0                              | 0                               | 8            | 184118820000000    | 54666000                 | 50461000                   | 105127000         | 0                    | 0                       | 0                        | 8                | 184110970000000        | 52911000                     | 48841000                       | 101752000             | 0                        | 0                           | 0                            | 12544057344    | 0                    | 0               | 0                 | 2727936        | 0                  | 0                            |
| 2022-12-22 18:05:37 | true            | false        | false             | false          | false              |               | running        |                  |                      | 171077595 | 166890387 | 6                      | -1                   | 0                        | 0                              | 0                               | 8            | 184196780000000    | 83039000                 | 72275000                   | 155315000         | 0                    | 0                       | 0                        | 8                | 184188900000000        | 82306000                     | 71636000                       | 153942000             | 0                        | 0                           | 0                            | 12544057344    | 0                    | 0               | 0                 | 2969600        | 0                  | 0                            |
| 2022-12-22 18:05:48 | true            | false        | false             | false          | false              |               | running        |                  |                      | 263693003 | 230767753 | 6                      | -1                   | 0                        | 0                              | 0                               | 8            | 184284240000000    | 141004000                | 125160000                  | 266165000         | 0                    | 0                       | 0                        | 8                | 184276070000000        | 132987000                    | 118391000                      | 251378000             | 0                        | 0                           | 0                            | 12544057344    | 0                    | 0               | 0                 | 3538944        | 0                  | 0                            |
| 2022-12-22 18:05:57 | true            | false        | false             | false          | false              |               | running        |                  |                      | 178493424 | 166287840 | 6                      | -1                   | 0                        | 0                              | 0                               | 8            | 184354520000000    | 196078000                | 181249000                  | 377327000         | 0                    | 0                       | 0                        | 8                | 184346570000000        | 189827000                    | 167112000                      | 356939000             | 0                        | 0                           | 0                            | 12544057344    | 0                    | 0               | 0                 | 4067328        | 0                  | 0                            |
| 2022-12-22 18:06:07 | true            | false        | false             | false          | false              |               | running        |                  |                      | 174253762 | 166177428 | 6                      | -1                   | 0                        | 0                              | 0                               | 8            | 184433270000000    | 255357000                | 251859000                  | 507217000         | 0                    | 0                       | 0                        | 8                | 184425320000000        | 253246000                    | 242839000                      | 496085000             | 0                        | 0                           | 0                            | 12544057344    | 0                    | 0               | 0                 | 4505600        | 0                  | 0                            |
| 2022-12-22 18:06:17 | true            | false        | false             | false          | false              |               | running        |                  |                      | 174244961 | 164435544 | 7                      | -1                   | 0                        | 0                              | 0                               | 8            | 184512140000000    | 334782000                | 312463000                  | 647245000         | 0                    | 0                       | 0                        | 8                | 184504170000000        | 325585000                    | 308840000                      | 634425000             | 0                        | 0                           | 0                            | 12544057344    | 0                    | 0               | 0                 | 5455872        | 0                  | 0                            |
| 2022-12-22 18:06:27 | true            | false        | false             | false          | false              |               | running        |                  |                      | 175264632 | 164967173 | 7                      | -1                   | 0                        | 0                              | 0                               | 8            | 184590980000000    | 388605000                | 364317000                  | 752922000         | 0                    | 0                       | 0                        | 8                | 184583060000000        | 385493000                    | 359421000                      | 744915000             | 0                        | 0                           | 0                            | 12544057344    | 0                    | 0               | 0                 | 5894144        | 0                  | 0                            |
| 2022-12-22 18:06:37 | true            | false        | false             | false          | false              |               | running        |                  |                      | 180010095 | 171461678 | 8                      | -1                   | 0                        | 0                              | 0                               | 8            | 184669860000000    | 442193000                | 434521000                  | 876714000         | 0                    | 0                       | 0                        | 8                | 184661890000000        | 442193000                    | 427139000                      | 869332000             | 0                        | 0                           | 0                            | 12544057344    | 0                    | 0               | 0                 | 6246400        | 0                  | 0                            |
| 2022-12-22 18:06:47 | true            | false        | false             | false          | false              |               | running        |                  |                      | 302193294 | 252033752 | 8                      | -1                   | 0                        | 0                              | 0                               | 8            | 184749540000000    | 510029000                | 510029000                  | 1020058000        | 0                    | 0                       | 0                        | 8                | 184741290000000        | 503274000                    | 503274000                      | 1006549000            | 0                        | 0                           | 0                            | 12544057344    | 0                    | 0               | 0                 | 6516736        | 0                  | 0                            |
| 2022-12-22 18:06:57 | true            | false        | false             | false          | false              |               | running        |                  |                      | 172970549 | 165068132 | 8                      | -1                   | 0                        | 0                              | 0                               | 8            | 184827440000000    | 589323000                | 579939000                  | 1169263000        | 0                    | 0                       | 0                        | 8                | 184819470000000        | 576633000                    | 574761000                      | 1151395000            | 0                        | 0                           | 0                            | 12544057344    | 0                    | 0               | 0                 | 6774784        | 0                  | 0                            |
| 2022-12-22 18:07:07 | true            | false        | false             | false          | false              |               | running        |                  |                      | 177549762 | 167924220 | 8                      | -1                   | 0                        | 0                              | 0                               | 8            | 184906300000000    | 658255000                | 637328000                  | 1295584000        | 0                    | 0                       | 0                        | 8                | 184898330000000        | 652189000                    | 633175000                      | 1285364000            | 0                        | 0                           | 0                            | 12544057344    | 0                    | 0               | 0                 | 6422528        | 0                  | 0                            |
| 2022-12-22 18:07:17 | true            | false        | false             | false          | false              |               | running        |                  |                      | 179934169 | 169975460 | 8                      | -1                   | 0                        | 0                              | 0                               | 8            | 184985180000000    | 731492000                | 693689000                  | 1425182000        | 0                    | 0                       | 0                        | 8                | 184977220000000        | 724967000                    | 690806000                      | 1415773000            | 0                        | 0                           | 0                            | 12544057344    | 0                    | 0               | 0                 | 6512640        | 0                  | 0                            |
| 2022-12-22 18:07:27 | true            | false        | false             | false          | false              |               | running        |                  |                      | 171513382 | 163851090 | 8                      | -1                   | 0                        | 0                              | 0                               | 8            | 185063940000000    | 806551000                | 754947000                  | 1561498000        | 0                    | 0                       | 0                        | 8                | 185056020000000        | 794624000                    | 747098000                      | 1541723000            | 0                        | 0                           | 0                            | 12544057344    | 0                    | 0               | 0                 | 6586368        | 0                  | 0                            |
| 2022-12-22 18:07:37 | true            | false        | false             | false          | false              |               | running        |                  |                      | 178869179 | 164901553 | 8                      | -1                   | 0                        | 0                              | 0                               | 8            | 185142840000000    | 874769000                | 811877000                  | 1686646000        | 0                    | 0                       | 0                        | 8                | 185134860000000        | 863826000                    | 808404000                      | 1672231000            | 0                        | 0                           | 0                            | 12544057344    | 0                    | 0               | 0                 | 6623232        | 0                  | 0                            |
| 2022-12-22 18:07:47 | true            | false        | false             | false          | false              |               | running        |                  |                      | 300618253 | 249422210 | 8                      | -1                   | 0                        | 0                              | 0                               | 8            | 185222770000000    | 941152000                | 880063000                  | 1821216000        | 0                    | 0                       | 0                        | 8                | 185214450000000        | 938897000                    | 866085000                      | 1804983000            | 0                        | 0                           | 0                            | 12544057344    | 0                    | 0               | 0                 | 6709248        | 0                  | 0                            |
| 2022-12-22 18:07:57 | true            | false        | false             | false          | false              |               | running        |                  |                      | 178222466 | 168048840 | 8                      | -1                   | 0                        | 0                              | 0                               | 8            | 185300740000000    | 997041000                | 949381000                  | 1946422000        | 0                    | 0                       | 0                        | 8                | 185292790000000        | 987244000                    | 941591000                      | 1928835000            | 0                        | 0                           | 0                            | 12544057344    | 0                    | 0               | 0                 | 6909952        | 0                  | 0                            |
| 2022-12-22 18:08:07 | true            | false        | false             | false          | false              |               | running        |                  |                      | 172771137 | 164513470 | 8                      | -1                   | 0                        | 0                              | 0                               | 8            | 185379560000000    | 1039701000               | 1032180000                 | 2071881000        | 0                    | 0                       | 0                        | 8                | 185371600000000        | 1027553000                   | 1029435000                     | 2056989000            | 0                        | 0                           | 0                            | 12544057344    | 0                    | 0               | 0                 | 7127040        | 0                  | 0                            |
| 2022-12-22 18:08:17 | true            | false        | false             | false          | false              |               | running        |                  |                      | 176381544 | 166040377 | 8                      | -1                   | 0                        | 0                              | 0                               | 8            | 185458440000000    | 1094075000               | 1109192000                 | 2203268000        | 0                    | 0                       | 0                        | 8                | 185450470000000        | 1082076000                   | 1099072000                     | 2181148000            | 0                        | 0                           | 0                            | 12544057344    | 0                    | 0               | 0                 | 7258112        | 0                  | 0                            |
| 2022-12-22 18:08:27 | true            | false        | false             | false          | false              |               | running        |                  |                      | 169803091 | 162617257 | 8                      | -1                   | 0                        | 0                              | 0                               | 8            | 185537320000000    | 1149111000               | 1179451000                 | 2328563000        | 0                    | 0                       | 0                        | 8                | 185529390000000        | 1144102000                   | 1166871000                     | 2310974000            | 0                        | 0                           | 0                            | 12544057344    | 0                    | 0               | 0                 | 7483392        | 0                  | 0                            |
| 2022-12-22 18:08:37 | true            | false        | false             | false          | false              |               | running        |                  |                      | 174729804 | 163236137 | 8                      | -1                   | 0                        | 0                              | 0                               | 8            | 185616390000000    | 1193547000               | 1234899000                 | 2428447000        | 0                    | 0                       | 0                        | 8                | 185608370000000        | 1181820000                   | 1232712000                     | 2414533000            | 0                        | 0                           | 0                            | 12544057344    | 0                    | 0               | 0                 | 7643136        | 0                  | 0                            |
| 2022-12-22 18:08:47 | true            | false        | false             | false          | false              |               | running        |                  |                      | 272978878 | 230336502 | 8                      | -1                   | 0                        | 0                              | 0                               | 8            | 185696140000000    | 1231568000               | 1303115000                 | 2534683000        | 0                    | 0                       | 0                        | 8                | 185687950000000        | 1231568000                   | 1290124000                     | 2521693000            | 0                        | 0                           | 0                            | 12544057344    | 0                    | 0               | 0                 | 7757824        | 0                  | 0                            |
| 2022-12-22 18:08:57 | true            | false        | false             | false          | false              |               | running        |                  |                      | 168274799 | 161101715 | 8                      | -1                   | 0                        | 0                              | 0                               | 8            | 185774320000000    | 1303667000               | 1348816000                 | 2652484000        | 0                    | 0                       | 0                        | 8                | 185766390000000        | 1299718000                   | 1343042000                     | 2642760000            | 0                        | 0                           | 0                            | 12544057344    | 0                    | 0               | 0                 | 7880704        | 0                  | 0                            |
| 2022-12-22 18:09:07 | true            | false        | false             | false          | false              |               | running        |                  |                      | 170555804 | 161260220 | 8                      | -1                   | 0                        | 0                              | 0                               | 8            | 185853390000000    | 1347750000               | 1409181000                 | 2756931000        | 0                    | 0                       | 0                        | 8                | 185845420000000        | 1346494000                   | 1396848000                     | 2743342000            | 0                        | 0                           | 0                            | 12544057344    | 0                    | 0               | 0                 | 7176192        | 0                  | 0                            |
| 2022-12-22 18:09:17 | true            | false        | false             | false          | false              |               | running        |                  |                      | 170159628 | 162693710 | 8                      | -1                   | 0                        | 0                              | 0                               | 8            | 185932420000000    | 1392506000               | 1446422000                 | 2838928000        | 0                    | 0                       | 0                        | 8                | 185924480000000        | 1383817000                   | 1443256000                     | 2827074000            | 0                        | 0                           | 0                            | 12544057344    | 0                    | 0               | 0                 | 7323648        | 0                  | 0                            |
| 2022-12-22 18:09:27 | true            | false        | false             | false          | false              |               | running        |                  |                      | 168085132 | 162079882 | 8                      | -1                   | 0                        | 0                              | 0                               | 8            | 186011390000000    | 1433485000               | 1491347000                 | 2924833000        | 0                    | 0                       | 0                        | 8                | 186003470000000        | 1424233000                   | 1489479000                     | 2913712000            | 0                        | 0                           | 0                            | 12544057344    | 0                    | 0               | 0                 | 7180288        | 0                  | 0                            |
| 2022-12-22 18:09:37 | true            | false        | false             | false          | false              |               | running        |                  |                      | 171131804 | 164846262 | 8                      | -1                   | 0                        | 0                              | 0                               | 8            | 186090440000000    | 1488872000               | 1532001000                 | 3020873000        | 0                    | 0                       | 0                        | 8                | 186082480000000        | 1483453000                   | 1524608000                     | 3008061000            | 0                        | 0                           | 0                            | 12544057344    | 0                    | 0               | 0                 | 7204864        | 0                  | 0                            |
| 2022-12-22 18:09:47 | true            | false        | false             | false          | false              |               | running        |                  |                      | 246547711 | 205213169 | 8                      | -1                   | 0                        | 0                              | 0                               | 8            | 186170070000000    | 1536031000               | 1575130000                 | 3111162000        | 0                    | 0                       | 0                        | 8                | 186161850000000        | 1533463000                   | 1570638000                     | 3104101000            | 0                        | 0                           | 0                            | 12544057344    | 0                    | 0               | 0                 | 7200768        | 0                  | 0                            |
| 2022-12-22 18:09:57 | true            | false        | false             | false          | false              |               | running        |                  |                      | 171709466 | 164836548 | 8                      | -1                   | 0                        | 0                              | 0                               | 8            | 186248500000000    | 1575590000               | 1627553000                 | 3203144000        | 0                    | 0                       | 0                        | 8                | 186240540000000        | 1573723000                   | 1620283000                     | 3194006000            | 0                        | 0                           | 0                            | 12544057344    | 0                    | 0               | 0                 | 7241728        | 0                  | 0                            |
| 2022-12-22 18:10:07 | true            | false        | false             | false          | false              |               | running        |                  |                      | 169028179 | 161384887 | 8                      | -1                   | 0                        | 0                              | 0                               | 8            | 186327640000000    | 1616116000               | 1679493000                 | 3295610000        | 0                    | 0                       | 0                        | 8                | 186319650000000        | 1613259000                   | 1674805000                     | 3288065000            | 0                        | 0                           | 0                            | 12544057344    | 0                    | 0               | 0                 | 7507968        | 0                  | 0                            |
| 2022-12-22 18:10:17 | true            | false        | false             | false          | false              |               | running        |                  |                      | 170431544 | 162779377 | 8                      | -1                   | 0                        | 0                              | 0                               | 8            | 186406680000000    | 1667976000               | 1740596000                 | 3408573000        | 0                    | 0                       | 0                        | 8                | 186398720000000        | 1659381000                   | 1740596000                     | 3399978000            | 0                        | 0                           | 0                            | 12544057344    | 0                    | 0               | 0                 | 7860224        | 0                  | 0                            |
| 2022-12-22 18:10:27 | true            | false        | false             | false          | false              |               | running        |                  |                      | 179464091 | 169394757 | 8                      | -1                   | 0                        | 0                              | 0                               | 8            | 186485630000000    | 1741757000               | 1819497000                 | 3561255000        | 0                    | 0                       | 0                        | 8                | 186477650000000        | 1735767000                   | 1808013000                     | 3543781000            | 0                        | 0                           | 0                            | 12544057344    | 0                    | 0               | 0                 | 8175616        | 0                  | 0                            |
| 2022-12-22 18:10:37 | true            | false        | false             | false          | false              |               | running        |                  |                      | 178936845 | 166344095 | 8                      | -1                   | 0                        | 0                              | 0                               | 8            | 186564540000000    | 1789673000               | 1883867000                 | 3673541000        | 0                    | 0                       | 0                        | 8                | 186556550000000        | 1787416000                   | 1874381000                     | 3661798000            | 0                        | 0                           | 0                            | 12544057344    | 0                    | 0               | 0                 | 8331264        | 0                  | 0                            |

# Example of crash capture

> Automatic function

```
      @2022-12-16T04:17:02.324250972Z 2022/12/16 04:17:02 IP: 10.0.0.6
      v2022-12-16T04:17:02.324469972Z 2022/12/16 04:17:02 [DEBUG] memberlist: Initiating push/pull sync with:  10.0.0.6:7946
      =2022-12-16T04:17:02.325176930Z Member: be594d5ade2e 10.0.0.6
      =2022-12-16T04:17:02.325184055Z Member: 48b8e00607b2 10.0.0.7
      F2022-12-16T04:17:02.326104639Z 2022/12/16 04:17:02 nats connection ok
      @2022-12-16T04:17:02.326115514Z 2022/12/16 04:17:02 chaos enable
      m2022-12-16T04:17:18.253631340Z 2022/12/16 04:17:18 [DEBUG] memberlist: Stream connection from=10.0.0.8:40948
      m2022-12-16T04:17:21.766587550Z 2022/12/16 04:17:21 [DEBUG] memberlist: Stream connection from=10.0.0.6:54688
      G2022-12-16T04:17:22.327208800Z 2022/12/16 04:17:22 you can restart now
      �2022-12-16T04:17:26.174066386Z 2022/12/16 04:17:26 [DEBUG] memberlist: Initiating push/pull sync with: be594d5ade2e 10.0.0.6:7946
      m2022-12-16T04:17:33.256542333Z 2022/12/16 04:17:33 [DEBUG] memberlist: Stream connection from=10.0.0.8:48176
      m2022-12-16T04:17:36.770945293Z 2022/12/16 04:17:36 [DEBUG] memberlist: Stream connection from=10.0.0.6:34038
      �2022-12-16T04:17:41.178559545Z 2022/12/16 04:17:41 [DEBUG] memberlist: Initiating push/pull sync with: 575857e427da 10.0.0.8:7946
      m2022-12-16T04:17:48.258598174Z 2022/12/16 04:17:48 [DEBUG] memberlist: Stream connection from=10.0.0.8:40284
      �2022-12-16T04:17:56.183927386Z 2022/12/16 04:17:56 [DEBUG] memberlist: Initiating push/pull sync with: 575857e427da 10.0.0.8:7946
      m2022-12-16T04:18:03.260538041Z 2022/12/16 04:18:03 [DEBUG] memberlist: Stream connection from=10.0.0.8:40656
      w2022-12-16T04:18:03.524478167Z 2022/12/16 04:18:03 [DEBUG] memberlist: Failed UDP ping: be594d5ade2e (timeout reached)
      x2022-12-16T04:18:04.325072375Z 2022/12/16 04:18:04 [INFO] memberlist: Suspect be594d5ade2e has failed, no acks received
      w2022-12-16T04:18:05.523893459Z 2022/12/16 04:18:05 [DEBUG] memberlist: Failed UDP ping: be594d5ade2e (timeout reached)
      �2022-12-16T04:18:05.542607542Z 2022/12/16 04:18:05 [INFO] memberlist: Marking be594d5ade2e as failed, suspect timeout reached (1 peer confirmations)
      x2022-12-16T04:18:06.327539501Z 2022/12/16 04:18:06 [INFO] memberlist: Suspect be594d5ade2e has failed, no acks received
      �2022-12-16T04:18:11.186402295Z 2022/12/16 04:18:11 [DEBUG] memberlist: Initiating push/pull sync with: 575857e427da 10.0.0.8:7946
      m2022-12-16T04:18:18.262939049Z 2022/12/16 04:18:18 [DEBUG] memberlist: Stream connection from=10.0.0.8:34016
      @2022-12-16T04:19:38.808491836Z 2022/12/16 04:19:38 IP: 10.0.0.6
      v2022-12-16T04:19:38.808835169Z 2022/12/16 04:19:38 [DEBUG] memberlist: Initiating push/pull sync with:  10.0.0.6:7946
      =2022-12-16T04:19:38.809361419Z Member: be594d5ade2e 10.0.0.6
      =2022-12-16T04:19:38.809370169Z Member: 48b8e00607b2 10.0.0.7
      ~2022-12-16T04:19:40.810843628Z 2022/12/16 04:19:40 nats connection error: read tcp 10.0.0.7:39924->10.0.0.2:4222: i/o timeout
      r2022-12-16T04:19:40.810940878Z 2022/12/16 04:19:40 bug: messageSystem.Subscribe().error: nats: invalid connection
```