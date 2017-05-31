package store

import (
	"encoding/json"
	"fmt"
	"graphql"
	"jdy/pkg/api"
	"jdy/pkg/api/reactlu"
	"jdy/pkg/backend/gql/handler"
	"jdy/pkg/lu"
	"jdy/pkg/site/store/gql"
	"log"
	"net/http"
	"strconv"
)

var pageTemplate = `<!doctype html>
<html>
	<head>
		<meta charset="utf-8">
		<meta content="width=device-width, user-scalable=no, initial-scale=1.0" name="viewport">
		<title>JDY</title>
		<link rel="shortcut icon" href="{{.AssetURL}}/favicon.ico">
		<link rel="stylesheet" href="//dn-applysquare-lib.qbox.me/bootstrap/3.3.6/css/bootstrap.min.css" />
    <link rel="stylesheet" href="//dn-applysquare-lib.qbox.me/font-awesome/4.6.3/css/font-awesome.min.css" />
	</head>
	<body>
		<div id="app" style="position:absolute;left:0;right:0;top:0;bottom:0;display:none"></div>
		<form action="/ajax/graphql" method="post">
			<p>First name: <input type="text" name="Query" /></p>
			<p>Last name: <input type="text" name="lname" /></p>
			<input type="submit" value="Submit" />
		</form>
		<script src="{{.AssetURL}}/vendor.js"></script>
		<script src="{{.AssetURL}}/app.js"></script>
	</body>
</html>`

// var UserType = graphql.NewObject(
// 	graphql.ObjectConfig{
// 		Name: "User",
// 		Fields: graphql.Fields{
// 			"id": &graphql.Field{
// 				Type: graphql.Int,
// 			},
// 			"email": &graphql.Field{
// 				Type: graphql.String,
// 			},
// 		},
// 	},
// )

var QueryType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Query",
	Fields: graphql.Fields{
		"user": &graphql.Field{
			Type: gql.UserType,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Description: "User ID",
					Type:        graphql.NewNonNull(graphql.ID),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				i := p.Args["id"].(string)
				id, err := strconv.Atoi(i)
				if err != nil {
					return nil, err
				}
				return gql.GetUserByID(id)
			},
		},
	},
})

// Routes returns routes for store.
func Routes() []api.Route {

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: QueryType,
		// Query:    gql.QueryType,
		// Mutation: gql.MutationType,
	})
	if err != nil {
		log.Fatal(err)
	}

	routes := []api.Route{
		{
			Method:      "POST",
			Pattern:     "/ajax/graphql",
			HandlerFunc: handleGraphQL(schema),
			// HandlerFunc: lu.PermissionRequired(inj, nil, gql.MakeGraphQLHandler("graphql.ccss")),
		},
		{
			Method:  "GET",
			Pattern: "/admin/graphiql",
			// HandlerFunc: lu.PermissionRequired(inj, nil, gql.MakeGraphIQLHandler("/ajax/graphql")),
		},
	}

	routes = append(routes, reactlu.NewRoutes(&reactlu.RouteParams{
		AssetName:          "jdy",
		RootPath:           "",
		PageTemplate:       pageTemplate,
		PermissionRequired: nil,
	})...)

	return routes
}

func handleGraphQL(schema graphql.Schema) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		opts, err := handler.NewRequestOptions(r)
		if err != nil {
			lu.Error(err)(w)
			return
		}

		fmt.Println("opts: ", opts)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		result := QueryFull(schema, opts.Query, opts.Variables, opts.OperationName)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(result)
	}
}

// QueryFull accepts more arguments for query.
func QueryFull(schema graphql.Schema, query string, vars map[string]interface{}, operationName string) *graphql.Result {
	params := graphql.Params{
		Schema:         schema,
		RequestString:  query,
		VariableValues: vars,
		OperationName:  operationName,
	}

	return graphql.Do(params)
}
