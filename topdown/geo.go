package topdown

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net"

	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/topdown/builtins"

	"github.com/oschwald/geoip2-golang"
	"github.com/posener/gitfs"
)

const (
	databaseRepo = "github.com/yashtewari/geoip"
	databasePath = "geolite2-city.mmdb"
)

var (
	memdb *geoip2.Reader
)

func builtinGeoFromIP(inp ast.Value) (ast.Value, error) {
	ip, err := builtins.StringOperand(inp, 1)
	if err != nil {
		return nil, err
	}

	if memdb == nil {
		if err := initMemdb(); err != nil {
			return nil, err
		}
	}

	geo, err := memdb.City(net.ParseIP(string(ip)))
	if err != nil {
		return nil, err
	}

	j, err := json.Marshal(geo)
	if err != nil {
		return nil, err
	}

	var x interface{}
	if err := json.Unmarshal(j, &x); err != nil {
		return nil, err
	}

	return ast.InterfaceToValue(x)
}

func initMemdb() error {
	ctx := context.Background()

	fs, err := gitfs.New(ctx, databaseRepo)
	if err != nil {
		return err
	}

	f, err := fs.Open(databasePath)
	if err != nil {
		return err
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	memdb, err = geoip2.FromBytes(b)
	if err != nil {
		return err
	}

	return nil
}

func init() {
	RegisterFunctionalBuiltin1(ast.GeoFromIP.Name, builtinGeoFromIP)
}
