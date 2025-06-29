---
title:          Go's defers should not return values
date:           2025-06-29
author:         Tom Wiesing 
authorLink:     https://tkw01536.de

description:    Why 'defer' statements in go can lead to silencing errors.

draft:          true
---

Consider an off-the-shelf go program which wants to marshal some data into a file.
It might implement this functionality using code like the following:

```go
func (data D) Export(name string) error {
    file, err := os.Create(name, 0644) // create the file
    if err != nil {
        return 0, fmt.Errorf("failed to create file: %w", err)
    }
    defer file.Close() // close the file on return

    if _, err := data.WriteTo(file); err != nil {
        return fmt.Errorf("failed to write to file: %w", err)
    }
    return nil
}
```

This relies on a second method `WriteTo` to perform the actual marshaling.
While the precise definition is not important, it can be assumed to implement the [io.WriterTo](https://pkg.go.dev/io#WriterTo) interface.
It might make various (direct or indirect) calls to [io.Writer.Write](https://pkg.go.dev/io#Writer).

```go
func (data D) WriteTo(w io.Writer) (int64, error) {
    // ...
}
```

<!-- Next: what is the problem here? -->