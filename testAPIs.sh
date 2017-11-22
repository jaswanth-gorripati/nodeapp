#!/bin/bash
#
# Copyright IBM Corp. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#

jq --version > /dev/null 2>&1
if [ $? -ne 0 ]; then
	echo "Please Install 'jq' https://stedolan.github.io/jq/ to execute this script"
	echo
	exit 1
fi
starttime=$(date +%s)

echo "POST request Enroll on Org1  ..."
echo
ORG1_TOKEN=$(curl -s -X POST \
  http://localhost:4000/users \
  -H "content-type: application/x-www-form-urlencoded" \
  -d 'username=scadmin&passkey=pass&orgName=org1')
echo eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MTAzMzQ2MTQsInVzZXJuYW1lIjoic2NhZG1pbiIsIm9yZ05hbWUiOiJvcmcxIiwiaWF0IjoxNTEwMjk4NjE0fQ.cb7oNNocRGYyH_feUwamJ6aryQ5jM5T1H_LPFj2Otek
ORG1_TOKEN=$(echo eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MTAzMzQ2MTQsInVzZXJuYW1lIjoic2NhZG1pbiIsIm9yZ05hbWUiOiJvcmcxIiwiaWF0IjoxNTEwMjk4NjE0fQ.cb7oNNocRGYyH_feUwamJ6aryQ5jM5T1H_LPFj2Otek | jq ".token" | sed "s/\"//g")
echo
echo "ORG1 token is eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MTAzMzQ2MTQsInVzZXJuYW1lIjoic2NhZG1pbiIsIm9yZ05hbWUiOiJvcmcxIiwiaWF0IjoxNTEwMjk4NjE0fQ.cb7oNNocRGYyH_feUwamJ6aryQ5jM5T1H_LPFj2Otek"
echo
echo "POST request Enroll on Org2 ..."
echo
ORG2_TOKEN=$(curl -s -X POST \
  http://localhost:4000/users \
  -H "content-type: application/x-www-form-urlencoded" \
  -d 'username=clgadmin&passkey=pass&orgName=org2')
echo $ORG2_TOKEN
ORG2_TOKEN=$(echo $ORG2_TOKEN | jq ".token" | sed "s/\"//g")
echo
echo "ORG2 token is $ORG2_TOKEN"
echo
echo
echo "POST request Create channel 2  ..."
echo
curl -s -X POST \
  http://localhost:4000/channels \
  -H "authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MTAzMzQ2MTQsInVzZXJuYW1lIjoic2NhZG1pbiIsIm9yZ05hbWUiOiJvcmcxIiwiaWF0IjoxNTEwMjk4NjE0fQ.cb7oNNocRGYyH_feUwamJ6aryQ5jM5T1H_LPFj2Otek" \
  -H "content-type: application/json" \
  -d '{
  "channelName":"mychannel",
  "channelConfigPath":"../artifacts/channel/mychannel.tx"
}'
echo
echo
sleep 5
echo
echo "POST request Create channel 1  ..."
echo
curl -s -X POST \
  http://localhost:4000/channels \
  -H "authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MTAzMzQ2MTQsInVzZXJuYW1lIjoic2NhZG1pbiIsIm9yZ05hbWUiOiJvcmcxIiwiaWF0IjoxNTEwMjk4NjE0fQ.cb7oNNocRGYyH_feUwamJ6aryQ5jM5T1H_LPFj2Otek" \
  -H "content-type: application/json" \
  -d '{
	"channelName":"mychannel",
	"channelConfigPath":"../artifacts/channel/mychannel.tx"
}'
echo
echo
sleep 5

echo "POST request Join channel on Org1"
echo
curl -s -X POST \
  http://localhost:4000/channels/mychannel/peers \
  -H "authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MTAzMzQ2MTQsInVzZXJuYW1lIjoic2NhZG1pbiIsIm9yZ05hbWUiOiJvcmcxIiwiaWF0IjoxNTEwMjk4NjE0fQ.cb7oNNocRGYyH_feUwamJ6aryQ5jM5T1H_LPFj2Otek" \
  -H "content-type: application/json" \
  -d '{
	"peers": ["peer1","peer2"]
}'
echo
echo

