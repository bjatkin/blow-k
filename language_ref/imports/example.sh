#!/bin/bash

# Example 1
# import sed as test
if [[ -z "$( which which )" ]]; then
    exit 213
fi

if [[ -z "$( which echo )" ]]; then
    exit 214
fi

# rename will happen in translation
if [[ -z "$( which sed )" ]]; then
    echo "imported command sed could not be found"
    exit 215
fi