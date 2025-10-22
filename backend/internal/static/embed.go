package static

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed dist
var StaticFiles embed.FS

// GetFileSystem 返回嵌入的静态文件系统
func GetFileSystem() http.FileSystem {
	// 获取dist子目录
	fsys, err := fs.Sub(StaticFiles, "dist")
	if err != nil {
		panic(err)
	}
	return http.FS(fsys)
}