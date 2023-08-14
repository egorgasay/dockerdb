package dockerdb

import "context"

func With(ctx context.Context, c Config, fn func(c Config, vdb *VDB) error) (err error) {
	vdb, err := New(ctx, c)
	if err != nil {
		if vdb != nil {
			return vdb.Clear(ctx)
		}
		return err
	}
	defer func() { err = vdb.Clear(ctx) }()

	return fn(c, vdb)
}

func WithPostgres(ctx context.Context, name string, fn func(c Config, vdb *VDB) error) (err error) {
	return With(ctx, PostgresConfig(name).Build(), fn)
}

func WithMySQL(ctx context.Context, name string, fn func(c Config, vdb *VDB) error) (err error) {
	return With(ctx, MySQLConfig(name).Build(), fn)
}

func WithRedis(
	ctx context.Context, name string,
	checkWakeUp func(c Config) (stop bool),
	fn func(c Config, vdb *VDB) error,
) (err error) {
	return With(ctx, RedisConfig(name, checkWakeUp).Build(), fn)
}

func WithKeyDB(
	ctx context.Context, name string,
	checkWakeUp func(c Config) (stop bool),
	fn func(c Config, vdb *VDB) error,
) (err error) {
	return With(ctx, KeyDBConfig(name, checkWakeUp).Build(), fn)
}

func WithScyllaDB(
	ctx context.Context, name string,
	checkWakeUp func(c Config) (stop bool),
	fn func(c Config, vdb *VDB) error,
) (err error) {
	return With(ctx, ScyllaDBConfig(name, checkWakeUp).Build(), fn)
}
