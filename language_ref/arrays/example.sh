#!/bin/bash

# Example 1
# a:[]string
# this is compiled away and only used for syntax validation,
# in bash variables do not need a type and do not need to be defined before use

# Example 2
# a:[]string:["hi", "hello"] or a :: ["hi", "hello"]
echo; echo "EXAMPLE 2"
a=( "hi" "hello" )
echo ${a[@]}

# Example 3
# b:[]int:[1, 2, 3] or b :: [1, 2, 3]
echo; echo "EXAMPLE 3"
b=( 1 2 3 )
echo ${b[@]}

# Example 4
# c:[]bool:[true, false, true] or b :: [true, false, true]
echo; echo "EXAMPLE 4"
c=( true false true )
echo ${c[@]}

# Example 5
# a:[][]string:[ ["hi", "world"], ["hello", "world"] ]
# nested arrays are not possible in bash so we need to use associatve arrays instead
echo; echo "EXAMPLE 5"
declare -A d
d=( [0,0]="hi" [0,1]="world" [1,0]="hello" [1,1]="world" )
echo ${d[@]}

# Example 6
# e :: ["one", "two"]
# e<-"three"
# e<-["four", "five"]
echo; echo "EXAMPLE 6"
e=( "one" "two" )
e+=( "three" )
e+=( "four" "five" )
echo ${e[@]}

# Example 7
# f :: [1, 2, 3, 4, 5]
# g :: a->
echo; echo "EXAMPLE 7"
f=( 1 2 3 4 5 )
f_last=$(( ${#f[@]} - 1 ))
g=${f[$f_last]}
f=( ${f[@]:0:$f_last} )
echo ${f[@]}
echo $g

# Example 8
# h :: [ true, false, false ]
# i :: h[:1]
echo; echo "EXAMPLE 8"
h=( true false true )
i=( ${h[@]:0:1} )
echo ${h[@]}
echo ${i[@]}

# Example 9
# j :: [1, 2, 3]
# j[1] += 1 or j[1]++
# j[0] += j[2]
echo; echo "EXAMPLE 9"
j=( 1 2 3 )
j[1]=$(( j[1] + 1 ))
j[0]=$(( ${j[0]} + ${j[2]} ))
echo ${j[@]}