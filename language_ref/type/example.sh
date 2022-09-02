#!/bin/bash

# Example 1
# a:string
# b:string: "hello"
echo; echo "EXAMPLE 1"
a=""
b="hello"
echo $a
echo $b

# Example 2
# c :: "hi"
# d :: "world"
# e :: a + b
echo; echo "EXAMPLE 2"
c="hi"
d="world"
e=$c$d
echo $e

# Example 3
# f :: 1
# g :: 2 + 3
# h :: f + g + 4 + g
echo; echo "EXAMPLE 3"
f=1
g=2
h=$(( $f + $g + 4 + $g ))
echo $h

# Example 4
# i :: true
# j :: false
echo; echo "EXAMPLE 4"
i=true
j=false
echo $i
echo $j