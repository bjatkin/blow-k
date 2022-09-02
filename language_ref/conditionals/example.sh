#!/bin/bash

# Example 1
# a :: true
# if a {
#   $echo["success"]
# } else {
#   $echo["failure"]
# }
#
echo; echo "EXAMPLE 1"
a=true
if [[ $a == true ]]
then
    echo "success"
else
    echo "failure"
fi

# Example 2
# b :: 15
# c :: 20
# if b > 17 || c > 17 {
#   $echo["b or c is big", b, c]
# }
#
echo; echo "EXAMPLE 2"
b=15
c=20
if (( $b > 17 )) || (( $c > 17 ))
then
    echo "b or c is big" $b $c
fi

# Example 3
# d :: "yes"
# if d != "no" {
#   $echo["not no", d]
# }
#
echo; echo "EXAMPLE 3"
d="yes"
if [[ $d != "no" ]]
then
    echo "not no" $d
fi