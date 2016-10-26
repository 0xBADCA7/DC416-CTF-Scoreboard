# IVScoreboard

A scoreboard application for capture the flag events.

In the future, the scoreboard application itself may include some (intentional) vulnerabilities to be exploited for bonus points.

![Scoreboard Screenshot](https://u.pomf.is/kwamcd.png)



## Running the application

The following instructions explain how to build, configure and run the application from source.

### Dependencies

#### SQLite3

The scoreboard uses SQLite to persist information about teams and the flags they have submitted. You will likely be able to install it via your operating system's package manager.

#### The Go Toolset

The easiest way to get the Go compiler and other tools is by installing the tool suite from the [official site](https://golang.org/dl/). Once installed, you should be able to run a command like the following and see the version of Go that you have installed.

```
$ go version
go version go1.7.3 darwin/amd64
```

You will probably also need to set the environment variable `GOROOT` to the path to where your base install of Go ended up in either your `.profile`, `.bashrc` or `.zshrc` file.

#### Third-Party Go Libraries

To install each of the third-party libraries required by the application, simply run the following command.

```
go get github.com/mattn/go-sqlite3
```

### Configuration

Once you have set up the dependencies for the project, the next step is to start configuring the server. First, clone the repository if you haven't already.

```
git clone https://github.com/DC416/IVScoreboard
cd IVScoreboard
```

Now, open `config/config.json`. The structure of the configuration for the service is as follows:

```javascript
{
  "ctfName": "Defcon Toronto Offline CTF",
  "bindAddress": "0.0.0.0:3000",
  "dbFile": "scoreboard.db",
  "htmlDir": "html",
  "flags": [
    {
      "id": 1,
      "secret": "secret1",
      "reward": 10
    },
    {
      "id": 2,
      "secret": "secret2",
      "reward": 20
    }
  ]
}
```

* `ctfName` is the name for your CTF event, and will appear as `<ctfName> Scoreboard` on the index page served by the application.
* `bindAddress` is the IP and port to bind the web server to.
* `dbFile` is the name of the file you'd like SQLite to create and use.
* `htmlDir` is the path to the `html` directory containing HTML templates.
* `flags` is an array of flags that the application needs to be able to reward contestants for finding.
  * `id` must be a unique integer identifier for the flag. It's easiest to just count from 1.
  * `secret` is the secret flag string that contestants will have to enter to earn points.
  * `reward` is the number of points to award teams for finding the flag.

### Building and Running

The server can be built very simply from the `IVScoreboard` directory by invoking

```
go build
```

You can run the server one of two ways. The simplest is to just run the newly compiled binary.

```
./IVScoreboard
```

This will run the web server and try to load a configuration from `config/config.json`, relative to the `IVScoreboard/` directory that you are running. Alternatively, you can specify another path to a configuration file of your choosing by setting the `CONFIG_FILE` environment variable to that path. For example:

```
CONFIG_FILE=./config.json ./IVScoreboard
```

You should see output like the following and be able to verify in your browser that going to `localhost:3000/` brings you to the scoreboard page.

```
./config/config.json
Listening on 0.0.0.0:3000
```


### Submitting flags using curl
```
curl --data "token=TOKEN&flag=FLAG" ctf.server.com:3000/submit
```
