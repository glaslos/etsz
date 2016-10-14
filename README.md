# etsz
Embedded golang time series database

Use with `import "github.com/glaslos/etsz"`

Create a new embedded db with `edb := etsz.New()`

Create entry: `edb.Insert(1.0, "test")` where "test" is the series name.


Read data: `edb.Read("test")`
Example output:
```JSON
[
  {
    "name": "test",
    "shards": [
      {
        "range": "20161014T14",
        "data": [
          {
            "date": 1476447064,
            "value": 1
}]}]}]
```

Series are split in shards with a currently hard-coded resolution of one hour (yolo).

`range` is the shard named by it's timestamp in hour resolution (yolo).

A shard contains data points consisting of unix timestamps (second resolution (yolo)) and values.

Shards are deleted on inserts if the shard is from the previous day (yolo).
