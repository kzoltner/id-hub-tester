# Testing the Identity Hub
An error occured during one test run of the local docker compose setup. 
Since it has occured once again, this repository contains a test to find out if there is actually something wrong.

The repository is using Go and assumes, that you already have a local docker compose setup up and running. 

It will create a whole lot of requests and it will log out runs that failed to logs/failed as well as some random runs to ok for comparison.

To run, you need Go installed (>= 1.25 at least). Then run:

~~~
go run .
~~~

from the project root folder. 

## Structure of jsons
The produced jsons in failed contain: 

- Did: DID used for request
- ParticipantRequest: The actual data that was sent (not encoded as it is from the struct directly)
- ParticipantRequestResponse: What the server responded with
- DidDocumentResponse: Just the did document of this run
- KeyPairResponse: left empty, but the [IdentityHub API](https://eclipse-edc.github.io/IdentityHub/openapi/identity-api/#/) has a KeyPair endpoint to get all of them - just remember to put ?limit=10000 to get all.

# Results
So far, I was able to reproduce the error about 20 times for every 5000 runs.

What I was able to see: 
None of them seem out of place at all. Apart from the fact, that all failed request have an "X" Property within their according keypairs.json entry (identified using just test-id-xxxx) which has one char less then the others. But that could simply be a symptom instead of the actual cause.

The logs of consumer_id_hub always show the same error - and a lot of warnings about some state? But that is always there, even on good ones.
