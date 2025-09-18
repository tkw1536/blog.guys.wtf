---
title:          Go's defers should not be allowed to return values
date:           2025-09-21
author:         Tom Wiesing 
authorLink:     https://tkw01536.de

description:    Why 'defer' statements in go can lead to silencing errors, and why they should not be allowed to return values. 

draft:          true
---

Consider an off-the-shelf go program which wants to marshal some data into a file.
It might implement this functionality using code like the following:

```go
func (d Data) Export(name string) error {
    file, err := os.Create(name, 0644)
    if err != nil {
        return 0, fmt.Errorf("failed to create file: %w", err)
    }
    defer file.Close()

    if _, err := d.WriteTo(d, file); err != nil {
        return fmt.Errorf("failed to write to file: %w", err)
    }
    return nil
}
```

In order, this method:

- creates a new file;
- `defer`s[^1] closing the file once the method returns; and finally
- delegates the actual writing to a `WriteTo` method.

At first glance, this function looks correct - it even performs proper error checking. 
But there is a subtle bug here: 
Data is not actually flushed to disk before the file is properly closed using the `(*os.File).Close` method. 
Important in this case: Calling `Close` can fail.
In this case the go implementation returns a non-nil `error`.
This error is not checked by our function above, possibly resulting in a `nil` return value despite the data not actually having been written. 
This is a well-established problem, see for example [^2] and [^3]. 

Linters such as `errcheck` [^4] can check if your code has checked these errors. 
But that only partially addresses the problem. 
What if instead of returning an error the function you call returns a boolean?
What if a function expects the caller to perform some more advanced error checking?
That is not easily caught by a linter. 

In go, unused variables are an error. 
The go language FAQ [^5] says:
> The presence of an unused variable may indicate a bug [...]. For these reasons, Go refuses to compile programs with unused variables [...], trading short-term convenience for long-term build speed and program clarity. 

But is an ignored return value from a `defer`ed function not the exact same thing?
I would argue that such a call should also immediately be a compiler error. 


[^1]: If you're not familiar with go, a `defer file.Close()` is effectively equivalent to inserting `file.Close()` immediately before any of the `return` statements.
This doesn't take into account `panic()` handling and receiver/argument evaluation order.
See the go specification on [defer statement](https://go.dev/ref/spec#Defer_statements) for details. 

[^2]: Joe Shaw. 12th June 2017. [Don't defer Close() on writable files](https://www.joeshaw.org/dont-defer-close-on-writable-files/). 
[^3]: SoByte. 15th Jan 2022. [Simply defer file.Close() is probably a misuse](https://www.sobyte.net/post/2022-01/golang-defer-file-close/).
[^4]: https://github.com/kisielk/errcheck
[^5]: https://go.dev/doc/faq#unused_variables_and_imports

<!-- spellchecker:words Errorf -->
