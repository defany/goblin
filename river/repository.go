package river

import (
	"cmp"
	"context"
	"fmt"

	"github.com/defany/goblin/errfmt"
	"github.com/defany/goblin/pg"
	"github.com/gookit/goutil/arrutil"
	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/rivertype"
)

var (
	ErrJobIDsNotProvided = fmt.Errorf("river: job ids not provided")
	ErrJobNotFound       = fmt.Errorf("river: job not found")
)

type Repository[T river.JobArgs] struct {
	client *river.Client[pgx.Tx]
}

func New[T river.JobArgs](client *river.Client[pgx.Tx]) *Repository[T] {
	return &Repository[T]{client: client}
}

func (r *Repository[T]) Insert(ctx context.Context, args T, opts ...*river.InsertOpts) (*rivertype.JobInsertResult, error) {
	opt := cmp.Or(opts...)

	if tx := pg.ExtractTx(ctx); tx != nil {
		res, err := r.client.InsertTx(ctx, tx, args, opt)
		return res, errfmt.WithSource(err)
	}

	res, err := r.client.Insert(ctx, args, opt)
	return res, errfmt.WithSource(err)
}

func (r *Repository[T]) InsertMany(ctx context.Context, args []T, opts ...*river.InsertOpts) ([]*rivertype.JobInsertResult, error) {
	opt := cmp.Or(opts...)

	params := arrutil.Map(args, func(arg T) (river.InsertManyParams, bool) {
		return river.InsertManyParams{
			Args:       arg,
			InsertOpts: opt,
		}, true
	})

	if tx := pg.ExtractTx(ctx); tx != nil {
		res, err := r.client.InsertManyTx(ctx, tx, params)
		return res, errfmt.WithSource(err)
	}

	res, err := r.client.InsertMany(ctx, params)
	return res, errfmt.WithSource(err)
}

func (r *Repository[T]) FetchJob(ctx context.Context, id int64) (*rivertype.JobRow, error) {
	jobs, err := r.FetchJobs(ctx, id)
	if err != nil {
		return nil, errfmt.WithSource(err)
	}

	if len(jobs) == 0 {
		return nil, ErrJobNotFound
	}

	return jobs[0], nil
}

func (r *Repository[T]) FetchJobs(ctx context.Context, ids ...int64) ([]*rivertype.JobRow, error) {
	if len(ids) == 0 {
		return nil, ErrJobIDsNotProvided
	}

	result, err := r.client.JobList(ctx, river.NewJobListParams().IDs(ids...))
	if err != nil {
		return nil, errfmt.WithSource(err)
	}

	return result.Jobs, nil
}

func (r *Repository[T]) CancelJobs(ctx context.Context, ids ...int64) error {
	if len(ids) == 0 {
		return ErrJobIDsNotProvided
	}

	tx := pg.ExtractTx(ctx)
	for _, id := range ids {
		var err error
		if tx != nil {
			_, err = r.client.JobCancelTx(ctx, tx, id)
		} else {
			_, err = r.client.JobCancel(ctx, id)
		}

		if err != nil {
			return errfmt.WithSource(err)
		}
	}

	return nil
}

func (r *Repository[T]) DeleteJobs(ctx context.Context, ids ...int64) error {
	if len(ids) == 0 {
		return ErrJobIDsNotProvided
	}

	const query = "DELETE FROM river_job WHERE id = ANY($1)"

	if tx := pg.ExtractTx(ctx); tx != nil {
		_, err := tx.Exec(ctx, query, ids)
		return errfmt.WithSource(err)
	}

	return errfmt.WithSource(r.client.Driver().GetExecutor().Exec(ctx, query, ids))
}
