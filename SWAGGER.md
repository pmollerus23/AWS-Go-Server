# OpenAPI/Swagger Documentation

Your API now has interactive OpenAPI/Swagger documentation! This provides a visual interface to explore and test your API endpoints.

## Accessing Swagger UI

Once your server is running, visit:

```
http://localhost:8080/swagger/
```

You'll see an interactive documentation page where you can:
- Browse all available endpoints
- See request/response schemas
- Try out endpoints directly from the browser
- View example requests and responses

## Quick Start

### 1. Start the Server

```bash
# Run the server
make run

# Or with hot reload
make dev
```

### 2. Open Swagger UI

Open your browser and navigate to:
```
http://localhost:8080/swagger/
```

### 3. Try an Endpoint

1. Click on an endpoint (e.g., `GET /api/v1/items`)
2. Click **"Try it out"**
3. Click **"Execute"**
4. View the response below

## Available Endpoints

### Health Check
- **GET** `/healthz` - Server health check

### Items (CRUD)
- **GET** `/api/v1/items` - List all items
- **POST** `/api/v1/items` - Create a new item

### AWS Services
- **GET** `/api/v1/aws/s3/buckets` - List S3 buckets
- **GET** `/api/v1/aws/dynamodb/tables` - List DynamoDB tables

## How It Works

### 1. Annotations in Code

Swagger docs are generated from special comments (annotations) in your Go code:

```go
// HandleHealthz returns a simple health check handler.
//
//	@Summary		Health Check
//	@Description	Check if the server is healthy and responding
//	@Tags			health
//	@Produce		plain
//	@Success		200	{string}	string	"OK"
//	@Router			/healthz [get]
func HandleHealthz(logger *slog.Logger) http.HandlerFunc {
    // Implementation
}
```

### 2. Generate Documentation

Annotations are compiled into OpenAPI spec files using the `swag` CLI:

```bash
make swagger
```

This generates:
- `docs/docs.go` - Go code with embedded docs
- `docs/swagger.json` - OpenAPI JSON spec
- `docs/swagger.yaml` - OpenAPI YAML spec

### 3. Serve with Swagger UI

The Swagger UI is served at `/swagger/` endpoint (configured in `internal/server/routes.go`):

```go
mux.Handle("GET /swagger/", http.StripPrefix("/swagger/", httpSwagger.WrapHandler))
```

## Annotation Reference

### General API Info (in main.go)

```go
//	@title						AWS Go Server API
//	@version					1.0
//	@description				A production-grade Go web server
//	@host						localhost:8080
//	@BasePath					/
//	@schemes					http https
```

### Handler Annotations

```go
//	@Summary		Short description
//	@Description	Detailed description
//	@Tags			category
//	@Accept			json
//	@Produce		json
//	@Param			name	path		type	true	"Description"
//	@Param			body	body		Type	true	"Description"
//	@Success		200		{object}	ResponseType
//	@Failure		400		{string}	string	"Error message"
//	@Router			/path [get]
```

### Struct Tags

Add examples and validation info to structs:

```go
type Item struct {
    ID          int64  `json:"id" example:"1"`
    Name        string `json:"name" example:"Sample" minLength:"1" maxLength:"100"`
    Description string `json:"description" example:"Description" maxLength:"500"`
}
```

## Common Annotation Types

### @Summary
Short one-line description of the endpoint.

### @Description
Detailed multi-line description.

### @Tags
Groups endpoints in Swagger UI. Examples: "items", "aws", "health"

### @Accept
Request content type. Common values: `json`, `xml`, `plain`

### @Produce
Response content type. Common values: `json`, `xml`, `plain`

### @Param
Request parameters:
- `@Param name path string true "Description"` - Path parameter
- `@Param name query string false "Description"` - Query parameter
- `@Param body body Type true "Description"` - Request body

### @Success
Successful response:
```go
@Success 200 {object} ResponseType
@Success 200 {array} ItemType
@Success 200 {string} string "OK"
```

### @Failure
Error responses:
```go
@Failure 400 {object} ErrorType "Bad Request"
@Failure 500 {string} string "Internal Server Error"
```

### @Router
Route path and HTTP method:
```go
@Router /api/v1/items [get]
@Router /api/v1/items [post]
@Router /api/v1/items/{id} [delete]
```

### @Security
Authentication requirement (defined in main.go):
```go
@Security BearerAuth
```

## Regenerating Documentation

Whenever you:
- Add new endpoints
- Change existing endpoints
- Modify request/response structures
- Update descriptions

Run:
```bash
make swagger
```

Or manually:
```bash
swag init -g cmd/server/main.go -o docs --parseDependency --parseInternal
```

## Development Workflow

### 1. Write Your Handler

```go
// HandleNewEndpoint handles the new feature
//
//	@Summary		Short description
//	@Description	Detailed description
//	@Tags			category
//	@Produce		json
//	@Success		200	{object}	ResponseType
//	@Router			/api/v1/new [get]
func HandleNewEndpoint(logger *slog.Logger) http.Handler {
    // Implementation
}
```

### 2. Add Route

