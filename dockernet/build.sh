#!/bin/bash

set -eu
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
source ${SCRIPT_DIR}/config.sh

BUILDDIR="$2"
mkdir -p $BUILDDIR
# Build the local binary
build_local() {
   set +e

   cli="$1"
   folder="$2"
   title=$(printf "$cli" | awk '{ print toupper($0) }')

   printf '%s' "Building $cli Locally...  "

   intento_home=$PWD
   cd $folder

   rm -f build/"$cli"


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

   set -e
   return $local_build_succeeded
}
# Build the Docker image
build_docker() {
   set +e

   module="$1"
   folder="$2"
   title=$(printf "$module" | awk '{ print toupper($0) }')

   echo "Building $title Docker... "
   if [[ "$module" == "into" ]]; then
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


# Build local binaries and Docker images
while getopts igdosehrn flag; do
   case "${flag}" in
      i)
         build_local intentod .
         build_docker into .
         ;;
      g)
         build_local gaiad deps/gaia
         build_docker gaia deps/gaia
         ;;
      o)
         build_local osmosisd deps/osmosis
         build_docker osmo deps/osmosis
         ;;
      n) continue ;; # build_local and build_docker {new-host-zone} deps/{new-host-zone}
      r)
         build_local rly deps/relayer
         build_docker relayer deps/relayer
         ;;
      h)
         echo "Building Hermes Docker... "
         docker build --tag intento:hermes -f dockernet/dockerfiles/Dockerfile.hermes .
         ;;
   esac
done