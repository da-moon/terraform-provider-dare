# Encrypt Artifact Resource

## Create

This resource, on creation, recursively searches a directory and encrypts artifacts annd stores it as a file with `.enc` extension at rest either in the same directory as the original file or under a directory tree with `output_dir` as root.

**Note:** in case file with same name already exists in destination, it would remove that file. 

## Read

this resource takes no action on read.

## Update

This resource first runs `Delete` action and then runs `Create` action with updated values.

## Delete

this resource deletes all encrypted artifacts on `Delete`.

## Example Usage

```hcl
```

## Argument Reference

- `path` - (Required) path of a directory or a single file for encryption
- `output_dir`- (Optional) root directory in which the encrypted files are stored. the plugin would try to replicate directory tree under input path . In case this is not passed, encrypted files would be stored in the same directory as their plain text source. all encrypted artifacts would have `.enc` extension.  
- `regex`- (Optional) in case input `path` is a directory, `regex can be used to recursively search all subdirectories for files that meets the requested pattern.
