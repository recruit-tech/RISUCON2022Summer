package client

type ClientType = string

var (
	InitializerType          ClientType = "initializer"
	CompatibilityCheckerType ClientType = "compatibility checker"
	LoaderType               ClientType = "loader"
)
