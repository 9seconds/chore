# chore

chore is an executor, environment, and management for scripts that help to
perform some routine operations.

Here is a short list of cool titles:

1. A Sometimes Better Management for Homebrew Scripts
2. Even Superheroes Wash Floors
3. Blazingly fast 🚀 script management

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
code unnecessarily. Even parsing of parameters could take more lines of code
than actual usage.

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

1. chore is a script runner that runs scripts organized under simple convention.
2. Each script could have an external configuration describing named parameters,
   its purpose etc.
3. Each script has its way to run other scripts from the same namespace so namespace
   is self contained and movable.
4. Each namespace has its secret vault, safe to move along with a namespace.
5. chore prepopulates tens of commonly used values like start time, ip address,
   machine id, git commit, geolocation and push them into scripts so they can
   immediately use them
6. maintains XDG directories for them and provides a correct temporary directory
   for each script run. You can think that all related directories are prepared
   beforehand, even temporary one.
7. It does not oblige your script to be a part of some frameworks or be implemented
   in any language

If you are interested, please proceed to [Wiki](https://github.com/9seconds/chore/wiki).