In `internal/server/routes.go`:
```go
mux.Handle("GET /api/v1/new", handlers.HandleNewEndpoint(s.logger))
```

### 3. Regenerate Docs

```bash
make swagger
```

### 4. Test in Swagger UI

1. Restart server
2. Open `http://localhost:8080/swagger/`
3. Find your new endpoint
4. Try it out!

## OpenAPI Spec Files

Generated files in `docs/` directory:

### docs/swagger.json
Complete OpenAPI 2.0 spec in JSON format. Can be imported into:
- Postman
- Insomnia
- API testing tools
- Code generators

### docs/swagger.yaml
Same spec in YAML format. Human-readable and git-friendly.

### docs/docs.go
Go code with embedded documentation. Imported automatically by the server.

## Advanced Features

### Multiple Response Types

```go
//	@Success	200	{object}	SuccessResponse
//	@Success	201	{object}	CreatedResponse
//	@Failure	400	{object}	ValidationError
//	@Failure	404	{string}	string	"Not found"
//	@Failure	500	{string}	string	"Server error"
```

### Path Parameters

```go
//	@Param	id	path	int	true	"Item ID"
//	@Router	/api/v1/items/{id} [get]
```

### Query Parameters

```go
//	@Param	limit	query	int		false	"Limit results"
//	@Param	offset	query	int		false	"Offset for pagination"
//	@Param	sort	query	string	false	"Sort order"
```

### Request Headers

```go
//	@Param	Authorization	header	string	true	"Bearer token"
//	@Param	X-Request-ID	header	string	false	"Request ID"
```

### Arrays and Objects

```go
// Array response
//	@Success	200	{array}	Item

// Object response
//	@Success	200	{object}	ResponseStruct

// Map response
//	@Success	200	{object}	map[string]interface{}
```

## Customization

### Change Swagger URL

In `internal/server/routes.go`:
```go
// Change from /swagger/ to /docs/
mux.Handle("GET /docs/", http.StripPrefix("/docs/", httpSwagger.WrapHandler))
```

### Change Host/BasePath

In `cmd/server/main.go`:
```go
//	@host		api.example.com
//	@BasePath	/v1
```

### Add Security Schemes

In `cmd/server/main.go`:
```go
//	@securityDefinitions.basic	BasicAuth
//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						X-API-Key
```

## Troubleshooting

### Swagger UI shows 404
- Check that server is running on port 8080
- Try `http://localhost:8080/swagger/index.html`
- Ensure `docs/` directory exists with generated files

### Changes not reflected
- Run `make swagger` to regenerate docs
- Restart the server
- Hard refresh browser (Ctrl+F5)

### "Failed to load API definition"
- Check that `docs` package is imported in `cmd/server/main.go`:
  ```go
  _ "github.com/pmollerus23/go-aws-server/docs"
  ```
- Verify docs were generated: `ls -la docs/`

### Annotations not recognized
- Ensure annotations start with `//` (double slash)
- Check indentation (use tabs, not spaces)
- Run `make swagger` to see error messages

## Best Practices

### 1. Keep Annotations Updated
Update annotations whenever you change endpoints.

### 2. Use Descriptive Names
Be clear about what each endpoint does.

### 3. Document All Responses
Include success and error cases.

### 4. Add Examples
Use `example:` tags in struct fields for realistic data.

### 5. Group Related Endpoints
Use consistent `@Tags` to organize API sections.

### 6. Version Control Docs
Commit generated `docs/` files so team members get latest specs.

### 7. Test in Swagger UI
Use "Try it out" to validate your API works as documented.

## Integration with CI/CD

### Pre-commit Hook

Create `.git/hooks/pre-commit`:
```bash
#!/bin/bash
make swagger
git add docs/
```

### GitHub Actions

```yaml
- name: Generate Swagger docs
  run: |
    go install github.com/swaggo/swag/cmd/swag@latest
    make swagger

- name: Check for changes
  run: git diff --exit-code docs/
```

## Export and Share

### Export OpenAPI Spec

```bash
# Copy spec file
cp docs/swagger.json api-spec.json

# Share with team
# Import into Postman, Insomnia, etc.
```

### Generate Client Libraries

Use OpenAPI Generator to create client SDKs:

```bash
# Install OpenAPI Generator
npm install -g @openapitools/openapi-generator-cli

# Generate Python client
openapi-generator-cli generate -i docs/swagger.json -g python -o clients/python

# Generate JavaScript client
openapi-generator-cli generate -i docs/swagger.json -g javascript -o clients/js
```

## Resources

- [Swag Documentation](https://github.com/swaggo/swag)
- [OpenAPI Specification](https://swagger.io/specification/)
- [Swagger UI](https://swagger.io/tools/swagger-ui/)
- [OpenAPI Generator](https://openapi-generator.tech/)

## Summary

You now have:
- ✅ Interactive API documentation at `/swagger/`
- ✅ Auto-generated from code annotations
- ✅ Try-it-out functionality for testing
- ✅ OpenAPI specs for tools and clients
- ✅ Easy to update with `make swagger`

Happy documenting! 📚
