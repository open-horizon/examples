# SDR service API

Assuming that `sdr` is the host name:

## `/freqs`
Get a list of the frequencies of strong radio stations.

`curl sdr:5427/freqs`

Example response:
If the SDR hardware is present it will return a list of string FM stations.
`{"origin":"sdr_hardware","freqs":[89700000,91100000,91900000,93300000,94500000,95700000,96100000,97700000,99100000,101500000,102300000,103700000,105100000,107900000]}`

If the SDR hardware is not present or can not be used for some reason it will return a single station of frequency 0.
`{"origin":"fake","freqs":[0]}`

## /audio/<freq>
Get a 30 second chunk of raw audio.
`curl sdr:5427/audio/99100000`

`curl sdr:5427/audio/99100000 | aplay -r 16000 -f S16_LE -t raw -c 1`
