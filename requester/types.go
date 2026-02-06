package requester

type ParticipantRequest struct {
	Active               bool   `json:"active"`
	ParticipantContextId string `json:"participantContextId"`
	DID                  string `json:"did"`
	Key                  Key    `json:"key"`
}

type Key struct {
	KeyID              string             `json:"keyId"`
	PrivateKeyAlias    string             `json:"privateKeyAlias"`
	KeyGeneratorParams KeyGeneratorParams `json:"keyGeneratorParams"`
}

type KeyGeneratorParams struct {
	Algorithm string `json:"algorithm"`
	Curve     string `json:"curve"`
}

func NewParticipantRequest(did string) ParticipantRequest {
	req := ParticipantRequest{
		Active:               true,
		ParticipantContextId: did,
		DID:                  did,
		Key: Key{
			KeyID:           did + "#test-key",
			PrivateKeyAlias: did + "-test-key-alias",
			KeyGeneratorParams: KeyGeneratorParams{
				Algorithm: "EdDSA",
				Curve:     "Ed25519",
			},
		},
	}

	return req
}
