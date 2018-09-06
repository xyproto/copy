# copy

# WORK IN PROGRESS, early stages

Copy a file locally or over ssh, and ask before overwriting

* Use `-a` to add a username@hostname:port alias on the form: `alias=username@hostname:port`
* Use `-r` to remove a remote host alias.

Takes two argument, the file to copy from and the file to copy to.

If the filename starts with an exclamation mark, it is interpreted as a remote host alias instead.
