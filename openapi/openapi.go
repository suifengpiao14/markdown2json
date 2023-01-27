package parseapi

import "github.com/suifengpiao14/markdown2json/parsemarkdown"

const (
	OPENAPI_INFO_TITLE       = "openapi.info.title"
	OPENAPI_INFO_DESCRIPTION = "openapi.info.description"
	OPENAPI_INFO_VERSION     = "openapi.info.version"
	OPENAPI_SERVERS_ARR_URL  = "openapi.servers[].url"
	OPENAPI_PATHS_METHOD     = "openapi.paths.method"
)

func ToOpenApi(records parsemarkdown.Records) (openapi string, err error) {

	return "", nil
}
