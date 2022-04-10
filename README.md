# Vend Test

## Application Structure 

(Ben Johnson's article)[https://medium.com/@benbjohnson/standard-package-layout-7cdbc8391fc1]

```
├── README.md
├── cmd
│   └── http_server
│       └── main.go // Binaries
├── go.mod
├── go.sum
├── http
│   └── // Contains http related stuff.
├── inmemory
│   └── // InMemory Database
└── product.go // Root is for domain objects.
```
