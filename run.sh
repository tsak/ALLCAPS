#!/bin/bash

if [ ! -f "./allcaps" ]; then
  echo "$0: missing ./allcaps binary"
  exit 1
fi

if [ ! -f ".env" ]; then
  echo "$0: missing .env file"
  exit 1
fi

source .env

if [ -z "$SLACKTOKEN" ]; then
  echo "$0: SLACKTOKEN env var not set"
  exit 1
fi

export SLACKTOKEN

./allcaps