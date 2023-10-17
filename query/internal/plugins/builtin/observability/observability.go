package clickhouse

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	ch "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/DataObserve/datav/query/internal/plugins/builtin/observability/api"
	pluginUtils "github.com/DataObserve/datav/query/internal/plugins/utils"
	"github.com/DataObserve/datav/query/pkg/colorlog"
	"github.com/DataObserve/datav/query/pkg/models"
	"github.com/gin-gonic/gin"
)

/* Query plugin for clickhouse database */

var datasourceName = "observability"

type ObservabilityPlugin struct{}

var conns = make(map[int64]ch.Conn)
var connsLock = &sync.Mutex{}

func (p *ObservabilityPlugin) Query(c *gin.Context, ds *models.Datasource) models.PluginResult {
	query := c.Query("query")
	if query == api.TestDatasourceAPI {
		return pluginUtils.TestClickhouseDatasource(c)
	}

	conn, ok := conns[ds.Id]
	if !ok {
		var err error
		conn, err = pluginUtils.ConnectToClickhouse(ds.URL, ds.Data["database"], ds.Data["username"], ds.Data["password"])
		if err != nil {
			colorlog.RootLogger.Warn("connect to clickhouse error:", err, "ds_id", ds.Id, "url", ds.URL)
			return models.GenPluginResult(models.PluginStatusError, err.Error(), nil)

		}
		connsLock.Lock()
		conns[ds.Id] = conn
		connsLock.Unlock()
	}
	route, ok := api.APIRoutes[query]
	if ok {
		paramStr := c.Query("params")
		params := make(map[string]interface{})
		err := json.Unmarshal([]byte(paramStr), &params)
		if err != nil {
			return models.GenPluginResult(models.PluginStatusError, fmt.Sprintf("decode params error: %s", err.Error()), nil)
		}

		start := time.Now()
		res := route(c, ds, conn, params)

		colorlog.RootLogger.Info("Excecute observability query api", "query", query, "time", time.Since(start).String())
		return models.PluginResult{
			Status: models.PluginStatusSuccess,
			Error:  "",
			Data:   res,
		}
	} else {
		return models.GenPluginResult(models.PluginStatusError, "api not found", nil)
	}
}

func init() {
	// register datasource
	models.RegisterPlugin(datasourceName, &ObservabilityPlugin{})
}
