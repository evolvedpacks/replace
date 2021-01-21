# replace(1)

## Usage

```
Replace strings in input streams.

Either you can provide mappings via command line arguments,
for example like following:

  cat data.txt | replace --map "replace this" --to "With this"

You can chain as many '--map' and '--to' bindings as you want
as long as the same ammount of mappings as of replacements is
provided. The first mappings is replaced with the frist
replacement and so on.
You can also provide a JSON file as mapping which looks like
following, for example:

  {
    "Replace this": "With this"
  }


replace v1.0.0
Usage: replace [--map MAP] [--to TO] [--mapfile MAPFILE] INPUT

Positional arguments:
  INPUT                  Input data (taken from STDIN when not provided)

Options:
  --map MAP, -m MAP      Values to be replaced
  --to TO, -t TO         Values to replace with
  --mapfile MAPFILE, -f MAPFILE
                         JSON file to read replacement mappings from
  --help, -h             display this help and exit
  --version              display version and exit
```

## Example

```
$ cat data/test.md \
    | replace --mapfile data/mappings.json --map "this" --to "that" \
    | tee data/result.md
```