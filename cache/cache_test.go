package cache_test

import (
    "context"
    "fmt"
	"os"
    "runtime"
    "testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "golang.org/x/sync/errgroup"

    "github.com/xmlking/toolkit/cache"
)

/********************
Models
********************/
type ServiceCache struct {
	Enabled bool
	Size    int
}

type Product struct {
	id    string
	name  string
	price float32
}

func fakeGetProductByID(productID string) (*Product, error) {
	time.Sleep(200 * time.Millisecond)
	return &Product{
		id:    productID,
		name:  fmt.Sprintf("prod %s name", productID),
		price: 32.21,
	}, nil
}

/********************
CatalogService
********************/
type CatalogService interface {
	GetProductByID(productID string) (prod *Product, err error)
}

// Cashed CatalogService
type cachedCatalogService struct {
	productCache cache.Cache
}

func (c *cachedCatalogService) GetProductByID(productID string) (*Product, error) {
	v, err := c.productCache.GetOrSet(productID, func() (interface{}, error) { return fakeGetProductByID(productID) })
	if err != nil {
		return nil, err
	}
	return v.(*Product), nil
}

// Raw CatalogService
type rawCatalogService struct {
}

func (d *rawCatalogService) GetProductByID(productID string) (*Product, error) {
	return fakeGetProductByID(productID)
}

func NewCatalogService(cacheConf ServiceCache) CatalogService {
	if cacheConf.Enabled {
		return &cachedCatalogService{cache.NewCache(cacheConf.Size)}
	} else {
		return &rawCatalogService{}
	}
}

/********************
Unit Tests
********************/
func TestMain(m *testing.M) {
	closer := func() {
		log.Info().Msg("closing the resource")
	}

	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	zerolog.TimeFieldFormat = time.RFC3339
	log.Debug().Msg("Logger set to debug")

	code := m.Run()
	closer() // cleanup
	os.Exit(code)
}

func TestCatalogService_NoCache(t *testing.T) {
	log.Info().Msg("called TestCatalogService_NoCache")
	catSrv := NewCatalogService(ServiceCache{Enabled: false})
	gotProduct, err := catSrv.GetProductByID("abc")

	assert.NoError(t, err)
	assert.Equal(t, "abc", gotProduct.id)
}

func TestCatalogService_Cache(t *testing.T) {
	log.Info().Msg("called TestCatalogService_Cache")
	catSrv := NewCatalogService(ServiceCache{Enabled: true, Size: 5})
	gotProduct, err := catSrv.GetProductByID("abc")
	gotProduct2, err2 := catSrv.GetProductByID("abc")

	assert.NoError(t, err)
	assert.NoError(t, err2)
	assert.Equal(t, "abc", gotProduct.id)
	assert.Equal(t, gotProduct, gotProduct2)
}

func TestCatalogService_CacheConcurrent(t *testing.T) {
    catSrv := NewCatalogService(ServiceCache{Enabled: true, Size: 5})

    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    g, ctx := errgroup.WithContext(ctx)

    productID := "abc"

    for i := 2; i <100; i++ {
        g.Go(func() error {
            product, err := catSrv.GetProductByID(productID)
            assert.NoError(t, err)
            assert.Equal(t, productID, product.id)
            return nil
        })
    }

    assert.NoError(t, g.Wait())
}

/********************
Benchmark Tests
********************/
// go test -bench 'BenchmarkCatalogService_GetProductByIDConcurrent' ./cache -benchtime=100x
func BenchmarkCatalogService_GetProductByIDConcurrent(b *testing.B) {
	serviceWithoutCache := NewCatalogService(ServiceCache{Enabled: false})
	serviceWithCache := NewCatalogService(ServiceCache{Enabled: true, Size: 5})

	b.Run("without_cache", func(b *testing.B) {
		benchmarkCatalogServiceGetProductByIDConcurrent(b, serviceWithoutCache)
	})
	b.Run("with_cache", func(b *testing.B) {
		benchmarkCatalogServiceGetProductByIDConcurrent(b, serviceWithCache)
	})
}

func benchmarkCatalogServiceGetProductByIDConcurrent(b *testing.B, s CatalogService) {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    g, ctx := errgroup.WithContext(ctx)

    workers := runtime.NumCPU()
    each := b.N / workers

    b.ReportAllocs()
	b.ResetTimer()
    for i := 0; i < workers; i++ {
        g.Go(func() error {
            for j := 0; j < each; j++ {
                if product, err := s.GetProductByID("abc"); err != nil {
                    return err
                } else {
                    require.Equal(b, "abc", product.id)
                }
            }
            return nil
        })
    }
    require.NoError(b, g.Wait(), "Failed to GetProductByID")
}
