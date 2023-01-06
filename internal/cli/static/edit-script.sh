#!/bin/bash
# vim: set ft=bash sw=2 ts=2 sts=2 et:
set -eu -o pipefail

[[ ${{"{"}}{{ .DebugVar }}:-{{"}"}} = "{{ .DebugEnabled }}" ]] && set -x


