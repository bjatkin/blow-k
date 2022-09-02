#!/bin/bash

# Example 1
# v2:struct: < x:int, y:int > or v2:struct: < x, y:int > or v2 :: < x, y:int >
# a:v2: {
#   x:5
#   y:10 
# }
#
# b:v2: {
#   x:5
#   y:10 
# }
echo; echo "EXAMPLE 1"
declare -A v2_ret
function v2 () {
    v2_ret=( [x]=$1 [y]=$2 )
}

declare -A a
v2 5 10
a=( [x]=${v2_ret[x]} [y]=${v2_ret[y]} )
echo ${a[@]}

declare -A b
v2 80 3
b=( [x]=${v2_ret[x]} [y]=${v2_ret[y]} )
echo ${b[@]}

# Inline Version
declare -A a_in
a_in=( [x]=5 [y]=10 )

# Example 2
# v2 :: < x, y:int >
# square :: < one, two:v2 >
# c:square: {
#   one: {
#     x: 1
#     y: 2
#   }
#   two: {
#     x: 3
#     y: 4
#   }
# }
echo; echo "EXAMPLE 2"
declare -A v2_ret
function v2 () {
    v2_ret=( [x]=$1 [y]=$2 )
}

declare -A square_ret
function square () {
    v2 $1 $2
    square_ret=( [one.x]=${v2_ret[x]} [one.y]=${v2_ret[y]} )

    v2 $3 $4
    square_ret+=( [two.x]=${v2_ret[x]} [two.y]=${v2_ret[y]} )
}

declare -A c
square 1 2 3 4
c=( [one.x]=${square_ret[one.x]} [one.y]=${square_ret[one.y]} [two.x]=${square_ret[two.x]} [two.y]=${square_ret[two.y]} )
echo ${c[@]}

# Inline Version
declare -A c_in
c_in=( [one.x]=1 [one.y]=2 [two.x]=3 [two.y]=4 )


# Example 3
# struct with embedded array
# greet :: < lang:string, greet:[]string, goodbye:[]string >
#
# english:greet: {
#   lang:    "english"
#   greet:   ["hi", "hello", "heya"]
#   goodbye: ["bye", "seeya"]
# }
#
# english.greet <- "sup"
# english.goodbye->
#
echo; echo "EXAMPLE 3"
declare -A greet_ret
function greet () {
    args=("$@")
    greet_ret=( [lang]=$1 [greet.len]=$2 [goodbye.len]=$4 )

    for (( i=0; i<$2; i++ ))
    do
        greet_ret+=( [greet.$i]=${args[ $(( $i + $3 )) ]} )
    done

    for (( i=0; i<$4; i++ ))
    do
        greet_ret+=( [goodbye.$i]=${args[ $(( $i + $5 )) ]})
    done

}

declare -A english
greet "english" 3 5 2 8 "hi" "hello" "heya" "bye" "seeya"
english=( [lang]=${greet_ret[lang]} [greet.len]=3 [goodbye.len]=2 )
for (( i=0; i<3; i++ ))
do
    english+=( [greet.$i]=${greet_ret[greet.$i]} )
done

for (( i=0; i<2; i++ ))
do
    english+=( [goodbye.$i]=${greet_ret[goodbye.$i]})
done

echo ${english[@]}

english[greet.${english[greet.len]}]="sup"
english[greet.len]=$(( ${english[greet.len]} + 1 ))
echo ${english[@]}

goodbye_last=$(( ${english[goodbye.len]} - 1))
unset english[goodbye.$goodbye_last]
english[goodbye.len]=$goodbye_last
echo ${english[@]}

# Inline Version (this is much simpler and should be used agressively)
declare -A english_in
english=( [lang]="english" [greet.len]=3 [greet.0]="hi" [greet.1]="hello" [greet.2]="heya" [goodbye.len]=2 [goodbye.0]="bye" [goodbye.1]="seeya" )

# Example 4 
# add defaults to the struct def
# v4 :: < x, y, z, w:int:1 >
# p1:v4: { x: 5, y: 10 }
echo; echo "EXAMPLE 4"
declare -A v4_ret
function v4 () {
    v4_ret=( [x]=$1 [y]=$2 [z]=$4 [w]=$4 )
}

declare -A p1
v4 5 10 0 1
p1=( [x]=${v4_ret[x]} [y]=${v4_ret[y]} [z]=${v4_ret[z]} [w]=${v4_ret[w]} )
echo ${p1[@]}

# Inline Version
declare -A p1_in
p1_in=( [x]=5 [y]=10 [z]=0 [w]=1 )

# Example 5
# struct with methods
# v2 :: < x, y:int
#   add:(b:v2): {
#       me.x = me.x + b.x
#       me.y = me.y + b.y
#   }
# >
#
# a:v2: {x: 4, y: 5}
# b:v2: {x: 6, y: 1}
# a.add(b)
#
echo; echo "EXAMPLE 5"
declare -A v2_ret
function v2 () {
    v2_ret=( [x]=$1 [y]=$2 )
}

function v2_add () {
    v2_add_ret=( [x]=$(( $1 + $2 )) [y]=$(( $3 + $4 )) )
}

declare -A a
v2 4 5
a=( [x]=${v2_ret[x]} [y]=${v2_ret[y]} )
echo ${a[@]}

declare -A b
v2 6 1
b=( [x]=${v2_ret[x]} [y]=${v2_ret[y]} )
echo ${b[@]}

v2_add ${a[x]} ${b[x]} ${a[y]} ${b[y]}
a=( [x]=${v2_add_ret[x]} [y]=${v2_add_ret[y]} )
echo ${a[@]}

# Inline Version
a=( [x]=4 [y]=5 )
echo ${a[@]}

b=( [x]=6 [y]=1 )
echo ${b[@]}

a_tmp=( [x]=$(( ${a[x]}+${b[x]} )) [y]=$(( ${a[y]}+${b[y]} )))
a=( [x]=${a_tmp[x]} [y]=${a_tmp[y]} )
echo ${a[@]}