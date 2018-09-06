# copy

A better copy command, with a progress bar, concurrent copying and the ability to use host aliases for username+host+port when coping over ssh.

Should always ask before overwriting files. `pscp` can not ask before overwriting.

# THIS IS A WORK IN PROGRESS AND NOT NEARLY DONE

* Use `-a` to add a username@hostname:port alias on the form: `alias=username@hostname:port`.
* Use `-r` to remove a remote host alias.
* Use `-l` to list host aliases.

Takes two argument, the file to copy from and the file to copy to.
