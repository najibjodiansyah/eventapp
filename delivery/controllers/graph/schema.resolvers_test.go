package graph

import (
	"eventapp/utils/graph/generated"

	"testing"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
)

func newResolver() generated.Config {
	r := Resolver{}

	return generated.Config{
		Resolvers: &r,
	}
}

func Test(t *testing.T) {
	c := client.New(handler.NewDefaultServer(generated.NewExecutableSchema(newResolver())))

	var resp struct {
		users []struct{ name string }
	}

	t.Run("users", func(t *testing.T) {
		c.MustPost(`{ users { name } }`, &resp)
	})
}
