package v1

import (
	"encoding/binary"
	"net/http"
	"strconv"

	. "github.com/chideat/pcc/pig/modules/config"
	. "github.com/chideat/pcc/pig/modules/id-gen"
	"github.com/gin-gonic/gin"
)

var gen *IDGen

const (
	AUTH_TOKEN = "H48K1VlJyqOcLqGmzEwSuuorW4qmHlQuDfY9bLSRKSw="
)

func init() {
	gen = NewIDGen(Config.Cache.Default)
}

func Auth(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")
	if token != AUTH_TOKEN {
		c.JSON(http.StatusNonAuthoritativeInfo, "非法请求")
		c.Abort()
		return
	}
}

// Route: /id/:type_id
// Method: GET
func Next(c *gin.Context) {
	typ, err := strconv.ParseUint(c.Params.ByName("type_id"), 10, 8)
	if err != nil {
		c.String(http.StatusNotFound, "不支持的类型")
		return
	}
	id, err := gen.Next(uint8(typ))
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	} else {
		buf := make([]byte, 8)
		binary.LittleEndian.PutUint64(buf, id)
		c.Data(http.StatusOK, "application/octet-stream", buf)
	}
}
