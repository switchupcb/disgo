# Coverage Integration Test Generator

Disgo uses an integration test that covers a majority of the Discord API to prove feature-complete functionality.
- The test calls 70+ endpoints (requests) and handles respective events.
- The test is run in the CI/CD pipeline used to approve a build for production-usage. 
- The test is also used by developers to debug issues.

So, this test must be optimized to _minimize the amount of time spent running tests_ in any given workflow.

## Guide

The test is run within the context of a Guild.

### Optimization

The test is optimized to run in less than 10 seconds. 

Discord's Global Rate Limit allows 50 requests per second. However, **certain requests are dependent on other requests**. So, the "order of the request" calls must be optimized to increase the speed of the test.

#### Generation

Use `go build -o coverage` from the [`./_gen/coverage`](/_gen/coverage/) directory to build the executable file for the "order of requests" generator. You may be required to set the `GOWORK` environment variable to `off`.

Calling the executable file generates a linear "order of requests".

##### How does it work?

1. Every endpoint is placed in a map (of dependent endpoints _to_ dependencies).
2. Requests that should **NOT** be included in the test are marked in the `unused` map.
3. A topological sort outputs a **valid order to call requests**.
4. The integration test is checked for requests which are not present in the current coverage test, but should be.

The coverage test can be modified manually at this point.
- Dependent requests may be moved to same goroutine as their dependencies, such that the dependent request is called synchronously AFTER the dependency.
- Otherwise, a dependency that maintains many dependents may need to be called on its own, with some capacity to ensure it was called BEFORE the dependent is called _(i.e `WaitGroups`)_.