#!/bin/bash
set -e

if [[ -f .env ]]; then
	for line in $(cat .env); do
		export $line
	done
fi

composer install --no-dev
php translate.php
