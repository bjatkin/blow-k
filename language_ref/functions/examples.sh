#!/bin/bash

# Example 1
# main :(args []string): {
#   $echo["hello", args[0]]
# }
#
# main functions get turned directly into a script
echo; echo "EXAMPLE 1"
args=($@)
echo "hello" ${args[0]}

# Example 2
# run :(do:fn(string, string)<string>)<string>: {
#   return do("hi", "hello")
# }
# 
# c :: "test"
#
# combine :: run((a, b:string)<string> {
#   return a + b + c
# })
#
# c : "sup"
# str :: run(combine)
#
echo; echo "EXAMPLE 2"
function run () {
    $1 "hi" "hello"
    run_ret=$( eval "echo \$$2" )
}

c="test"

function combine () {
    combine_ret=$1$2$c
}

c="sup"
run combine combine_ret
str=$run_ret
echo $str

# Inline Version
c="test"
c="sup"
str="hi""hello"$c
echo $str

# Example 3
# main :(args []string): {
#   fn : test()
#   $echo[fn()]
#
#   cap1 : "goodbye"
#   #echo[fn()]
# }
#
# cap1 :: "hello"
#
# test :(): <()<string>> {
#   cap2 :: "world"
#   return :():<string> {
#       return cap1 + " " + cap2
#   }
# }
#
echo; echo "EXAMPLE 3"
cap1="hello"

function test () {
    cap2="guys"
    test_ret=( anon_1 $cap1 $cap2 )
}

function anon_1 () {
    anon_1_ret=$1" "$2
}

test 
fn=( ${test_ret[@]} )
fn[1]=$cap1
${fn[@]}
echo $anon_1_ret

cap1="goodbye"
fn[1]=$cap1
${fn[@]}
echo $anon_1_ret

# Example 4
# ask :(question string)|resp|<string>: {
#   $echo[question+":"]
#   return resp
# }
#
# name :: ask("what is your name?")
# $echo["hello", name]
#
echo; echo "EXAMPLE 4"
function ask () {
    echo $1":"
    read resp
    ask_ret=$resp
}

ask "what is your name?"
name=$ask_ret
echo "hello" $name