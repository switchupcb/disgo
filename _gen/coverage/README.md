# Coverage Integration Test Generator

Disgo contains an integration test that covers a majority of the Discord API to ensure feature-complete functionality. This test entails calling 100+ endpoints (requests) and dealing with respective events as necessary. The test is run in the CI/CD pipeline used to approve a build for production-usage. It can also be used by developers to debug issues. As a result, it's important to optimize this task in order to _minimize the amount of time spent running tests_ in any given workflow.

## Guide

The test is run within the context of a Guild.

### Optimization

Running the entire test is optimized to approximately 2 - 4 seconds. Discord's Global Rate Limit allows 50 requests per second. However, **certain requests are dependent on other requests**. This means that the order of the request calls must be optimized in order to increase the performance of the test.

#### Generation

A general order of requests is generated using the `main.go` file in this directory. It consists of every endpoint placed in a map (of dependent endpoints _to_ dependencies), along with a topological sort to output a **valid order to call requests**. Requests that are **NOT** included in the test are marked in the `unused` map. 

Once the request order is output, the test can be created and/or modified. Dependent requests may be moved to same goroutine as their dependencies, such that the dependent request is called synchronously AFTER the dependency. In other cases, a dependency that maintains many dependents may need to be called on its own, with some capacity to ensure it was called BEFORE the dependent is called _(i.e `WaitGroups`)_.