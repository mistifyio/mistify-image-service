#/usr/bin/env bash
HOST=${1:-127.0.0.1}
PORT=${2:-20000}
FILE=${3:-/home/vagrant/ubuntu.zfs.gz}

prefix () {
    NOW=$(date +"%Y/%m/%d %H:%M:%S")
    echo "[ $NOW ] "
}

log () {
    echo -e "$(prefix)$@"
}

indent () {
    PF=$(prefix)
    LENGTH=${#PF}
    printf -v SPACE '%*s' "$LENGTH"
    sed "s/^/$SPACE/";
}

header () {
    echo ""
    log "===== $@ ====="
}

http () {
    METHOD=$1
    ENDPOINT=$2
    shift 2

    if [[ -n "$XIT" ]]
    then
        log "HEADER\t$XIT"
    fi

    if [[ -n "$XIC" ]]
    then
        log "HEADER\t$XIC"
    fi

    URL="http://$HOST:$PORT/$ENDPOINT"
    log "$METHOD\t$URL" 
    OUTPUT=$(curl --fail -s -X $METHOD -H "$XIT" -H "$XIC" -H 'Content-Type: application/json' $URL "$@" | jq .)
    log "Result:"
    echo "$OUTPUT" | indent

    #Unset headers
    XIT=
    XIC=
}

clean () {
    CLEANIDS=($(curl --fail -s -X GET "http://$HOST:$PORT/images" | jq -r .[].ID))
    for CLEANID in "${CLEANIDS[@]}"
    do
        _=$(curl -s -X DELETE "http://$HOST:$PORT/images/$CLEANID")
    done
}

clean

header "LIST IMAGES"
http GET images

header "UPLOAD IMAGE"
XIT="X-Image-Type: kvm"
XIC="X-Image-Comment: uploaded image"
http PUT images --data-binary "@$FILE"
ID=$(echo "$OUTPUT" | jq -r .ID)

header "GET IMAGE INFO"
http GET images/$ID

header "FETCH IMAGE"
URL="http://$HOST:$PORT/images/$ID/download"
http POST images --data-binary '{"source":"'$URL'","type":"container","comment":"fetched image"}'

header "LIST IMAGES"
http GET images

header "LIST KVM IMAGES"
http GET 'images?type=kvm'

header "LIST CONTAINER IMAGES"
http GET 'images?type=container'

header "DELETE IMAGE"
http DELETE images/$ID

clean
