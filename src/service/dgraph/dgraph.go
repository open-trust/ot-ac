package dgraph

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/open-trust/ot-ac/src/conf"
	"github.com/open-trust/ot-ac/src/logging"
	"github.com/open-trust/ot-ac/src/util"

	"github.com/dgraph-io/dgo/v200"
	"github.com/dgraph-io/dgo/v200/protos/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/encoding/gzip"
)

func init() {
	util.DigProvide(NewDgraph)
}

// Nquads ...
type Nquads struct {
	ID    string
	Type  string
	UKkey string
	UKval string
	KV    map[string]interface{}
}

func scalarVal(v interface{}) string {
	switch val := v.(type) {
	case string:
		l := len(val)
		switch {
		case val == "*":
			return val
		case l > 1 && val[:1] == "<":
			return val
		case l > 2 && val[:2] == "_:":
			return val
		case l > 4 && (val[:4] == "uid(" || val[:4] == "val("):
			return val
		}
		return strconv.Quote(val)
	case time.Time:
		return fmt.Sprintf("\"%s\"", val.UTC().Format(time.RFC3339))
	default:
		return fmt.Sprintf("\"%v\"", val)
	}
}

// WithLang ...
type WithLang struct {
	V string
	L string
}

func (l WithLang) String() string {
	return fmt.Sprintf("%s@%s", strconv.Quote(l.V), l.L)
}

// WithFacets ...
type WithFacets struct {
	V  interface{}
	KV map[string]interface{}
}

func (fs WithFacets) String() string {
	if len(fs.KV) == 0 {
		return scalarVal(fs.V)
	}
	kv := make([]string, 0, len(fs.KV))
	for k, v := range fs.KV {
		switch val := v.(type) {
		case string:
			kv = append(kv, fmt.Sprintf("%s=%s", k, strconv.Quote(val)))
		case time.Time:
			kv = append(kv, fmt.Sprintf("%s=%s", k, val.UTC().Format(time.RFC3339)))
		default:
			kv = append(kv, fmt.Sprintf("%s=%v", k, v))
		}
	}
	return fmt.Sprintf("%s (%s)", scalarVal(fs.V), strings.Join(kv, ", "))
}

// Bytes ...
func (ns Nquads) Bytes() ([]byte, error) {
	if ns.ID == "" {
		return nil, errors.New("ID are required")
	}

	l := len(ns.ID)
	if ns.ID[:1] != "<" && (l > 2 && ns.ID[:2] != "_:") && (l > 4 && ns.ID[:4] != "uid(") {
		ns.ID = fmt.Sprintf("<%s>", ns.ID)
	}

	b := new(bytes.Buffer)
	if ns.Type != "" {
		if err := writeNquad(b, ns.ID, "dgraph.type", ns.Type); err != nil {
			return nil, err
		}
	}

	for k, v := range ns.KV {
		if err := writeNquad(b, ns.ID, k, v); err != nil {
			return nil, err
		}
	}
	return b.Bytes(), nil
}

func writeNquad(w io.Writer, subject, predicate string, object interface{}) error {
	var err error
	if predicate != "*" {
		predicate = fmt.Sprintf("<%s>", predicate)
	}
	switch val := object.(type) {
	case string, bool, int, int64, float64, time.Time:
		_, err = fmt.Fprintf(w, "%s %s %s .\n", subject, predicate, scalarVal(val))
	case WithLang:
		_, err = fmt.Fprintf(w, "%s %s %s .\n", subject, predicate, val.String())
	case WithFacets:
		_, err = fmt.Fprintf(w, "%s %s %s .\n", subject, predicate, val.String())
	case []string:
		for _, v := range val {
			_, err = fmt.Fprintf(w, "%s %s %s .\n", subject, predicate, scalarVal(v))
		}
	case []int:
		for _, v := range val {
			_, err = fmt.Fprintf(w, "%s %s \"%v\" .\n", subject, predicate, v)
		}
	case []int64:
		for _, v := range val {
			_, err = fmt.Fprintf(w, "%s %s \"%v\" .\n", subject, predicate, v)
		}
	case []float64:
		for _, v := range val {
			_, err = fmt.Fprintf(w, "%s %s \"%v\" .\n", subject, predicate, v)
		}
	case []time.Time:
		for _, v := range val {
			_, err = fmt.Fprintf(w, "%s %s \"%s\" .\n", subject, predicate, v.Format(time.RFC3339))
		}
	case []WithLang:
		for _, v := range val {
			_, err = fmt.Fprintf(w, "%s %s %s .\n", subject, predicate, v.String())
		}
	case []WithFacets:
		for _, v := range val {
			_, err = fmt.Fprintf(w, "%s %s %s .\n", subject, predicate, v.String())
		}
	default:
		err = fmt.Errorf("invalid value: %#v", object)
	}
	return err
}

