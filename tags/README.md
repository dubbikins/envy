# Tags

The directories under `/tags` are each their own go module. This is done to keep the core envy library slim while allow others to import tag processing modules and only import the dependencies for the ones they need.

## Installing a Tag Submodule

```shell
go get github.com/dubbikins/envy/v2/tags/{tag_module}@latest
```