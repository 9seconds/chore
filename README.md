# chore

chore is an executor, environment, and management for scripts that help to
perform some routine operations.

Here is a short list of cool titles:

1. A Sometimes Better Management for Homebrew Scripts
2. Even Superheroes Wash Floors
3. Blazingly fast 🚀 script management

And here is the elevator pitch.

## Elevator pitch

I've been doing software engineering for some time and understood a couple of things about myself. One example: I used to abuse `Ctrl+R` in a terminal and to routinely perform a
set of actions. Sooner or later I turn them into scripts even if all they do is
execute another command. I followed the same princinple in a couple of very different companies in
the past and will do it in the future: this is not about projects, this is how I do
things. I like to have wrappers for boring ceremonies.

Most of such scripts are quick-and-dirty dumps that expect only happy paths and
do not manage errors. They usually do not have tests (do you test your Bash
scripts?), primitive error management, or incoming parameter validation. But
this is not because I am sloppy: I like to have such things, they just bloat a
code unnecessarily and sometimes I just do not want to pollute one-liner with a huge boilerplate of getopts etc. Even parsing of parameters could take more lines of code
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
be a bloated boilerplate that will dilute its purpose.

What if you want to manage some secrets? How would you store it? Can you move a
directory with such scripts from one machine to another?

Usually, scripts do repetitive tasks: get current data, extract a git commit
hash, etc. What if some execution environment provides you with some data?

This is what chore is. This is just a harness expressed in a form of external tool that does exactly what I need:

1. All scripts are organized in namespaced, each script has its name.
2. Each namespace should be linkable, distributable and probably stored in a VCS aside.
3. Each script has its way to run other scripts from the same namespace so namespace
   is self contained and movable.
4. Each script could have an external configuration describing named parameters,
   their purpose and validation strategies.
5. Each namespace has its own secret vault, safe to move along with a namespace.
6. Usually chore scripts require common stuff like git commit hashes, current branch names etc. 
   Harness extracts them and supply a script with environment variables.
8. It helps to maintain a set of directories in XDG styles. Scripts are executed having
   everything prepared, even temporary directories.
10. It does not oblige your script to be a part of some frameworks or be implemented
   in any language

If you are interested, please proceed to [Wiki](https://github.com/9seconds/chore/wiki).
