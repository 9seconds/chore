# chore

chore is an executor, environment, and management for scripts that help to
perform some routine operations.

Here is a list of cool titles:

1. A Sometimes Better Management for Homebrew Scripts
2. Even Superheroes Wash Floors
3. Management of a set of self-contained scripts and their secrets.
4. Blazingly fast 🚀 script management

And here is the elevator pitch.

## Elevator pitch

I've been doing software engineering for almost 15 years and understood some
things about myself: I tend to abuse `Ctrl+R` in a terminal and routinely perform a
set of actions. Sooner or later I turn them into scripts even if all they do is
execute another command. I did it in a couple of very different companies in
the past and will do it in the future: this is not about projects, this is how I do
things. I like to have wrappers for boring ceremonies.

Most of such scripts are quick-and-dirty dumps that expect only happy paths and
do not manage errors. They usually do not have tests (do you test your Bash
scripts?), primitive error management, or incoming parameter validation. But
this is not because I am sloppy: I like to have such things, they just bloat a
code unnecessarily.

For example, if I have a script

```shell
$ plog https://lalala.blablabla 1 2
```

How would you validate that the first parameter is always a correct URL to lalala
host, 2nd is range limited to 1..10 and 3rd is a string that starts with a digit?

On one hand, for a fast'n'dirty script, you do not need it. But sometimes you
want it.

If you do it within a script, the code managing these parameters will probably
be a bloated boilerplate that will dilute its purpose of it.

What if you want to manage some secrets? How would you store it? Can you move a
directory with such scripts from one machine to another?

Usually, scripts do repetitive tasks: get current data, extract a git commit
hash, etc. What if some execution environment provides you with some data?

This is what chore is. It is an attempt to express how I manage
these tasks:

1. chore is a script runner that knows where your scripts are placed.
2. Each script has its name and belongs to some namespace.
3. namespaces are movable between different machines.
4. Each namespace has its secret vault.
5. Each script has a defined set of named parameters with a description of how
   to validate them.
6. chore generates many environment variables with data commonly required
   by such scripts (IP addresses, dates and times, commit hashes, machine ids
   , etc) so all you need is just to get it from the environment
7. maintains XDG directories for them and provides a correct temporary directory
   for each script run
8. It does not oblige your script to be a part of some frameworks or be implemented
   in any language

Please proceed to Wiki if you are interested.
