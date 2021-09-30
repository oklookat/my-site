package routerica

import (
	"fmt"
	"net/http"
)

const (
	//
	Path1          = "/hello"
	Path1Response1 = "TPR_1_GET"
	Path1Response2 = "TPR_1_POST"
	Path1Response3 = "TPR_1_PUT"
	Path1Response4 = "TPR_1_DELETE"
	//
	Path2          = "/hello/big/world"
	Path2Response1 = "TPR_2_GET"
	Path2Response2 = "TPR_2_POST"
	Path2Response3 = "TPR_2_PUT"
	Path2Response4 = "TPR_2_DELETE"
	//
	Path3          = "/wow/very/big/world/with/long/paths/and/its/good/or/maybe/not"
	Path3Response1 = "TPR_3_GET"
	Path3Response2 = "TPR_3_POST"
	Path3Response3 = "TPR_3_PUT"
	Path3Response4 = "TPR_3_DELETE"
	//
	Path4 = "/route/with/{username}/and/{password}"
	Path4Params = "/route/with/12/and/84"
	Path4Response1 = "GET_username_12_password_84"
	Path4Response2 = "POST_username_12_password_84"
	Path4Response3 = "PUT_username_12_password_84"
	Path4Response4 = "DELETE_username_12_password_84"
	//
	Path5 = "/{article}/and/{tentacle}"
	Path5Params = "/77/and/1024"
	Path5Response1 = "GET_article_77_tentacle_1024"
	Path5Response2 = "POST_article_77_tentacle_1024"
	Path5Response3 = "PUT_article_77_tentacle_1024"
	Path5Response4 = "DELETE_article_77_tentacle_1024"
	//
	GroupPath1 = "/group"
	GroupSubPath1 = "/hello/big/and/small/world/im/too/lazy/to/write/tests/and/its/boring"
	GroupSubPath1Request = "/group/hello/big/and/small/world/im/too/lazy/to/write/tests/and/its/boring"
	GroupSubPath1Response1 = "TPG_1_GET"
	GroupSubPath1Response2 = "TPG_1_POST"
	GroupSubPath1Response3 = "TPG_1_PUT"
	GroupSubPath1Response4 = "TPG_1_DELETE"
	//
	GroupPath2 = "/group2"
	GroupSubPath2 = "/something/and/{username}/no/{id}"
	GroupSubPath2Request = "/group2/something/and/99/no/2048"
	GroupSubPath2Response1 = "GET_username_99_id_2048"
	GroupSubPath2Response2 = "POST_username_99_id_2048"
	GroupSubPath2Response3 = "PUT_username_99_id_2048"
	GroupSubPath2Response4 = "DELETE_username_99_id_2048"
	//
)

// TestingEndpoint - basic endpoint for test requests.
func TestingEndpoint(response http.ResponseWriter, request *http.Request) {
	switch request.URL.Path {
	case Path1:
		pathHelper(response, request, false, Path1Response1, Path1Response2, Path1Response3, Path1Response4)
		break
	case Path2:
		pathHelper(response, request, false, Path2Response1, Path2Response2, Path2Response3, Path2Response4)
		break
	case Path3:
		pathHelper(response, request, false, Path3Response1, Path3Response2, Path3Response3, Path3Response4)
		break
	case Path4Params:
		pathHelper(response, request, true, "username", "password", "", "")
		break
	case Path5Params:
		pathHelper(response, request, true, "article", "tentacle", "", "")
		break
	case GroupSubPath1Request:
		pathHelper(response, request, false, GroupSubPath1Response1, GroupSubPath1Response2, GroupSubPath1Response3, GroupSubPath1Response4)
		break
	case GroupSubPath2Request:
		pathHelper(response, request, true, "username", "id", "","")
		break
	}
}

// pathHelper - depending on request method send response (watch TestingEndpoint for example).
func pathHelper(response http.ResponseWriter, request *http.Request, paramsMode bool, TPR1 string, TPR2 string, TPR3 string, TPR4 string){
	if paramsMode {
		var params = GetParams(request)
		var p1 = params[TPR1]
		var p2 = params[TPR2]
		var formatted = fmt.Sprintf("%v_%v_%v_%v_%v", request.Method, TPR1, p1, TPR2, p2)
		response.WriteHeader(200)
		response.Write([]byte(formatted))
		return
	}
	switch request.Method {
	case http.MethodGet:
		response.WriteHeader(200)
		response.Write([]byte(TPR1))
		break
	case http.MethodPost:
		response.WriteHeader(200)
		response.Write([]byte(TPR2))
		break
	case http.MethodPut:
		response.WriteHeader(200)
		response.Write([]byte(TPR3))
		break
	case http.MethodDelete:
		response.WriteHeader(200)
		response.Write([]byte(TPR4))
		break
	}
}