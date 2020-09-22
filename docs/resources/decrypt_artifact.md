# Decrypt Artifact Resource

## Create

This resource, on creation, recursively searches a directory and decrypts encrypted artifacts (files with `.enc` extension ) and store it as a file with `.enc` extension stripped in the original files directory or under directory tree with `output_dir` as root.

**Note:** in case file with same name already exists in destination, it would remove that file. 

## Read

this resource takes no action on read.

## Update

This resource first runs `Delete` action and then runs `Create` action with updated values.

## Delete

On `Delete`, it re-encrypts the data that was decrypted. 
**Note:** decrypted data is stored in same dir as the plaintext form. in case there is a file with identical name, it would overwrite existing file.

## Example Usage

```hcl
```

## Argument Reference

- `path` - (Required) path of a directory or a single file for decryption. The plugin expects encrypted files to have `.enc` extension
- `nonce`- (Required) a random nonce generated at encryption time.
- `output_dir`- (Optional) root directory in which the decrypted files are stored. the plugin would try to replicate directory tree under input path . In case this is not passed, decrypted files would be stored in the same directory as their plain text source.