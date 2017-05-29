# andra

andra is a plugin extending goa to generate boilerplate code required to access nosql databases, primarily cassandra, 
aiding developers to focus on business logic from very early stages of development.

Install:
`go get gihub.com/ztx/andra`

run:
`goagen gen --pkg-path github.com/ztx/andra -d path-to/design-pkg`

eg:
`goagen gen --pkg-path github.com/ztx/andra -d github.com/ztx/andra/example/design`
