# Casbin Authorization Refactoring

## Overview

The Casbin authorization system has been refactored from a **route-pattern based approach** to a **resource-action based approach**. This change makes the authorization system more reliable, maintainable, and flexible.

## The Problem with the Old Approach

The previous implementation had a critical flaw:

```go
// OLD APPROACH - Route Pattern Matching
{"admin", "/api/v1/users/*", "GET"}
{"admin", "/api/v1/novels/:id", "PUT"}
```

**Issues:**
1. **Brittle Pattern Matching**: Required exact URL patterns to match
2. **Hard to Maintain**: Any route change required policy updates
3. **Pattern Ambiguity**: Complex paths like `/novels/123/chapters/456` were hard to normalize
4. **Tight Coupling**: Authorization logic was coupled to URL structure
5. **Error-Prone**: Easy to forget patterns or create mismatches

## The New Approach - Resource-Action Model

The new implementation uses semantic resource names and actions:

```go
// NEW APPROACH - Resource-Action Model
{"admin", "user", "read"}
{"admin", "novel", "update"}
{"translator", "novel_translation", "create"}
```

### Benefits

1. ✅ **Decoupled from URLs**: Authorization logic is independent of route structure
2. ✅ **Semantic and Clear**: Uses business domain terms (novel, chapter, user)
3. ✅ **Flexible**: Easy to add new resources and actions
4. ✅ **Maintainable**: Policy changes don't require code changes
5. ✅ **Testable**: Much easier to test permissions logic
6. ✅ **Self-Documenting**: Policies are easy to understand

## Architecture

### Resources

Resources represent the entities in your domain:

- `user` - User accounts
- `novel` - Novel entries
- `novel_translation` - Novel translations
- `chapter` - Novel chapters
- `chapter_translation` - Chapter translations

### Actions

Actions represent operations on resources:

- `list` - Retrieve multiple resources
- `read` - Retrieve a single resource
- `create` - Create a new resource
- `update` - Modify an existing resource
- `delete` - Remove a resource

### Roles

Predefined roles with specific permissions:

- **admin**: Full access to all resources
- **author**: Can manage novels and chapters
- **translator**: Can manage translations
- **user**: Basic read access

## Implementation Details

### 1. Middleware Enhancement

The middleware now uses context-based resource/action metadata:

```go
// Set resource and action in context
middleware.SetReSource404NovelDownloaderction(c, "novel", "create")

// Or use the convenience function
middleware.RequirePermission("novel", "create", enforcer, roleGetter)
```

### 2. Route Configuration

Routes explicitly declare their resource and action:

```go
novels.POST("", 
    middleware.RequirePermission("novel", "create", cfg.Enforcer, roleGetter), 
    cfg.NovelHandler.Create)

novels.PUT("/:id", 
    middleware.RequirePermission("novel", "update", cfg.Enforcer, roleGetter), 
    cfg.NovelHandler.Update)
```

### 3. Policy Definition

Policies are defined using resource-action pairs:

```go
// Admin can do everything with users
{"admin", "user", "list"},
{"admin", "user", "read"},
{"admin", "user", "update"},
{"admin", "user", "delete"},

// Translator can manage translations
{"translator", "novel_translation", "create"},
{"translator", "novel_translation", "update"},
```

## Usage Examples

### Adding a New Resource

1. Define the resource name (e.g., `comment`)
2. Add policies in `pkg/casbin/casbin.go`:

```go
{"user", "comment", "create"},
{"admin", "comment", "delete"},
```

3. Use in routes:

```go
comments.POST("", 
    middleware.RequirePermission("comment", "create", cfg.Enforcer, roleGetter),
    cfg.CommentHandler.Create)
```

### Adding a New Role

1. Define policies for the role:

```go
{"moderator", "comment", "read"},
{"moderator", "comment", "update"},
{"moderator", "comment", "delete"},
```

2. Assign role to users in the database

### Custom Authorization Logic

For complex authorization (e.g., "only owner can edit"), use handler-level checks:

