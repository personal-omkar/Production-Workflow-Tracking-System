#!/bin/bash

goProjects=("../")
         
for i in "${goProjects[@]}"
do
    echo --------------------$i------------------ | tr [a-z] [A-Z]
    cd $i 
    git branch   
    git fetch
    git pull
    cd -
done