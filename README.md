# GoBencode
Library used for reading and writing bencode, usually, .torrent files

Decoder results can be printed out in JSON-like style

# Examples
Take a look at the [examples](https://github.com/MrPythoneer/GoBencode/blob/main/example/main.go)

## Decoder output
>["hello", "world", 12345, {"type": 1, "animal": "cat"}]

## Encoder output
>l5:hello5:worldi12345ed6:animal3:cat4:typei1eee

## Go example
```go
item := DecodeFile(filepath)
fmt.Println(item)               // ["hello", "world", 12345, {"type": 1, "animal": "cat"}]

list := item.Value.([]BencodeItem)
fmt.Println(list[0])            // "hello"

dict := list[3].Value.(map[string]BencodeItem)
fmt.Println(dict["animal"])     // "cat"
```
# BencodeCLI
Take a look at the CLI tool based on this library. [Link](https://github.com/MrPythoneer/BencodeCLI)
