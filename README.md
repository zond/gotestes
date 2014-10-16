gotestes
========

A tool to simplify running just the go tests you WANT to run.

Say, for example, that you have a test suite with 200 tests that takes 30 minutes to run. Now, your test suite explodes after 20 minutes, and 40 tests failed.

What to do? Run each one individually using `-run=TestXXX`? Gaawd how irksome.

No! Now we have a solution!

`go test -v -timeout=30m -run=$(gotestes -from=TestXXX -to=TestYYY)`

Tam tam!

