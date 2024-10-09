package gen

//go:generate wget -O openapi.json https://api.artifactsmmo.com/openapi.json
//go:generate ogen -target=oas openapi.json
