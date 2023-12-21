package hundApiV1

import (
	"context"
	"net/http"
)

func Expand(expansions ...string) RequestEditorFn {
	return expansionOption("expand", expansions...)
}

func Unexpand(unexpansions ...string) RequestEditorFn {
	return expansionOption("unexpand", unexpansions...)
}

func expansionOption(optionName string, expansions ...string) RequestEditorFn {
	return func(ctx context.Context, req *http.Request) error {
		query := req.URL.Query()

		for _, v := range expansions {
			query.Add(optionName+"[]", v)
		}

		req.URL.RawQuery = query.Encode()

		return nil
	}
}
