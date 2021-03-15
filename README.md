## safe-trash

"safe-trash" is a safe way of running the `rm` command. 
Deleted files are moved to `~/.safe-trash` folder from the current folder.

### Usage

    -l, --list      List the files under the current directory.
    <FILE_NAME>     The filenames to be deleted. Multiple filename is allowed.
    -h, --help      Displays this help message.

### How to run? 

Ensure that your Go variables set properly.
If your go env variables are set properly, you can run `sh install_go.sh` or `bash install_go.sh` to use safe-trash as;

```shell script
$ sh install_go.sh
$ safe-trash server.py
```

Or, you can;
```shell script
$ go build -o $GOPATH/bin/safe-trash
$ safe-trash index.js
```

### Inspiration
A simple replica of [trash-cli](https://github.com/sindresorhus/trash-cli) written in Golang.