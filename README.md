# chore

chore is an executor, environment, and management for scripts that help to
perform some routine operations.

Here is a short list of cool titles:

1. A Sometimes Better Management for Homebrew Scripts
2. Even Superheroes Wash Floors
3. Blazingly fast 🚀 script management

And here is the elevator pitch.

## Elevator pitch

I've been doing software engineering for a while now and have learned a few things about myself.
For example, I used to abuse 'Ctrl+R' in a terminal to perform a series of actions on a regular basis.
Even if all they do is execute another command, I eventually turn them into scripts.
I followed the same principle in a couple of very different companies in the past and will continue to 
do so in the future: this isn't about projects, it's just how I work. For dull ceremonies, I like to have wrappers.

The majority of such scripts are quick-and-dirty dumps that only expect happy paths and 
do not handle errors. They typically lack tests (do you test your Bash scripts?), primitive error 
handling, and incoming parameter validation. But this isn't because I'm a slacker: I like having such 
things; they just bloat code unnecessarily, and sometimes I don't want to pollute a one-liner with a huge 
boilerplate of getopts and the like. Even parameter parsing may require more lines of code than actual usage. 

Imagine, you have a script:

```shell
$ plog https://lalala.blablabla 1 2
```

How would you ensure that the first parameter is always a valid URL to _lalala_ host, the second is 
limited to _1..10_, and the third is a string that begins with a digit? 

On the one hand, you don't need it for a quick'n'dirty script.
But there are times when you crave it.

If you do it in a script, the code that manages these parameters will most likely be bloated 
boilerplate that dilutes its purpose.

What if you need to keep some secrets? How would you keep it? Is it possible to transfer a directory 
containing such scripts from one machine to another?

Usually, scripts do repetitive tasks: get current data, extract a git commit hash, etc. What 
if an execution environment gives you some data?

This is what a chore entails.

1. All scripts are namespaced, and each script has its own name.
2. Each namespace should be linkable, distributable and probably stored in a VCS aside.
3. Because each script can run other scripts from the same namespace, the namespace is 
   self-contained and movable.
5. An external configuration describing named parameters, their purpose, and validation 
   strategies could be included in each script.
7. Each namespace has its own secret vault, which can be moved with the namespace.
8. Typically, chore scripts require common information such as git commit hashes, 
   current branch names, and so on. Harness extracts them and provides environment 
   variables.
9. It also helps to keep and maintain a set of XDG-styled directories for each script. Even
   temporary directory is prepopulated beforehand.
10. It does not require your script to be a part of any frameworks or to be written in any language.
    Actually, you can put any binary there. Even a whole chore binary, why not.

In case if you are interested, here is a wiki: https://github.com/9seconds/chore/wiki.
