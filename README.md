<!-- start title -->

# GitHub Action: Duplicates Checker

<!-- end title -->
<!-- start description -->

## Duplicate-text Detector

The duplicate-text detector is a Github Action that provides a configurable utility for detecting duplicate text in files after a commit. 

It comes in handy to block hard-to-find bugs with serious adverse impacts that come from duplicating certain content from one file to the next -- for example, merging the same state key for two different terraform contexts. 


## Running Locally
You can modify the make file in the root folder and just run 

```
make [command]
```

<!-- end description -->

<!-- start usage -->
## Usage

```yaml
- uses: heap/duplicate-checker-action@main
  with:
    # A name for this check
    check_name: ""

    # Checks for duplicates will be evaluated in the scope of this directory
    directory_scope: ""

    # Checks for duplicates will be done only for files with this extension
    check_file_extension: ""

    # The regular expression to use in checking for content
    content_regex: ""

    # This threshold limits the number of files to scan. It's a safeguard against
    # locking up your build by scanning too many files. Default is 500
    # Default: 500
    max_files_to_process: ""

    # This threshold limits the size of files to scan, to prevent memory overload.
    # Default is 200Kb
    # Default: 200000
    max_file_size_bytes: ""

    # A comma delimited list of extensions to skip (extensions can include the dot
    # prefix or not). This configuration allows skipping over files with certain
    # extensions
    # Default:
    excluded_extensions: ""

    # A comma delimited list of filenames to skip. This configuration allows skipping
    # certain file names. Useful if you use the same settings in multiple environment
    # configuration files.
    # Default:
    ignore_files: ""

    # A comma delimited list of paths to skip. This configuration allows skipping
    # certain paths. Useful if you have generated paths that have copies of files you
    # do not want to check for duplicates.
    ignore_paths_containing: ""

    # The message to output if duplicates are found
    error_message: ""
```

<!-- end usage -->
<!-- start inputs -->

| **Input**                     | **Description**                                                                                                                                                                                  | **Default** | **Required** |
| ----------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | ----------- | ------------ |
| **`check_name`**              | A name for this check                                                                                                                                                                            |             | **false**    |
| **`directory_scope`**         | Checks for duplicates will be evaluated in the scope of this directory                                                                                                                           |             | **true**     |
| **`check_file_extension`**    | Checks for duplicates will be done only for files with this extension                                                                                                                            |             | **true**     |
| **`content_regex`**           | The regular expression to use in checking for content                                                                                                                                            |             | **true**     |
| **`max_files_to_process`**    | This threshold limits the number of files to scan. It's a safeguard against locking up your build by scanning too many files. Default is 500                                                     | `500`       | **false**    |
| **`max_file_size_bytes`**     | This threshold limits the size of files to scan, to prevent memory overload. Default is 200Kb                                                                                                    | `200000`    | **false**    |
| **`excluded_extensions`**     | A comma delimited list of extensions to skip (extensions can include the dot prefix or not). This configuration allows skipping over files with certain extensions                               |             | **false**    |
| **`ignore_files`**            | A comma delimited list of filenames to skip. This configuration allows skipping certain file names. Useful if you use the same settings in multiple environment configuration files.             |             | **false**    |
| **`ignore_paths_containing`** | A comma delimited list of paths to skip. This configuration allows skipping certain paths. Useful if you have generated paths that have copies of files you do not want to check for duplicates. |             | **true**     |
| **`error_message`**           | The message to output if duplicates are found                                                                                                                                                    |             | **true**     |

<!-- end inputs -->
<!-- start outputs -->

| \***\*Output\*\*** | \***\*Description\*\***                                       | \***\*Default\*\*** | \***\*Required\*\*** |
| ------------------ | ------------------------------------------------------------- | ------------------- | -------------------- |
| `has_duplicates`   | True or false depending on if the check found duplicates      | undefined           | undefined            |
| `result_title`     | Title caption of result summary                               | undefined           | undefined            |
| `result`           | The output after evaluating all checks                        | undefined           | undefined            |
| `result_escaped`   | The output after evaluating all checks in escaped json format | undefined           | undefined            |

<!-- end outputs -->
