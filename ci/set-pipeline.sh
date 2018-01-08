#!/usr/bin/env bash

TARGET="ci"
FLY="fly -t $TARGET"
SCRIPT_DIR="$(dirname $0)"
PIPELINE="$SCRIPT_DIR/pipeline.yml"
CREDENTIALS="$SCRIPT_DIR/credentials.yml"
CREDENTIALS_ENC="$CREDENTIALS.asc"

# decrypt credentials only on demand
if [[ ! -f "$CREDENTIALS" ]]; then
    gpg -d "${CREDENTIALS_ENC}"
    GPG_EXIT_STATUS=$?
    if [[ $GPG_EXIT_STATUS -ne 0 ]]; then
        echo "$CREDENTIALS decryption failed. Please, make sure you are one of the recipients."
        exit $GPG_EXIT_STATUS
    fi
fi

$FLY set-pipeline -p test -c $PIPELINE -l $CREDENTIALS

# cleanup
if [[ "$REMOVE_DECRYPTED_FILE" == 1 ]]; then
    rm -f "$CREDENTIALS"
fi

# todo: encrypt credentials.yml
