# gredisd

> GRedis server is simple version of original Redis server with limited functionality and intended for playing
with Go lang only.

[![License][License-Image]][License-Url] [![ReportCard][ReportCard-Image]][ReportCard-Url] [![Build Status][Travis-Image]][Travis-Url]

## Usage

### Command line

    Usage: gredisd [options]
    
    Server Options:
        -a, --addr <host>                Bind to host address (default: 0.0.0.0)
        -p, --port <port>                Use port for clients (default: 16379)
            --databases <count>          Set the number of databases. The default database is DB 0, you can select
                                         a different one on a per-connection basis using SELECT <dbid> where
                                         dbid is a number between 0 and 'databases'-1
            --trace_protocol <bool>      Trace low level read/write operations

    Authorization Options:
            --auth <token>               Authorization token required for connections

    Common Options:
        -h, --help                       Show this message
        -v, --version                    Show version

### Environment variables

    GREDIS_PORT

      Use port for clients

    GREDIS_HOST

      Bind to host address

## Docker

    docker build -t=gredisd .
    docker run -p 16379:16379 -it --rm --network=bridge --name gredisd gredisd

## Docker Compose

    docker-compose up

## Installation

You can install the GRedis server Docker image or build the server from source.

### Build

You can build the latest version of the server from the `master` branch.

