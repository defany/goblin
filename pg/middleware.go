package pg

import "context"

type Request struct {
	Query string
	Args  []any
}

type Middleware func(ctx context.Context, req Request) (context.Context, Request, error)

func applyMiddlewares(ctx context.Context, mws []Middleware, query string, args []any) (context.Context, string, []any, error) {
	if len(mws) == 0 {
		return ctx, query, args, nil
	}

	req := Request{Query: query, Args: args}

	var err error
	for _, mw := range mws {
		ctx, req, err = mw(ctx, req)
		if err != nil {
			return ctx, "", nil, err
		}
	}

	return ctx, req.Query, req.Args, nil
}
