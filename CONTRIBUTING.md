# Contributing to Choochoo

Thank you for your interest in contributing to the Choochoo GitHub webhook server! This document outlines the development practices and guidelines to ensure code quality and maintainability.

## Development Guidelines

### Testing Requirements

**All code changes MUST include comprehensive tests.** This ensures functionality continues to work as expected and prevents regressions.

#### Test Coverage Requirements

- **Unit Tests**: All public functions and methods must have unit tests
- **Integration Tests**: HTTP handlers and routing must have integration tests
- **Edge Cases**: Tests must cover error conditions and edge cases
- **Security Tests**: Authentication and validation logic must be thoroughly tested

#### Test Categories

1. **Function-level Tests**
   - Test all public functions with various inputs
   - Test error conditions and edge cases
   - Test configuration handling

2. **Handler Tests**
   - Test all HTTP endpoints (`/webhook`, `/health`, `/`)
   - Test HTTP method validation
   - Test request/response formats
   - Test error responses

3. **Security Tests**
   - Test webhook signature validation with valid/invalid signatures
   - Test authentication bypass scenarios
   - Test malformed input handling

4. **Integration Tests**
   - Test complete request/response cycles
   - Test routing behavior
   - Test environment variable handling

#### Running Tests

```bash
# Run all tests
go test -v

# Run tests with coverage
go test -v -cover

# Run tests with detailed coverage report
go test -v -coverprofile=coverage.out
go tool cover -html=coverage.out
```

#### Test Standards

- Use descriptive test names that explain what is being tested
- Follow the pattern: `TestFunctionName_Scenario`
- Include both positive and negative test cases
- Use table-driven tests for multiple scenarios when appropriate
- Mock external dependencies appropriately
- Test helper functions should be clearly documented

### Code Quality Standards

#### Code Structure

- Keep functions focused and single-purpose
- Use meaningful variable and function names
- Include appropriate error handling
- Add comments for complex logic

#### Security Requirements

- **Always validate webhook signatures** when secrets are configured
- Use constant-time comparison for signature validation
- Validate all inputs before processing
- Log security-related events appropriately

#### Configuration Management

- Use environment variables for configuration
- Provide sensible defaults
- Document all configuration options
- Validate configuration values

### Development Workflow

1. **Before Making Changes**
   - Run existing tests: `go test -v`
   - Ensure all tests pass
   - Review current code to understand structure

2. **While Developing**
   - Write tests for new functionality first (TDD approach recommended)
   - Run tests frequently: `go test -v`
   - Ensure new code follows existing patterns

3. **Before Submitting**
   - Run full test suite: `go test -v`
   - Check test coverage: `go test -v -cover`
   - Verify code builds: `go build .`
   - Update documentation if needed

### Code Review Checklist

- [ ] All new code has comprehensive tests
- [ ] All tests pass
- [ ] Code follows existing patterns and conventions
- [ ] Security considerations are addressed
- [ ] Documentation is updated if needed
- [ ] Error handling is appropriate
- [ ] Configuration is properly managed

### Example Test Structure

```go
func TestFunctionName_Scenario(t *testing.T) {
    // Arrange - set up test data and conditions
    server := createTestWebhookServer("test-secret")
    payload := []byte(`{"test": "data"}`)
    
    // Act - perform the action being tested
    result := server.validateSignature(payload, "valid-signature")
    
    // Assert - verify the expected outcome
    if !result {
        t.Error("Expected validation to pass with valid signature")
    }
}
```

### Test Helper Functions

When creating test helpers:
- Keep them focused and reusable
- Document their purpose clearly
- Place them at the top of test files
- Use consistent naming patterns

### Continuous Integration

This repository uses automated CI via GitHub Actions to ensure code quality. All changes must pass:

- **Automated Tests**: All existing tests must pass (`make test`)
- **Test Coverage**: Coverage reports are generated (`make coverage`) 
- **Build Verification**: Code must compile successfully (`make build`)
- **Integration Checks**: Full verification suite (`make verify`)

#### CI Process

1. **On Pull Requests**: CI automatically runs on all pull requests to `main`
2. **On Push to Main**: CI runs on direct pushes to the main branch
3. **Branch Protection**: The main branch requires passing CI checks before merging
4. **Status Checks**: Pull requests show CI status and block merging if tests fail

#### Local Development

Before submitting a pull request, ensure your changes pass locally:

```bash
# Run the same checks as CI
make verify

# Or run individual steps
make deps    # Install dependencies
make test    # Run tests
make coverage # Generate coverage report
make build   # Build application
```

This mirrors the CI pipeline and helps catch issues early.

## Questions?

If you have questions about testing requirements or development practices, please open an issue for discussion.

## Remember

**Testing is not optional** - it's a requirement for maintaining a reliable and secure webhook server. Every contribution should include appropriate tests to ensure the functionality continues to work as expected.