### File deleter processor
The deleter processor marks files for deletion by the internal deletion output. This allows files with certain extensions, empty directories, or that match specified regexps to be deleted after Pachinko has finished sorting.

#### Configuration
The default deleter plugin configuration is:
```yaml
- categories: []
  directories: true
  extensions:
  - 7z
  - gz
  - gzip
  - rar
  - tar
  - zip
  - bmp
  - gif
  - heic
  - jpeg
  - jpg
  - png
  - tiff
  - info
  - nfo
  - txt
  - website
  matchers: []
  name: deleter
```

||||
|-|-|-|
|`categories`|`[]string`|[unimplemented] list of categories of file such as text or archive to remove.|
|`directories`|`bool`|whether to remove directories. Even if true, only *empty* dirs will be removed.|
|`extensions`|`[]string` | list of file extensions to remove.|
|`matchers`|`[]string` | regexps to match files to remove.|
