# User-Role M2M Refactoring with Custom Join Table

## Overview

Successfully refactored the User-Role relationship to use GORM's M2M associations while maintaining a **custom join table** with CASCADE delete constraints.

## Key Features

### Custom Join Table with CASCADE Delete

Your custom `UserRole` model in `/internal/domain/user_role/user_role_model.go`:

```go
package userrole

import "time"

type UserRole struct {
    UserID    string    `gorm:"type:uuid;not null;index;constraint:OnDelete:CASCADE"`
    RoleID    string    `gorm:"type:uuid;not null;index;constraint:OnDelete:CASCADE"`
    CreatedAt time.Time `gorm:"autoCreateTime"`
}

func (UserRole) TableName() string {
    return "user_roles"
}
```

**Benefits of Custom Join Table:**
- ✅ **CASCADE Delete**: When a user or role is deleted, related entries are automatically removed
- ✅ **CreatedAt Timestamp**: Track when role was assigned
- ✅ **Custom Constraints**: Full control over database constraints
- ✅ **GORM Compatibility**: Works seamlessly with GORM M2M associations

## Implementation

### User Model

```go
type User struct {
    // ... other fields ...
    Roles []role.Role `gorm:"many2many:user_roles;joinForeignKey:UserID;joinReferences:RoleID" json:"-"`
}
```

**Important Tags:**
- `many2many:user_roles` - Specifies the join table name
- `joinForeignKey:UserID` - Foreign key from User to join table
- `joinReferences:RoleID` - Foreign key from Role to join table
- `json:"-"` - Don't expose roles in JSON (fetch separately if needed)

### Database Migration

```go
err := db.AutoMigrate(
    &user.User{},
    &role.Role{},
    &userrole.UserRole{}, // Custom join table with CASCADE
    // ... other models
)
```

**Note:** The custom `UserRole` model must be migrated explicitly to ensure CASCADE constraints are created.

## Architecture: Separation of Concerns

### User Repository (Owns the Relationship)

```go
type UserRepository interface {
    // ... CRUD operations ...
    
    // Role association operations
    AssignRole(ctx context.Context, userID, roleID string) error
    RemoveRole(ctx context.Context, userID, roleID string) error
    GetRoles(ctx context.Context, userID string) ([]role.Role, error)
}
```

**Implementation:**
```go
// AssignRole uses GORM Association API
func (r *userRepository) AssignRole(ctx context.Context, userID, roleID string) error {
    var u user.User
    if err := r.db.WithContext(ctx).First(&u, "id = ?", userID).Error; err != nil {
        return fmt.Errorf("user not found: %w", err)
    }

    var ro role.Role
    if err := r.db.WithContext(ctx).First(&ro, "id = ?", roleID).Error; err != nil {
        return fmt.Errorf("role not found: %w", err)
    }

    // GORM automatically handles the custom join table
    return r.db.WithContext(ctx).Model(&u).Association("Roles").Append(&ro)
}
```

### Role Repository (No User Dependencies)

```go
type RoleRepository interface {
    Create(ctx context.Context, role *role.Role) error
    GetByID(ctx context.Context, id string) (*role.Role, error)
    GetByName(ctx context.Context, name string) (*role.Role, error)
    GetAll(ctx context.Context) ([]role.Role, error)
    
    // Read-only: Get users with this role
    GetRoleUsers(ctx context.Context, roleID string) ([]string, error)
}
```

## CASCADE Delete Behavior

### When User is Deleted
```sql
DELETE FROM users WHERE id = 'user-uuid';
-- Automatically triggers:
DELETE FROM user_roles WHERE user_id = 'user-uuid';
```

### When Role is Deleted
```sql
DELETE FROM roles WHERE id = 'role-uuid';
-- Automatically triggers:
DELETE FROM user_roles WHERE role_id = 'role-uuid';
```

**This prevents orphaned records in the join table!**

## Files Modified

1. ✅ `internal/domain/user/user_model.go` - Added Roles field with custom join table tags
2. ✅ `pkg/database/database.go` - Added `userrole.UserRole` to AutoMigrate
3. ✅ `internal/repository/user_repository.go` - Added role operations
4. ✅ `internal/repository/gormrepo/user_repository.go` - Implemented role operations
5. ✅ `internal/repository/role_repository.go` - Removed user-role operations
6. ✅ `internal/repository/gormrepo/role_repository.go` - Removed user-role methods
7. ✅ `internal/service/auth_service.go` - Use UserRepository.AssignRole
8. ✅ `internal/service/user_service.go` - Use UserRepository.GetRoles
9. ✅ `seeds/seed_users.go` - Use Association API

## Comparison: Custom vs Default Join Table

### Default GORM Join Table
```go
// User model
Roles []role.Role `gorm:"many2many:user_roles"`

// GORM creates:
// - user_id (no CASCADE)
// - role_id (no CASCADE)
// - No timestamps
```

### Your Custom Join Table ✨
```go
// User model with custom table
Roles []role.Role `gorm:"many2many:user_roles;joinForeignKey:UserID;joinReferences:RoleID"`

// Custom UserRole model with:
// - UserID with CASCADE delete
// - RoleID with CASCADE delete  
// - CreatedAt timestamp
// - Custom constraints
```

## Usage Examples

### Assign Role to User
```go
err := userRepo.AssignRole(ctx, userID, roleID)
// GORM automatically uses the custom join table
```

### Remove Role from User
```go
err := userRepo.RemoveRole(ctx, userID, roleID)
// CASCADE ensures referential integrity
```

### Get User's Roles
```go
roles, err := userRepo.GetRoles(ctx, userID)
// Uses Preload with custom join table
```

### Seeding with Association API
```go
// Create user
db.Create(&user)

// Assign role using Association
var role role.Role
db.Where("name = ?", "admin").First(&role)
db.Model(&user).Association("Roles").Append(&role)
// Automatically populates custom UserRole table with CASCADE constraints
```

## Benefits Summary

1. ✅ **CASCADE Delete Protection**: No orphaned join table records
2. ✅ **Audit Trail**: CreatedAt shows when role was assigned
3. ✅ **Clean Architecture**: User owns relationship, no circular dependencies
4. ✅ **GORM Compatible**: Full support for Association API
5. ✅ **Type Safety**: Strongly typed operations
6. ✅ **Maintainable**: Clear separation of concerns
7. ✅ **Database Integrity**: Custom constraints at DB level

## Testing

Verify CASCADE behavior:
```sql
-- Create test user with role
INSERT INTO users (id, email) VALUES ('test-id', 'test@test.com');
INSERT INTO user_roles (user_id, role_id) VALUES ('test-id', 'role-id');

-- Delete user - should cascade to user_roles
DELETE FROM users WHERE id = 'test-id';

-- Verify join table is cleaned up
SELECT * FROM user_roles WHERE user_id = 'test-id';
-- Should return 0 rows
```

## Conclusion

Your custom join table approach provides the best of both worlds:
- GORM's convenient M2M association API
- Database-level CASCADE delete protection
- Audit trail with timestamps
- Full control over constraints

This is a production-ready implementation that maintains data integrity while keeping code clean and maintainable! 🎉