echo "POST request Join channel on Org2"
echo
curl -s -X POST \
  http://localhost:4000/channels/mychannel/peers \
  -H "authorization: Bearer $ORG2_TOKEN" \
  -H "content-type: application/json" \
  -d '{
	"peers": ["peer1","peer2"]
}'
echo
echo

echo "POST request Join channel2 on Org2"
echo
curl -s -X POST \
  http://localhost:4000/channels/registration/peers \
  -H "authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MTAzMzQ2MTQsInVzZXJuYW1lIjoic2NhZG1pbiIsIm9yZ05hbWUiOiJvcmcxIiwiaWF0IjoxNTEwMjk4NjE0fQ.cb7oNNocRGYyH_feUwamJ6aryQ5jM5T1H_LPFj2Otek" \
  -H "content-type: application/json" \
  -d '{
  "peers": ["peer1","peer2"]
}'
echo
echo

echo "POST Install chaincode on Org1"
echo
curl -s -X POST \
  http://localhost:4000/chaincodes \
  -H "authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MTAzMzQ2MTQsInVzZXJuYW1lIjoic2NhZG1pbiIsIm9yZ05hbWUiOiJvcmcxIiwiaWF0IjoxNTEwMjk4NjE0fQ.cb7oNNocRGYyH_feUwamJ6aryQ5jM5T1H_LPFj2Otek" \
  -H "content-type: application/json" \
  -d '{
	"peers": ["peer1", "peer2"],
	"chaincodeName":"mycc",
	"chaincodePath":"github.com/example_cc",
	"chaincodeVersion":"v0"
}'
echo
echo


echo "POST Install chaincode on Org2"
echo
curl -s -X POST \
  http://localhost:4000/chaincodes \
  -H "authorization: Bearer $ORG2_TOKEN" \
  -H "content-type: application/json" \
  -d '{
	"peers": ["peer1","peer2"],
	"chaincodeName":"mycc",
	"chaincodePath":"github.com/example_cc",
	"chaincodeVersion":"v0"
}'
echo
echo

echo "POST instantiate chaincode on peer1 of Org1"
echo
curl -s -X POST \
  http://localhost:4000/channels/mychannel/chaincodes \
  -H "authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MTAzMzQ2MTQsInVzZXJuYW1lIjoic2NhZG1pbiIsIm9yZ05hbWUiOiJvcmcxIiwiaWF0IjoxNTEwMjk4NjE0fQ.cb7oNNocRGYyH_feUwamJ6aryQ5jM5T1H_LPFj2Otek" \
  -H "content-type: application/json" \
  -d '{
	"chaincodeName":"mycc",
	"chaincodeVersion":"v0",
	"args":["a","100","b","200"]
}'
echo
echo
sleep 3
echo "POST invoke chaincode on peers of Org1 and Org2"
echo
TRX_ID=$(curl -s -X POST \
  http://localhost:4000/channels/mychannel/chaincodes/mycc \
  -H "authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MTAzMzQ2MTQsInVzZXJuYW1lIjoic2NhZG1pbiIsIm9yZ05hbWUiOiJvcmcxIiwiaWF0IjoxNTEwMjk4NjE0fQ.cb7oNNocRGYyH_feUwamJ6aryQ5jM5T1H_LPFj2Otek" \
  -H "content-type: application/json" \
  -d '{
  "peers": ["peer1","peer2"],
    "fcn":"register",
    "args":["100","emoji","jash","1994-05-03","male","admin","2017-05-04","ssc","stateboard","rkp","2010","8","jash","2017-02-05"]
}')
echo "Transacton ID is $TRX_ID"
echo
echo
echo "get channel information"
echo
TRX_ID=$(curl -s -X GET \
  http://localhost:4000/channels/mychannel/ \
  -H "authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MTAzMzQ2MTQsInVzZXJuYW1lIjoic2NhZG1pbiIsIm9yZ05hbWUiOiJvcmcxIiwiaWF0IjoxNTEwMjk4NjE0fQ.cb7oNNocRGYyH_feUwamJ6aryQ5jM5T1H_LPFj2Otek" \
  -H "content-type: application/json" \
  -d '{
}')
echo "Transacton ID is $TRX_ID"
echo
curl -s -X GET \
  http://localhost:4000/chaincode?peer='peer1'&&type='installed' \
  -H "authorization: Bearer  " \
  -H "content-type: application/json" \
  -d '{
}'