// Dgraph ...
type Dgraph struct {
	*dgo.Dgraph
	dc api.DgraphClient
}

// NewDgraph ...
func NewDgraph() (*Dgraph, error) {
	opts := []grpc.DialOption{
		grpc.WithTimeout(time.Second * 5),
		grpc.WithDefaultCallOptions(grpc.UseCompressor(gzip.Name)),
	}
	if conf.Config.Dgraph.Insecure {
		opts = append(opts, grpc.WithInsecure())
	}
	conn, err := grpc.Dial(conf.Config.Dgraph.GRPCEndpoint, opts...)

	if err != nil {
		return nil, err
	}

	dc := api.NewDgraphClient(conn)
	dg := dgo.NewDgraphClient(dc)
	return &Dgraph{Dgraph: dg, dc: dc}, nil
}

// CheckHealth ...
func (dg *Dgraph) CheckHealth(ctx context.Context) (interface{}, error) {
	_, err := dg.dc.CheckVersion(ctx, &api.Check{})
	if err != nil {
		return nil, err
	}
	return map[string]string{"status": "OK"}, nil
}

// Query ...
func (dg *Dgraph) Query(ctx context.Context, query string, vars map[string]string, out interface{}) error {
	txn := dg.NewReadOnlyTxn()
	resp, err := loggingDgraph(ctx, func() (*api.Response, error) {
		return txn.QueryWithVars(ctx, query, vars)
	})
	if err == nil && out != nil && len(resp.Json) > 0 {
		err = json.Unmarshal(resp.Json, out)
	}
	return err
}

// QueryBestEffort ...
func (dg *Dgraph) QueryBestEffort(ctx context.Context, query string, vars map[string]string, out interface{}) error {
	txn := dg.NewReadOnlyTxn().BestEffort()
	resp, err := loggingDgraph(ctx, func() (*api.Response, error) {
		return txn.QueryWithVars(ctx, query, vars)
	})
	if err == nil && out != nil && len(resp.Json) > 0 {
		err = json.Unmarshal(resp.Json, out)
	}
	return err
}

// Do ...
func (dg *Dgraph) Do(ctx context.Context, query string, vars map[string]string, out interface{}, mus ...*api.Mutation) error {
	if len(mus) == 0 {
		return dg.Query(ctx, query, vars, out)
	}

	txn := dg.NewTxn()
	defer txn.Discard(ctx)

	req := &api.Request{
		Query:     query,
		Vars:      vars,
		Mutations: mus,
		CommitNow: true,
	}
	resp, err := loggingDgraph(ctx, func() (*api.Response, error) {
		return txn.Do(ctx, req)
	})
	if err == nil && out != nil && len(resp.Json) > 0 {
		err = json.Unmarshal(resp.Json, out)
	}
	return err
}

// DoRaw ...
func (dg *Dgraph) DoRaw(ctx context.Context, query string, vars map[string]string, mus ...*api.Mutation) (*api.Response, error) {
	if len(mus) == 0 {
		txn := dg.NewReadOnlyTxn().BestEffort()
		return txn.QueryWithVars(ctx, query, vars)
	}

	txn := dg.NewTxn()
	defer txn.Discard(ctx)

	req := &api.Request{
		Query:     query,
		Vars:      vars,
		Mutations: mus,
		CommitNow: true,
	}

	return loggingDgraph(ctx, func() (*api.Response, error) {
		return txn.Do(ctx, req)
	})
}

func loggingDgraph(ctx context.Context, fn func() (*api.Response, error)) (*api.Response, error) {
	startT := time.Now()
	resp, err := fn()
	end := time.Now().Sub(startT) / 1000000
	if end > 10 {
		logging.SetKV(ctx, "dgraphClientLatency", end)
		if resp != nil {
			logging.SetKV(ctx, "dgraphInternalLatency", resp.Latency.TotalNs/1000000)
		}
	}
	return resp, err
}
