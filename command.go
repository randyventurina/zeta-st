package main

type Alias = string

type list struct {
	Add Alias
	Get Alias
}

// Enum for public use
var Command = &list{
	Add: "add",
	Get: "get",
}