```go
func (h *NovelHandler) Update(c *gin.Context) {
    userID, _ := middleware.GetUserID(c)
    novelID := c.Param("id")
    
    // Casbin already checked that user has "novel:update" permission
    // Now check if user owns the novel
    novel, _ := h.novelService.GetByID(c.Request.Context(), novelID)
    
    if novel.CreatedBy != userID {
        // Check if user is admin
        roles, _ := h.userService.GetUserRoles(c.Request.Context(), userID)
        isAdmin := contains(roles, "admin")
        
        if !isAdmin {
            response.Error(c, http.StatusForbidden, "Can only edit your own novels", nil)
            return
        }
    }
    
    // Proceed with update...
}
```

## Migration Guide

### For Existing Policies

Old policies are automatically replaced when you restart the application. The `InitializeDefaultPolicies` function will:

1. Clear old route-based policies
2. Add new resource-action policies
3. Save to database

### Testing After Migration

1. Test each endpoint with different roles:

```bash
# Test as admin
curl -H "Authorization: Bearer $ADMIN_TOKEN" \
     -X POST http://localhost:8080/api/v1/novels \
     -d '{"novel": {...}, "translation": {...}}'

# Test as translator
curl -H "Authorization: Bearer $TRANSLATOR_TOKEN" \
     -X POST http://localhost:8080/api/v1/novels/translations \
     -d '{"novel_id": "123", "language": "es", ...}'

# Test as regular user (should fail)
curl -H "Authorization: Bearer $USER_TOKEN" \
     -X DELETE http://localhost:8080/api/v1/novels/123
```

2. Check Casbin logs for authorization decisions

## Policy Management

### View Current Policies

```go
policies := enforcer.GetPolicy()
for _, policy := range policies {
    fmt.Printf("Role: %s, Resource: %s, Action: %s\n", 
               policy[0], policy[1], policy[2])
}
```

### Add Policy at Runtime

```go
enforcer.AddPolicy("moderator", "comment", "delete")
enforcer.SavePolicy()
```

### Remove Policy

```go
enforcer.RemovePolicy("user", "novel", "delete")
enforcer.SavePolicy()
```

## Best Practices

1. **Keep Resources Granular**: Separate concerns (e.g., `novel` vs `novel_translation`)
2. **Use Consistent Naming**: Follow a naming convention (lowercase, underscore-separated)
3. **Document Policies**: Keep this documentation updated with new resources
4. **Test Thoroughly**: Test all role-permission combinations
5. **Audit Regularly**: Review policies for security implications
6. **Use Handler-Level Checks**: For ownership and business logic validation

## Troubleshooting

### "Authorization configuration error"

**Cause**: Resource and action not set in context

**Solution**: Ensure you're using `RequirePermission` middleware or manually setting resource/action

### "Access denied" when it should work

**Cause**: Missing policy or incorrect role assignment

**Solution**: 
1. Check user roles: `SELECT * FROM user_roles WHERE user_id = ?`
2. Check policies: `SELECT * FROM casbin_rule`
3. Verify policy matches resource and action exactly

### Policies not persisting

**Cause**: Not calling `SavePolicy()` after changes

**Solution**: Always call `enforcer.SavePolicy()` after policy modifications

## Performance Considerations

- Casbin caches policies in memory for fast lookups
- Policy checks are O(1) for most cases
- Database queries only happen on policy modifications or server restart
- No performance impact from complex URL parsing anymore

## Future Enhancements

Potential improvements to consider:

1. **Policy API**: REST endpoints to manage policies dynamically
2. **Policy UI**: Admin interface for policy management
3. **Policy Versioning**: Track policy changes over time
4. **Attribute-Based Access Control (ABAC)**: Add support for attributes (time, location, etc.)
5. **Policy Testing Framework**: Automated tests for policy configurations

## References

- [Casbin Documentation](https://casbin.org/docs/overview)
- [RBAC Model](https://casbin.org/docs/rbac)
- [Best Practices](https://casbin.org/docs/best-practice)
