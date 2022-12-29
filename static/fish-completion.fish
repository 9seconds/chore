#!/usr/bin/env fish


set -l CHORE_COMMANDS run edit-script edit-config show


function __chore_require_command_completion
  if test (count (commandline -poc)) -eq 1
    return 0
  end

  return 1
end


function __chore_require_namespace_completion
  set -l args (commandline -poc)

  if test (count $args) -ne 2
    return 1
  end

  if contains -- $args[2] $CHORE_COMMANDS
    return 1
  end

  return 0
end


function __chore_require_script_completion
  set -l args (commandline -poc)

  if test (count $args) -ne 3
    return 1
  end

  if contains -- $args[2] $CHORE_COMMANDS
    return 1
  end

  return 0
end


function __chore_require_arguments_completion
  set -l args (commandline -poc)

  if test (count $args) -lt 4
    return 1
  end

  if contains -- $args[2] $CHORE_COMMANDS
    return 1
  end

  return 0
end


function __chore_require_script_completion
  set -l args (commandline -poc)

  if test (count $args) -ne 3
    return 1
  end

  if contains -- $args[2] $CHORE_COMMANDS
    return 1
  end

  return 0
end


function __chore_namespace_completion
  chore show
end


function __chore_script_completion
  set -l args (commandline -poc)

  chore show $args[3]
end


function __chore_arguments_completion
  set -l args (commandline -poc)

  chore fish-completion $args[3..(count $args)]
end


complete -x  -c chore -n '__fish_use_subcommand'                -s h -l help    -d 'Show help'
complete -x  -c chore -n '__fish_use_subcommand'                -s d -l debug   -d 'Run in debug mode'
complete -x  -c chore -n '__fish_use_subcommand'                -s V -l version -d 'Show version'
complete -x  -c chore -n "__chore_require_command_completion"   -a run          -d 'Run chore script'
complete -x  -c chore -n "__chore_require_command_completion"   -a edit-script  -d 'Edit chore script'
complete -x  -c chore -n "__chore_require_command_completion"   -a edit-config  -d 'Edit chore script config'
complete -x  -c chore -n "__chore_require_command_completion"   -a show         -d 'Show details on namespaces or scripts'
complete -x  -c chore -n "__chore_require_namespace_completion" -a "(__chore_namespace_completion)"
complete -x  -c chore -n "__chore_require_script_completion"    -a "(__chore_script_completion)"
complete -Fr -c chore -n "__chore_require_arguments_completion" -a "(__chore_arguments_completion)"
