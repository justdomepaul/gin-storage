package cloud

import (
	"cloud.google.com/go/storage"
	"context"
	"github.com/justdomepaul/gin-storage/pkg/config"
	"github.com/justdomepaul/gin-storage/pkg/panicerrorhandler"
	"go.uber.org/zap"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
)

var (
	newClient = storage.NewClient
)

// ISession interface
type ISession interface {
	Bucket(name string) *storage.BucketHandle
	Buckets(ctx context.Context, projectID string) *storage.BucketIterator
	HMACKeyHandle(projectID, accessID string) *storage.HMACKeyHandle
	CreateHMACKey(ctx context.Context, projectID, serviceAccountEmail string, opts ...storage.HMACKeyOption) (*storage.HMACKey, error)
	ListHMACKeys(ctx context.Context, projectID string, opts ...storage.HMACKeyOption) *storage.HMACKeysIterator
	Close() error
	ServiceAccount(ctx context.Context, projectID string) (string, error)
}

// NewSession method
var NewSession = func(ctx context.Context, opt config.Cloud) (*storage.Client, error) {
	options := make([]option.ClientOption, 0)
	if opt.EndPoint != "" {
		options = append(options, option.WithEndpoint(opt.EndPoint))
	}
	if opt.WithoutAuthentication {
		options = append(options, option.WithoutAuthentication())
	}
	if opt.GRPCInsecure {
		options = append(options, option.WithGRPCDialOption(grpc.WithInsecure()))
	}
	return newClient(ctx, options...)
}

func NewExtendStorageDatabase(logger *zap.Logger, opt config.Cloud) (*storage.Client, func(), error) {
	session, err := NewSession(context.Background(), opt)
	if err != nil {
		return nil, nil, err
	}
	logger.Info("Storage init complete", zap.String("system", "Database"))

	return session,
		func() {
			defer panicerrorhandler.PanicErrorHandler("Core", "")
			if err := session.Close(); err != nil {
				panic(err)
			}
		}, nil
}
