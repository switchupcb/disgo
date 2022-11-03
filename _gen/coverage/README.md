# Coverage Integration Test Generator

Disgo contains an integration test that covers nearly 100% of the Discord API to ensure feature-complete functionality. This test entails calling 160+ endpoints (requests) and dealing with respective events as necessary. The test is run in the CI/CD pipeline used to approve a build for production-usage. It can also be used by developers to debug issues. As a result, it's important to optimize this task in order to _minimize the amount of time spent running tests_ in any given workflow.

## Guide

The test itself represents a story between an ADMIN and USER from a "new account". As a result, the ADMIN must setup its Guild, Channels, Roles, and other features programmatically. The USER can be used to verify that these actions occur, or to trigger events that may not be possible to trigger with one account.

### Optimization

Running the entire story can be optimized to approximately 2 - 4 seconds. Discord's Global Rate Limit allows 50 requests per second. In an ideal run, both bots can send 50 - 100 requests at any given time in batches. However, **certain requests are dependent on other requests**. This means that we must optimize the order of the request calls in order to optimize the performance of the test.

## Generation

A general order of request is generated using the `main.go` file in this directory. It consists of every endpoint placed in a map (of dependent endpoints _to_ dependencies), along with a topological sort to output a **valid order to call requests**. Requests that are **NOT** included in the test should be marked in the `unused` map.

Copygen is used to speed up the test file's creation. First, requests are placed in a [`setup.go`](/wrapper/copygen/integration/setup.go) file _(using the same definitions as the [`send.go` setup file](/wrapper/copygen/requests/setup.go)_) in the specified order _(which can be adjusted manually)_. Then, Copygen is run using the customized template to output a test that calls functions in an `errgroup`.

Once the general test is output, it can be further modified where necessary. In the case of dependent requests, it may be more efficient to move them to the same goroutine as their dependency, such that the dependent is called synchronously AFTER the dependency. In other cases, a dependency that maintains many dependents may need to be called on its own, with some capacity to ensure it was called BEFORE the dependent is called _(i.e `WaitGroups`)_.