module github.com/poy/assistant

go 1.20

require (
	cloud.google.com/go/compute v1.19.3 // indirect
	cloud.google.com/go/compute/metadata v0.2.3 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/go-react v0.0.0-20230606162745-b3508c8c73ee
	github.com/poy/go-dependency-injection v0.0.0
	golang.org/x/net v0.10.0 // indirect
	golang.org/x/oauth2 v0.8.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
)

replace github.com/google/go-react => ../go-react

replace github.com/poy/go-dependency-injection => ../go-dependency-injection
