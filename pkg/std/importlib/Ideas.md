# Ideas for the import of the language

- There should be the possibility of importing a script directly from the working directory
- There should be possible to differentiate from scripts and modules
- There should be a way to expand the symbols in a module to the importer's table
- There should be the possibility to import scripts from a `sites-packages`
- There should be a way to include resources in the modules
- There should be a way to standardized modules
- It should be possible to have multiple version of the same module

## Modules information based in jsons

```json
{
  "name": "TheModuleName",
  "version": "0.0.0",
  "resources": "path/to/the/relative/resources/directory",
  "entry-script": "path/to/the/relative/entry/script/of/the/module",
  "dependencies": [
    "dependency_1",
    "dependency_2"
  ]
}
```

Dependencies will be also installed the same way modules are.

## Modules resource access

This is the way to relatively acceded the resources specified in the resources' entry of the **`JSON`**

### Resources read

This should return an object with a file like interface but read only:

- Read
- Seek
- Close

```ruby
resource = get_resource("relative/reource/path")
```

### Resource full path

This should return a string with the full path of a resource in the filesystem object handler. This way the user can do
what ever he wants with the resource.

```ruby
resoutce_path = get_resource_path("relative/reource/path")
```

## Imports code

### Relative imports

This kind of imports should always be relative to the script that execute them.

```ruby
script = import("path/to/script.pm")
```

### Module imports

This kind of imports should always be made from the `sites-packages` and are always of modules

```ruby
script = import_module("module_name") # Import the latest version of module.
script2 = import_module("module_name@0.0.0") # Import an specific version of a module.
```

## Command Line

This is not part of the library itself but part of the CLI.

### Sites packages set

By this way the module search can be set to a different environment

```shell
plasma -env path/to/sites/packages
```

If the flag is not set, it will search in the current directory for a `sites-packages` directory, if there is not one,
it will create the directory.

### Module initialization

```shell
plasma module init MODULE_NAME
```

### Module installation

This should be executed inside a module, or should be given to it the path the module.

When module path is a GitHub repository, it will be cloned and then manipulated locally.

```shell
plasma module install [MODULE_PATH]
```

### Module deletion

This function should be able to delete a module from the `sites-packages`.

```shell
plasma module uninstall MODULE_NAME [MODULE_VERSION]
```

### Module update

This should delete all the old versions of the module and reinstall the new one.

```shell
plasma module update [MODULE_PATH]
```

## Sites packages structure

```
site-packages
│
└───module
    └───0.0.0
        │   file.pm
        │   module.json
        │
        └───resources
```

Every child folder in the `sites-packages` is a name space reserved to the module.

The child folders of the reserved namespace are the versions of the module installed.

The file `module.json` is the **`JSON`** with all the information of the module.