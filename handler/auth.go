package handler

// token 验证拦截器
import(
	"net/http"
)
// HTTPInterceptor：http 请求拦截器
func HTTPInterceptor(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request){
			r.ParseForm();
			username := r.Form.Get("username")
			token := r.Form.Get("token")

			if len(username)<3 || !IsTokenValid(token){
				w.WriteHeader(http.StatusForbidden) // 403
				return
			}
			h(w,r)
		})
}