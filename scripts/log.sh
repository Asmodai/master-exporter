#!/bin/bash
# -*- Mode: Shell-script -*-
#
# log.sh --- Log file pretty-printer.
#
# Copyright (c) 2022 Paul Ward <asmodai@gmail.com>
#
# Author:     Paul Ward <asmodai@gmail.com>
# Maintainer: Paul Ward <asmodai@gmail.com>
# Created:    11 Mar 2022 02:56:19
#
# {{{ License:
#
# This program is free software: you can redistribute it
# and/or modify it under the terms of the GNU General Public
# License as published by the Free Software Foundation,
# either version 3 of the License, or (at your option) any
# later version.
#
# This program is distributed in the hope that it will be
# useful, but WITHOUT ANY  WARRANTY; without even the implied
# warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR
# PURPOSE.  See the GNU General Public License for more
# details.
#
# You should have received a copy of the GNU General Public
# License along with this program.  If not, see
# <http://www.gnu.org/licenses/>.
#
# }}}
# {{{ Commentary:
#
# }}}

if [ $# -ne 1 ]
then
    echo "${0} <logfile>"
fi

if [ ! -f ${1} ]
then
    echo "Uhhh, ${1}, like, isn't a file or something, dude."
    exit 1
fi

/usr/bin/jq -c '.ts |= (strftime("%Y-%m-%dT%H:%M:%Sz"))' ${1}

# log.sh ends here.
