[![Go Report Card](https://goreportcard.com/badge/github.com/joeke80215/safetyrecevier)](https://goreportcard.com/report/github.com/joeke80215/safetyrecevier)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/joeke80215/safetyrecevier/blob/master/LICENSE)
# safetyrecevier

golang package for receiving large file

## install
```
go get -u github.com/joeke80215/safetyrecevier
```

## receive
```go
maxMemoeryBufSize := 4 << 20
safeReceive := softrecevier.New(maxMemoeryBufSize)
defer safeReceive.CloseReceive()

for {
    chunk := make([]byte, 3) // set any chunk size
    n, err := exampleReader.Read(chunk)
    if err := safeReceive.Append(chunk,n);err != nil {
        if err == io.EOF {
            break
        }
    	// handle error
    }
}
```

## read
safeReceive implement io.ReadSeeker interface    
example
```go
b, _ := ioutil.ReadAll(safeReceive)
.
.
.
```
