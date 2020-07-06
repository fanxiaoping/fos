package router

import (
	"github.com/fanxiaoping/fos/service/v1"
	"github.com/fanxiaoping/fos/util/config"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"net/http"
	"path/filepath"
)

//	RegisterServer 注册grpc
func RegisterServer(server *grpc.Server) {
}

//	RegisterHandlerFromEndpoint 注册geteway
func RegisterHandlerFromEndpoint(ctx context.Context, gwmux *runtime.ServeMux, port string, dopts []grpc.DialOption) {

}

// RegisterMux 注册form提交路由
func RegisterMux() *mux.Router {
	r := mux.NewRouter()

	//文件下载
	r.Handle("/file/down/{id}", &v1.DownService{})
	//上传文件
	r.Handle("/file/upload", &v1.UploadService{})
	//上传文件测试
	r.Handle("/file/test/upload", http.StripPrefix("/file/test/upload", http.FileServer(http.Dir(filepath.Join(config.GetCurrentDirectory(), "static","test")))))

	return r
}
