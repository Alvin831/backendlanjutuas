# Simple Unit Tests

This directory contains simple function-level unit tests for the Achievement Management System.

## Testing Approach

We use a **simple testing strategy** that focuses on:
- Testing individual functions in isolation
- No complex mocking or external dependencies
- Direct function calls with various inputs
- Assertion-based validation

## Test Files

### `password_test.go`
Tests password utility functions:
- `HashPassword()` - Password hashing
- `CheckPasswordHash()` - Password verification

### `jwt_test.go`
Tests JWT utility functions:
- `GenerateToken()` - JWT token generation
- `ParseToken()` - JWT token parsing and validation

### `response_test.go`
Tests API response utility functions:
- `SuccessResponse()` - Success response formatting
- `ErrorResponse()` - Error response formatting

### `validation_test.go`
Tests validation utility functions:
- `ValidateEmail()` - Email format validation
- `ValidateUsername()` - Username validation
- `ValidatePassword()` - Password strength validation
- `IsEmptyString()` - Empty string checking

### `math_test.go`
Tests mathematical utility functions:
- `CalculatePoints()` - Achievement points calculation
- `CalculateAverage()` - Average calculation
- `GetMaxValue()` - Maximum value finder
- `GetMinValue()` - Minimum value finder

### `service_helpers_test.go`
Tests service helper functions:
- `validateUserInput()` - User input validation
- `calculatePagination()` - Pagination calculation
- `calculateTotalPages()` - Total pages calculation
- `formatUserResponse()` - User response formatting

### `achievement_logic_test.go`
Tests achievement business logic:
- `validateAchievementData()` - Achievement data validation
- `calculateAchievementPointsByLevel()` - Points calculation by level
- `canUserAccessAchievement()` - Role-based access control
- `validateAchievementStatus()` - Status validation
- `canChangeAchievementStatus()` - Status change validation

### `student_logic_test.go`
Tests student business logic:
- `validateStudentData()` - Student data validation
- `validateNIM()` - NIM format validation
- `formatStudentName()` - Name formatting
- `calculateStudentGPA()` - GPA calculation
- `validateStudentProgram()` - Program validation
- `calculateStudentSemester()` - Semester calculation

### `auth_logic_test.go`
Tests authentication business logic:
- `validateLoginCredentials()` - Login validation
- `validateRegistrationData()` - Registration validation
- `checkUserPermissions()` - Permission checking
- `validateRoleID()` - Role ID validation
- `isTokenExpired()` - Token expiration check
- `generateRefreshToken()` - Refresh token generation
- `validateUserRole()` - Role validation

## Running Tests

```bash
# Run all simple tests
go test -v ./tests/simple/...

# Run with coverage
go test -v -cover ./tests/simple/...

# Run using the test runner
go run run_tests.go
```

## Test Coverage

The simple tests cover:
- ✅ Core utility functions
- ✅ Input validation
- ✅ Edge cases (empty inputs, invalid data)
- ✅ Expected behavior verification
- ✅ Error handling

## Benefits of Simple Testing

1. **Easy to understand** - No complex setup or mocking
2. **Fast execution** - Direct function calls
3. **Reliable** - Tests actual function behavior
4. **Maintainable** - Simple test structure
5. **Focused** - One function per test case