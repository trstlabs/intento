#!/bin/bash

set -eu
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
source ${SCRIPT_DIR}/config.sh

BUILDDIR="$2"
mkdir -p $BUILDDIR

# Build the local binary and docker image
build_local_and_docker() {
   set +e

   module="$1"
   folder="$2"
   title=$(printf "$module" | awk '{ print toupper($0) }')

   printf '%s' "Building $title Locally...  "

   intento_home=$PWD
   cd $folder

   # Clear any previously build binaries, otherwise the binary can get corrupted
   if [[ "$module" == "intento" ]]; then
      rm -f build/intentod
   else
      rm -f build/*
   fi

   # Many projects have a "check_version" in their makefile that prevents building
   # the binary if the machine's go version does not match exactly,
   # however, we can relax this constraint
   # The following command overrides the check_version using a temporary Makefile override
   BUILDDIR=$BUILDDIR make -f Makefile -f <(echo -e 'check_version: ;') build --silent 
   local_build_succeeded=${PIPESTATUS[0]}
   cd $intento_home

   # Some projects have a hard coded build directory, while others allow the passing of BUILDDIR
   # In the event that they have it hard coded, this will copy it into our build directory
   mv $folder/build/* $BUILDDIR/ > /dev/null 2>&1
   mv $folder/bin/* $BUILDDIR/ > /dev/null 2>&1

   if [[ "$local_build_succeeded" == "0" ]]; then
      echo "Done" 
   else
      echo "Failed"
      return $local_build_succeeded
   fi

   echo "Building $title Docker... "
   if [[ "$module" == "intento" ]]; then
      image=Dockerfile
   else
      image=dockernet/dockerfiles/Dockerfile.$module
   fi

   DOCKER_BUILDKIT=1 docker build --tag intento:$module -f $image . 
   docker_build_succeeded=${PIPESTATUS[0]}

   if [[ "$docker_build_succeeded" == "0" ]]; then
      echo "Done" 
   else
      echo "Failed"
   fi

   set -e
   return $docker_build_succeeded
}

# build docker images and local binaries
while getopts tgjosehrn flag; do
   case "${flag}" in
      t) build_local_and_docker intento . ;;
      g) build_local_and_docker gaia deps/gaia ;;
      j) build_local_and_docker juno deps/juno ;;
      o) build_local_and_docker osmo deps/osmosis ;;
      s) build_local_and_docker stars deps/stargaze ;;
      e) build_local_and_docker evmos deps/evmos ;;
      n) continue ;; # build_local_and_docker {new-host-zone} deps/{new-host-zone} ;;
      r) build_local_and_docker relayer deps/relayer ;;  
      h) echo "Building Hermes Docker... ";
         docker build --tag intento:hermes -f dockernet/dockerfiles/Dockerfile.hermes . ;

         printf '%s' "Building Hermes Locally... ";
         cd deps/hermes; 
         cargo build --release --target-dir $BUILDDIR/hermes; 
         cd ../..
         echo "Done" ;;
   esac
done