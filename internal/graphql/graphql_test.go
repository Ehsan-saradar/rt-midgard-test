package graphql

import (
	"context"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	"gitlab.com/thorchain/midgard/internal/graphql/models"
	"gitlab.com/thorchain/midgard/internal/timeseries/stat"
	. "gopkg.in/check.v1"
)

func TestPoolByID(t *testing.T) {
	poolStakesLookup = func(poolID string, w stat.Window) (stat.PoolStakes, error) {
		if poolID != "test-asset" {
			t.Errorf("lookup for pool %q, want test-asset", poolID)
		}
		if !w.Start.IsZero() || !w.End.IsZero() {
			t.Errorf("lookup with time constraints %+v", w)
		}

		return stat.PoolStakes{
			AssetE8Total: 1,
			RuneE8Total:  2,
			UnitsTotal:   3,
		}, nil
	}

	req := httptest.NewRequest("GET", `/?query={pool(poolId:"test-asset"){asset%20poolStakedTotal%20runeStakedTotal%20poolUnits}}`, nil)
	resp := httptest.NewRecorder()
	Server.ServeHTTP(resp, req)
	if resp.Code != 200 {
		t.Fatalf("HTTP %d: %s", resp.Code, resp.Body)
	}

	const want = `{"data":{"pool":{"asset":"test-asset","poolStakedTotal":1,"runeStakedTotal":2,"poolUnits":3}}}`
	if got := resp.Body.String(); got != want {
		t.Errorf("got %q\nwant %q", got, want)
	}
}

type GraphQLSuite struct{}

var _ = Suite(&GraphQLSuite{})

type Stub struct {
	QueryResolver struct {
		Pool        func(ctx context.Context, poolID string) (*models.Pool, error)
		Pools       func(ctx context.Context, orderBy *models.PoolOrderAttribute, limit *int) ([]*models.Pool, error)
		PoolHistory func(ctx context.Context, from *int, until *int, interval *models.Interval, poolID *string) (*models.PoolHistory, error)
	}
}
type StubQuery struct{ *Stub }
type StubPoolHistory struct{ *Stub }

func (r *Stub) Query() QueryResolver {
	return &StubQuery{r}
}

func (r *Stub) PoolHistory() PoolHistoryResolver {
	return &StubPoolHistory{r}
}

func (r *StubQuery) Pools(ctx context.Context, orderBy *models.PoolOrderAttribute, limit *int) ([]*models.Pool, error) {
	return r.Stub.QueryResolver.Pools(ctx, orderBy, limit)
}

func (r *StubQuery) Pool(ctx context.Context, poolID string) (*models.Pool, error) {
	return r.Stub.QueryResolver.Pool(ctx, poolID)
}

func (r *StubQuery) PoolHistory(ctx context.Context, from *int, until *int, interval *models.Interval, poolID *string) (*models.PoolHistory, error) {
	return r.Stub.QueryResolver.PoolHistory(ctx, from, until, interval, poolID)
}

func (r *StubPoolHistory) Swaps(ctx context.Context, obj *models.PoolHistory) (*models.PoolSwaps, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *StubPoolHistory) Fees(ctx context.Context, obj *models.PoolHistory) (*models.PoolFees, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *StubPoolHistory) Slippage(ctx context.Context, obj *models.PoolHistory) (*models.PoolSlippage, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *GraphQLSuite) TestGraphQL(c *C) {
	resolvers := &Stub{}
	handler := handler.NewDefaultServer(NewExecutableSchema(Config{
		Resolvers:  resolvers,
		Directives: DirectiveRoot{},
		Complexity: ComplexityRoot{},
	}))
	client := client.New(handler)

	poolTest := &models.Pool{
		Asset:            "BNB.BNB",
		Status:           "enable",
		Price:            1,
		AssetStakedTotal: 2,
		RuneStakedTotal:  3,
		PoolStakedTotal:  3,
		AssetDepth:       4,
		RuneDepth:        5,
		PoolDepth:        6,
		PoolUnits:        7,
		CurrentAssetROI:  -0.1,
		CurrentRuneROI:   -0.2,
	}

	resolvers.QueryResolver.Pool = func(ctx context.Context, poolID string) (*models.Pool, error) {
		return poolTest, nil
	}
	resolvers.QueryResolver.Pools = func(ctx context.Context, orderBy *models.PoolOrderAttribute, limit *int) ([]*models.Pool, error) {
		return []*models.Pool{poolTest}, nil
	}

	query := `
     query {
        pool(poolId: "BNB.BNB") {
			asset
			status
			price
			assetStakedTotal
			runeStakedTotal
			poolStakedTotal
			assetDepth
			runeDepth
			poolDepth
			poolUnits
			currentAssetROI
			currentRuneROI
		}
		pools(orderBy: DEPTH, limit: 2) {
			asset
			status
			price
			assetStakedTotal
			runeStakedTotal
			poolStakedTotal
			assetDepth
			runeDepth
			poolDepth
			poolUnits
			currentAssetROI
			currentRuneROI
		}
    }`

	var resp struct {
		models.Pool
		Pools []models.Pool
	}

	err := client.Post(query, &resp)
	c.Assert(err, IsNil)
	c.Assert(poolTest, Equals, &resp.Pool)
	c.Assert(poolTest, Equals, &resp.Pools[0])
}
