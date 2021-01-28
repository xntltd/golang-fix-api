# golang-fix-api

Golang FIX API library for XNT Ltd

### Building and installation

To install golang-fix-api library, use ```go get```:
```
$ go get github.com/xntltd/golang-fix-api
```

### Staying up to date

To update to the latest version, use
```
$ go get -u github.com/xntltd/golang-fix-api
```

### Developing golang-fix-api

If you wish to work on golang-fix-api itself, you will first need Go installed and configured on your machine (version 1.14+ is preferred, but the minimum required version is 1.8).

Next, using Git, clone the repository via ```git clone``` command 
```
git clone https://github.com/xntltd/golang-fix-api.git
```

## Basic usage

```
func main() {
	f, _ := NewFixAPI("cfg/session.conf", 100)
	if err := f.Run(); err != nil {
		panic("Cant start initiator")
	}
	err := f.SecurityList("WWB1220_TRADE_UAT", "EXANTE_TRADE_UAT")
	print(err)
	for {
		print(<- f.RespChan)
	}
	//f.Stop()
}
```

### License

Released under the [GNU GPL License](https://github.com/xntltd/golang-fix-api/blob/main/LICENSE)