You need [*Go*](http://golang.org/) version 1.6+ [installed](https://golang.org/doc/install) to build the GRedis server.

- Run `go version` to verify that you are running Go 1.7+. (Run `go help` for more guidance.)
- Clone the <https://github.com/valery-barysok/gredisd> repository.
- Run `go build` inside the `/valery-barysok/gredisd` directory. A successful build produces no messages and creates the server executable `gredisd` in the directory.
- Run `go test ./...` to run the unit regression tests.

## Running

To start the GRedis server with default settings (and no authentication), you can invoke the `gredisd` binary with no [command line options](#command-line-arguments) or [configuration file](#configuration-file).

```sh
> ./gredisd
2016/12/20 20:43:14 GRedis version 0.0.1
2016/12/20 20:43:14 Listening for client connections on 0.0.0.0:16379
```

The server is started and listening for client connections on port 16379 (the default) from all available interfaces. The logs are displayed to stdout as shown above in the server output.

### Clients

Technically you can use official clients for Redis with Gredis that implements compatible subset of commands with using the same RESP protocol

For playing you can use

    redis-cli -p 16379

and type `COMMANDS` to see list of all supported commands by GRedis

### Protocol

The GRedis server uses a [*RESP (REdis Serialization Protocol)*](https://redis.io/topics/protocol#request-response-model), so interacting with it can be as simple as using telnet as shown below.

```sh
> telnet localhost 16379
Trying localhost...
Connected to localhost.
Escape character is '^]'.
*1
$4
PING
+PONG
*2
$6
EXISTS
$3
key
:0
*3
$3
SET
$3
key
$5
value
+OK
*2
$3
GET
$3
key
$5
value
*2
$6
EXISTS
$3
key
:1
```

### Networking layer

A client connects to a GRedis server creating a TCP connection to the port 16379.
While RESP is technically non-TCP specific, in the context of GRedis the protocol is only used with TCP connections.

### Request-Response model

GRedis accepts commands composed of different arguments. Once a command is received, it is processed and a reply is sent back to the client.

## Securing GRedis

### Authentication

The GRedis server supports token authentication.

    gredisd --auth S3Cr3t

Clients after connection has to execute `AUTH` command with required token before any other command.

## Supported commands

### Basic Commands

##### [**AUTH password**](https://redis.io/commands/auth)

  Request for authentication in a password-protected GRedis server. GRedis can be instructed to
  require a password before allowing clients to execute commands. This is done using the pass
  option.

  If password matches the password configured for GRedis, the server replies with the OK status
  code and starts accepting commands. Otherwise, an error is returned and the clients needs to
  try a new password.

##### [**SELECT index**](https://redis.io/commands/select)

  Select the DB with having the specified zero-based numeric index. New connections always use DB 0.

##### [**ECHO message**](https://redis.io/commands/echo)

  Returns message.

##### [**PING [message]**](https://redis.io/commands/ping)

  Returns `PONG` if no argument is provided, otherwise return a copy of the argument as a bulk.

##### [**SHUTDOWN**](https://redis.io/commands/shutdown)

  The command behavior is the following:
     Stop all the clients.
     Quit the server.

  > Note: it is simple version of redis `SHUTDOWN` command

##### [**COMMANDS**](https://redis.io/commands/command)

  Returns Array of all supported commands.

  > Note: it has custom name due to incompatibility with original command from redis that used by redis-cli

##### [**KEYS pattern**](https://redis.io/commands/keys)

  Returns all keys matching **regexp** pattern.

  > Note: original matching pattern from redis for `KEYS` command differ from regexp pattern used by this one

##### [**EXISTS key [key ...]**](https://redis.io/commands/exists)

  Returns if key exists.

  It is possible to specify multiple keys instead of a single one. In such a case, it returns the total
  number of keys existing.

  The user should be aware that if the same existing key is mentioned in the arguments multiple times,
  it will be counted multiple times. So if `somekey` exists, `EXISTS somekey somekey` will return 2.

##### [**EXPIRE key seconds**](https://redis.io/commands/expire)

  Expire sets a timeout on key. After the timeout has expired, the key will automatically be deleted.

  - 1 if the timeout was set.
  - 0 if key does not exist or the timeout could not be set.

### Key Value Commands

##### [**SET key value [EX seconds] [PX milliseconds] [NX|XX]**](https://redis.io/commands/set)

  Set key to hold the string value. If key already holds a value, it is overwritten, regardless of its type.
  Any previous time to live associated with the key is discarded on successful `SET` operation.

##### [**GET key**](https://redis.io/commands/get)

  Get the value of key. If the key does not exist the special value nil is returned. An error is returned
  if the value stored at key is not a string, because `GET` only handles string values.

##### [**DEL key [key ...]**](https://redis.io/commands/del)

  Removes the specified keys. A key is ignored if it does not exist.

### Key Value List Commands

##### [**LPUSH key value [value ...]**](https://redis.io/commands/lpush)

  Insert all the specified values at the head of the list stored at key. If key does not exist, it is
  created as empty list before performing the push operations. When key holds a value that is not a list,
  an error is returned.

  It is possible to push multiple elements using a single command call just specifying multiple arguments
  at the end of the command. Elements are inserted one after the other to the head of the list, from the
  leftmost element to the rightmost element. So for instance the command `LPUSH mylist a b c` will result
  into a list containing `c` as first element, `b` as second element and `a` as third element.

##### [**RPUSH key value [value ...]**](https://redis.io/commands/rpush)

  Insert all the specified values at the tail of the list stored at key. If key does not exist, it is
  created as empty list before performing the push operation. When key holds a value that is not a list,
  an error is returned.

  It is possible to push multiple elements using a single command call just specifying multiple arguments
  at the end of the command. Elements are inserted one after the other to the tail of the list, from the
  leftmost element to the rightmost element. So for instance the command `RPUSH mylist a b c` will result
  into a list containing `a` as first element, `b` as second element and `c` as third element.

##### [**LPOP key**](https://redis.io/commands/lpop)

  Removes and returns the first element of the list stored at key.

##### [**RPOP key**](https://redis.io/commands/rpop)

  Removes and returns the last element of the list stored at key.

##### [**LLEN key**](https://redis.io/commands/llen)

  Returns the length of the list stored at key. If key does not exist, it is interpreted as an empty list
  and `0` is returned. An error is returned when the value stored at key is not a list.

##### [**LINSERT key BEFORE|AFTER pivot value**](https://redis.io/commands/linsert)

  Inserts value in the list stored at key either before or after the reference value pivot.

  When key does not exist, it is considered an empty list and no operation is performed.

  An error is returned when key exists but does not hold a list value.

##### [**LINDEX key index**](https://redis.io/commands/lindex)

  Returns the element at index index in the list stored at key. The index is zero-based, so 0 means the
  first element, 1 the second element and so on. Negative indices can be used to designate elements
  starting at the tail of the list. Here, -1 means the last element, -2 means the penultimate and so
  forth.

  When the value at key is not a list, an error is returned.

##### [**LRANGE key start stop**](https://redis.io/commands/lrange)

  Returns the specified elements of the list stored at key. The offsets start and stop are zero-based
  indexes, with 0 being the first element of the list (the head of the list), 1 being the next element
  and so on.

  These offsets can also be negative numbers indicating offsets starting at the end of the list.
  For example, -1 is the last element of the list, -2 the penultimate, and so on.

### Key Value Dict Commands

##### [**HSET key field value**](https://redis.io/commands/hset)

  Sets field in the hash stored at key to value. If key does not exist, a new key holding a hash is
  created. If field already exists in the hash, it is overwritten.

##### [**HGET key field**](https://redis.io/commands/hget)

  Returns the value associated with field in the hash stored at key.

##### [**HDEL key field [field ...]**](https://redis.io/commands/hdel)

  Removes the specified fields from the hash stored at key. Specified fields that do not exist
  within this hash are ignored. If key does not exist, it is treated as an empty hash and this
  command returns 0.

##### [**HLEN key**](https://redis.io/commands/hlen)

  Returns the number of fields contained in the hash stored at key.

##### [**HEXISTS key field**](https://redis.io/commands/hexists)

  Returns if field is an existing field in the hash stored at key.
  
[License-Url]: http://opensource.org/licenses/Apache-2.0
[License-Image]: https://img.shields.io/badge/License-Apache%202.0-blue.svg?style=flat-square
[ReportCard-Url]: http://goreportcard.com/report/valery-barysok/gredisd
[ReportCard-Image]: http://goreportcard.com/badge/github.com/valery-barysok/gredisd?style=flat-square
[Travis-Image]: https://img.shields.io/travis/valery-barysok/gredisd/master.svg?style=flat-square
[Travis-Url]: https://travis-ci.org/valery-barysok/gredisd