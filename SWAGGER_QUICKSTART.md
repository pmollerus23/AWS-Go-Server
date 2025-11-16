# Swagger/OpenAPI Quick Start

## View Your API Documentation

### 1. Start the Server
```bash
make run
```

### 2. Open Swagger UI
Visit in your browser:
```
http://localhost:8080/swagger/
```

### 3. Try It Out!
1. Click on any endpoint (e.g., `GET /api/v1/items`)
2. Click **"Try it out"**
3. Click **"Execute"**
4. See the response!

## After Making Changes

### When you modify endpoints or add new ones:

```bash
# Regenerate Swagger docs
make swagger

# Restart the server
make run
```

## Available Documentation

- **Swagger UI**: `http://localhost:8080/swagger/` - Interactive web interface
- **JSON Spec**: `docs/swagger.json` - Import into Postman, Insomnia, etc.
- **YAML Spec**: `docs/swagger.yaml` - Human-readable spec
- **Full Guide**: [SWAGGER.md](./SWAGGER.md) - Complete documentation

## Quick Commands

```bash
make swagger        # Generate/update Swagger docs
make run           # Run server and view at /swagger/
make build         # Build with Swagger included
make install-tools # Install swag CLI tool
```

## What You Can Do

✅ **Browse API** - See all endpoints organized by tags
✅ **Test Endpoints** - Try requests directly from browser
✅ **View Schemas** - See request/response data structures
✅ **Export Spec** - Use `docs/swagger.json` in other tools
✅ **Generate Clients** - Create SDKs from OpenAPI spec

## Example: Testing an Endpoint

### Create an Item
1. Go to `http://localhost:8080/swagger/`
2. Find `POST /api/v1/items` under "items" section
3. Click **"Try it out"**
4. Modify the request body:
   ```json
   {
     "name": "My Test Item",
     "description": "Testing from Swagger"
   }
   ```
5. Click **"Execute"**
6. See the response with the created item ID!

### List All Items
1. Find `GET /api/v1/items`
2. Click **"Try it out"**
3. Click **"Execute"**
4. See your created items!

That's it! You now have full OpenAPI/Swagger documentation for your API. 🎉
