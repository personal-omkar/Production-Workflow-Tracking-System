#!/bin/bash

# Get the name of the currently active branch
current_branch=$(git branch --show-current)

# Loop through all branches and delete those not matching `main` or the current branch
for branch in $(git branch | grep -vE "^\*|main|$current_branch"); do
  git branch -D $branch
done
