# softrecevier
golang package for receiving large file

## receive
```go
safeReceive := softrecevier.New()
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
