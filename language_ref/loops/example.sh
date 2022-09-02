#!/bin/bash

# Example 1
# stealling the idea of a single loop keyword from go
# but using the 'loop' keyword instead of 'for'
#
# loop i::0, i<10, i++ {
#   $echo["loop", i]
# }
#
echo; echo "EXAMPLE 1"
for (( i=0; i<10; i++ )) {
    echo "loop" $i
}

# Example 2
# 
# a:int
# b :: 3
# loop a < b {
#   $echo["a", a, "b", b]
#   a++
# }
#
echo; echo "EXAMPLE 2"
a=0
b=3
while (( $a < $b ))
do
    echo "a" $a "b" $b
    a=$(( $a + 1 ))
done

# Example 3
#
# c:int
# d :: 3
# loop c <= d, c++ {
#   $echo["c", c, "d", d]
# }
echo; echo "EXAMPLE 3"
c=0
d=3
while (( $c <= $d ))
do
    echo "c" $c "d" $d
    c=$(( $c + 1 ))
done