# socketbro

## Description

This daemon does nothing but establish unix sockets (i.e., UNIX domain sockets) and close them up when it receives SIGINT or SIGTERM.  It will try to remove any existing socket file before binding (in case it crashed or something else happened), but it will raise an error if you specify a non-socket location (directory, regular file, device file, etc.).

I wrote this because Go programs can't fork() like C programs.  In a normal setup, your master process binds the unix socket and fork()s workers which read from it.  The process which binds the socket must remain running for the socket to be usable.  This makes it impossible to have N workers (in Go) reading from the same shared unix socket without a third party maintaining it... thus, socketbro was born.

Created sockets are chmod 755.  To make things easy, just run this as whatever user will be writing to the socket e.g., if you're using nginx, you'll run this as the same user nginx runs as (probably "nobody").  Any other user can read from it.

## Feedback

Please feel free to send pull requests or file issues.
