# go get -d github.com/99designs/gqlgen@v0.17.36
GOTOOLCHAIN=go1.22.0 go run github.com/99designs/gqlgen@0b9bd5e generate --config ./user.gqlgen.yml

# go get -d github.com/99designs/gqlgen@v0.17.36
GOTOOLCHAIN=go1.22.0 go run github.com/99designs/gqlgen@0b9bd5e generate --config ./admin.gqlgen.yml